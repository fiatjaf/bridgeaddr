package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/fiatjaf/go-lnurl"
	decodepay "github.com/fiatjaf/ln-decodepay"
	"github.com/gorilla/mux"
)

func lnurlParams(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	json.NewEncoder(w).Encode(lnurl.LNURLPayResponse1{
		LNURLResponse:   lnurl.LNURLResponse{Status: "OK"},
		Callback:        fmt.Sprintf("%s/lnurl/%s/values", s.ServiceURL, id),
		MinSendable:     1000,
		MaxSendable:     100000000,
		EncodedMetadata: getMetadata(),
		Tag:             "payRequest",
	})
}

func lnurlValues(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	var kind string
	var jdata string
	err = pg.QueryRowx("SELECT kind, data::text FROM users WHERE id = $1", id).
		Scan(&kind, &jdata)
	if err != nil {
		json.NewEncoder(w).Encode(lnurl.ErrorResponse("user doesn't exist"))
		return
	}

	amount := r.URL.Query().Get("amount")
	msat, err := strconv.Atoi(amount)
	if err != nil {
		json.NewEncoder(w).Encode(lnurl.ErrorResponse("amount is not integer"))
		return
	}

	bolt11, err := makeInvoice(kind, jdata, msat)
	if err != nil {
		json.NewEncoder(w).Encode(
			lnurl.ErrorResponse("failed to create invoice: " + err.Error()))
		return
	}

	// check node id
	inv, err := decodepay.Decodepay(bolt11)
	if err != nil {
		json.NewEncoder(w).Encode(
			lnurl.ErrorResponse("failed to parse invoice: " + err.Error()))
		return
	}
	if inv.Payee != id {
		log.Warn().Msg("generated invoice is not from the correct node id")
		json.NewEncoder(w).Encode(
			lnurl.ErrorResponse("got an invoice from the wrong node id"))
		return
	}

	json.NewEncoder(w).Encode(lnurl.LNURLPayResponse2{
		LNURLResponse: lnurl.LNURLResponse{Status: "OK"},
		PR:            bolt11,
		SuccessAction: lnurl.Action("Thank you!", ""),
		Routes:        make([][]lnurl.RouteInfo, 0),
	})
}
