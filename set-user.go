package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	decodepay "github.com/fiatjaf/ln-decodepay"
	"github.com/gorilla/mux"
	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
)

func jsonErrorf(str string, args ...interface{}) (j struct {
	Ok    bool   `json:"ok"`
	Error string `json:"error"`
}) {
	j.Error = fmt.Sprintf(str, args)
	return
}

func setUser(w http.ResponseWriter, r *http.Request) {
	kind := mux.Vars(r)["kind"]

	defer r.Body.Close()
	jdata, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(400)
		json.NewEncoder(w).Encode(jsonErrorf("failed to read request: %w", err))
		log.Warn().Err(err).Msg("failed to read input")
		return
	}

	// make a test invoice so we can get the node id
	bolt11, err := makeInvoice(kind, string(jdata), 1000)
	if err != nil {
		w.WriteHeader(401)
		json.NewEncoder(w).Encode(jsonErrorf("failed to create test invoice: %w", err))
		log.Warn().Err(err).Msg("failed to create test invoice")
		return
	}

	// get node id from invoice
	inv, err := decodepay.Decodepay(bolt11)
	if err != nil {
		w.WriteHeader(417)
		json.NewEncoder(w).Encode(jsonErrorf("failed to parse test invoice: %w", err))
		log.Warn().Err(err).Msg("failed to parse test invoice")
		return
	}
	id := inv.Payee

	// get only relevant fields
	data := gjson.ParseBytes(jdata)
	inputdata := "{}"
	inputdata, _ = sjson.Set(inputdata, "endpoint", data.Get("endpoint").String())
	if data.Get("cert").Exists() {
		inputdata, _ = sjson.Set(inputdata, "cert", data.Get("cert").String())
	}
	if kind == "lnd" {
		inputdata, _ = sjson.Set(inputdata, "macaroon", data.Get("macaroon").String())
	} else if kind == "sparko" {
		inputdata, _ = sjson.Set(inputdata, "key", data.Get("key").String())
	}

	// save to database
	_, err = pg.Exec(`
INSERT INTO users (id, kind, data)
VALUES ($1, $2, $3)
ON CONFLICT (id) DO UPDATE
  SET kind = $2, data = $3
    `, id, kind, inputdata)
	if err != nil {
		w.WriteHeader(400)
		json.NewEncoder(w).Encode(jsonErrorf("failed to save: %w", err))
		log.Warn().Err(err).Msg("failed to insert or update user")
		return
	}

	w.WriteHeader(200)
	json.NewEncoder(w).Encode(struct {
		Ok bool   `json:"ok"`
		Id string `json:"id"`
	}{true, id})
	log.Info().Str("id", id).Msg("saved an user")
}
