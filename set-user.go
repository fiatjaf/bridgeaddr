package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	decodepay "github.com/fiatjaf/ln-decodepay"
	"github.com/gorilla/mux"
)

func jsonErrorf(str string, args ...interface{}) (j struct {
	Ok    bool   `json:"ok"`
	Error string `json:"error"`
}) {
	j.Error = fmt.Sprintf(str, args)
	return
}

func setUser(w http.ResponseWriter, r *http.Request) {
	kind := mux.Vars(r)["id"]

	defer r.Body.Close()
	jdata, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(400)
		json.NewEncoder(w).Encode(jsonErrorf("failed to read request: %w", err))
		log.Warn().Err(err).Msg("failed to read input")
		return
	}
	sjdata := string(jdata)

	// make a test invoice so we can get the node id
	bolt11, err := makeInvoice(kind, sjdata, 1000)
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

	_, err = pg.Exec(`
INSERT INTO users (id, kind, data)
VALUES ($1, $2, $3)
ON CONFLICT (id) DO UPDATE
  SET kind = $2, data = $3
    `, id, kind, sjdata)
	if err != nil {
		w.WriteHeader(400)
		json.NewEncoder(w).Encode(jsonErrorf("failed to save: %w", err))
		log.Warn().Err(err).Msg("failed to insert or update user")
		return
	}

	w.WriteHeader(200)
}
