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
	"net/http"
	"time"

	lightning "github.com/fiatjaf/lightningd-gjson-rpc"
	"github.com/lucsky/cuid"
	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
)

func getMetadata(nodeid string) string {
	metadata, _ := sjson.Set("[]", "0.0", "text/plain")
	metadata, _ = sjson.Set(metadata, "0.1", "An arbitrary payment to the Lightning node "+nodeid+" proxied by https://tip.bigsun.xyz/.")
	return metadata
}

func makeInvoice(nodeid, kind, jdata string, msat int) (bolt11 string, err error) {
	data := gjson.Parse(jdata)

	defer func(prevTransport http.RoundTripper) {
		http.DefaultClient.Transport = prevTransport
	}(http.DefaultClient.Transport)
	if data.Get("cert").Exists() {
		caCertPool := x509.NewCertPool()
		caCertPool.AppendCertsFromPEM([]byte(data.Get("cert").String()))

		http.DefaultClient.Transport = &http.Transport{
			TLSClientConfig: &tls.Config{RootCAs: caCertPool},
		}
	} else {
		http.DefaultClient.Transport = &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
	}

	h := sha256.Sum256([]byte(getMetadata(nodeid)))
	hexh := hex.EncodeToString(h[:])
	b64h := base64.StdEncoding.EncodeToString(h[:])

	switch kind {
	case "sparko":
		spark := &lightning.Client{
			SparkURL:    data.Get("endpoint").String(),
			SparkToken:  data.Get("key").String(),
			CallTimeout: time.Second * 3,
		}
		inv, err := spark.Call("invoicewithdescriptionhash", msat, "tip.bigsun.xyz/"+cuid.Slug(), hexh)
		fmt.Println(msat, "tip.bigsun.xyz/"+cuid.Slug(), hexh)
		if err != nil {
			return "", fmt.Errorf("lnurlinvoice call failed: %w", err)
		}
		return inv.Get("bolt11").String(), nil

	case "lnd":
		body, _ := sjson.Set("{}", "description_hash", b64h)
		body, _ = sjson.Set(body, "value", msat/1000)

		req, err := http.NewRequest("POST",
			data.Get("endpoint").String()+"/v1/invoices",
			bytes.NewBufferString(body),
		)
		if err != nil {
			return "", err
		}

		req.Header.Set("Grpc-Metadata-macaroon", data.Get("macaroon").String())
		resp, err := http.DefaultClient.Do(req)
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
	}

	return "", errors.New("unsupported lightning server kind: " + kind)
}
