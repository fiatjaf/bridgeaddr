package main

import (
	"crypto/sha256"
	"errors"
	"net"
	"strconv"
	"time"

	"github.com/fiatjaf/makeinvoice"
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
	// grab all the necessary data from DNS
	var (
		kind     string
		cert     string
		host     string
		key      string
		macaroon string
		pak      string
		waki     string
		nodeid   string
		rune_    string
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
	// description_hash
	h := sha256.Sum256([]byte(makeMetadata(username, domain)))

	// prepare params
	var backend makeinvoice.BackendParams
	switch kind {
	case "sparko":
		if v, err := net.LookupTXT("_key." + domain); err == nil && len(v) > 0 {
			key = v[0]
		}

		backend = makeinvoice.SparkoParams{
			Cert: cert,
			Host: host,
			Key:  key,
		}
	case "commando":
		if v, err := net.LookupTXT("_nodeid." + domain); err == nil && len(v) > 0 {
			nodeid = v[0]
		}
		if v, err := net.LookupTXT("_rune." + domain); err == nil && len(v) > 0 {
			rune_ = v[0]
		}

		backend = makeinvoice.CommandoParams{
			Host:   host,
			NodeId: nodeid,
			Rune:   rune_,
		}
	case "eclair":
		backend = makeinvoice.EclairParams{
			Cert: cert,
			Host: host,
		}
	case "lnd":
		if v, err := net.LookupTXT("_macaroon." + domain); err == nil && len(v) > 0 {
			macaroon = v[0]
		}

		backend = makeinvoice.LNDParams{
			Cert:     cert,
			Host:     host,
			Macaroon: macaroon,
		}
	case "lnbits":
		if v, err := net.LookupTXT("_key." + domain); err == nil && len(v) > 0 {
			key = v[0]
		}

		backend = makeinvoice.LNBitsParams{
			Cert: cert,
			Host: host,
			Key:  key,
		}
	case "lnpay":
		if v, err := net.LookupTXT("_pak." + domain); err == nil && len(v) > 0 {
			pak = v[0]
		}

		if v, err := net.LookupTXT("_waki." + domain); err == nil && len(v) > 0 {
			waki = v[0]
		}

		backend = makeinvoice.LNPayParams{
			PublicAccessKey:  pak,
			WalletInvoiceKey: waki,
		}
	}

	// actually generate the invoice
	return makeinvoice.MakeInvoice(makeinvoice.Params{
		Msatoshi:        int64(msat),
		DescriptionHash: h[:],
		Backend:         backend,

		Label: "bridgeaddr/" + strconv.FormatInt(time.Now().Unix(), 16),
	})
}
