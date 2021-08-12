package main

import (
	"bytes"
	"crypto/sha256"
	"crypto/tls"
	"crypto/x509"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"strconv"
	"time"

	lightning "github.com/fiatjaf/lightningd-gjson-rpc"
	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
)

func makeMetadata(username, domain string) string {
	metadata, _ := sjson.Set("[]", "0.0", "text/identifier")
	metadata, _ = sjson.Set(metadata, "0.1", username+"@"+domain)

	metadata, _ = sjson.Set(metadata, "1.0", "text/plain")
	if v, err := net.LookupTXT("_description." + domain); err == nil && len(v) > 0 {
		metadata, _ = sjson.Set(metadata, "1.1", v[0])
	} else {
		metadata, _ = sjson.Set(metadata, "1.1", "Satoshis to "+username+"@"+domain+".")
	}

	if v, err := net.LookupTXT("_image." + domain); err == nil && len(v) > 0 {
		if b64, err := base64ImageFromURL(v[0]); err == nil {
			metadata, _ = sjson.Set(metadata, "2.0", "image/jpeg;base64")
			metadata, _ = sjson.Set(metadata, "2.1", b64)
		}
	}
	return metadata
}

func makeInvoice(username, domain string, msat int) (bolt11 string, err error) {
	defer func(prevTransport http.RoundTripper) {
		http.DefaultClient.Transport = prevTransport
	}(http.DefaultClient.Transport)

	// grab all the necessary data from DNS
	var (
		kind     string
		cert     string
		host     string
		key      string
		macaroon string
	)
	if v, err := net.LookupTXT("_kind." + domain); err == nil && len(v) > 0 {
		kind = v[0]
	} else {
		return "", errors.New("missing kind")
	}
	if v, err := net.LookupTXT("_cert." + domain); err == nil && len(v) > 0 {
		cert = v[0]
	}
	if v, err := net.LookupTXT("_host." + domain); err == nil && len(v) > 0 {
		host = v[0]
	}
	if v, err := net.LookupTXT("_key." + domain); err == nil && len(v) > 0 {
		key = v[0]
	}
	if v, err := net.LookupTXT("_macaroon." + domain); err == nil && len(v) > 0 {
		macaroon = v[0]
	}

	// use a cert or skip TLS verification?
	if cert != "" {
		caCertPool := x509.NewCertPool()
		caCertPool.AppendCertsFromPEM([]byte(cert))
		http.DefaultClient.Transport = &http.Transport{
			TLSClientConfig: &tls.Config{RootCAs: caCertPool},
		}
	} else {
		http.DefaultClient.Transport = &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
	}

	// prepare these things
	h := sha256.Sum256([]byte(makeMetadata(username, domain)))
	hexh := hex.EncodeToString(h[:])
	b64h := base64.StdEncoding.EncodeToString(h[:])

	// actually generate the invoice
	switch kind {
	case "sparko":
		spark := &lightning.Client{
			SparkURL:    host,
			SparkToken:  key,
			CallTimeout: time.Second * 3,
		}
		inv, err := spark.Call("invoicewithdescriptionhash", msat,
			"lightningaddr/"+strconv.FormatInt(time.Now().UnixNano(), 16), hexh)
		if err != nil {
			return "", fmt.Errorf("invoicewithdescriptionhash call failed: %w", err)
		}
		return inv.Get("bolt11").String(), nil

	case "lnd":
		body, _ := sjson.Set("{}", "description_hash", b64h)
		body, _ = sjson.Set(body, "value", msat/1000)

		req, err := http.NewRequest("POST",
			host+"/v1/invoices",
			bytes.NewBufferString(body),
		)
		if err != nil {
			return "", err
		}

		req.Header.Set("Grpc-Metadata-macaroon", macaroon)
		resp, err := (&http.Client{Timeout: 25 * time.Second}).Do(req)
		if err != nil {
			return "", err
		}
		if resp.StatusCode >= 300 {
			return "", errors.New("call to lnd failed")
		}

		defer resp.Body.Close()
		b, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return "", err
		}

		return gjson.ParseBytes(b).Get("payment_request").String(), nil
	case "lnpay":
	case "lnbits":
	}

	return "", errors.New("unsupported lightning server kind: " + kind)
}
