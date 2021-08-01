package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"strconv"
	"time"

	"github.com/fiatjaf/go-lnurl"
	"github.com/gorilla/mux"
	"github.com/tidwall/sjson"
)

func lnurlParams(w http.ResponseWriter, r *http.Request) {
	username := mux.Vars(r)["username"]
	domain := r.Header.Get("Host")

	// check if the receiver accepts comments
	var commentLength int64 = 0
	if v, err := net.LookupTXT("_webhook." + domain); err == nil && len(v) > 0 {
		commentLength = 500
	}

	json.NewEncoder(w).Encode(lnurl.LNURLPayResponse1{
		LNURLResponse:   lnurl.LNURLResponse{Status: "OK"},
		Callback:        fmt.Sprintf("%s/lnurl/%s", domain, username),
		MinSendable:     1000,
		MaxSendable:     100000000,
		EncodedMetadata: makeMetadata(username, domain),
		CommentAllowed:  commentLength,
		Tag:             "payRequest",
	})
}

func lnurlValues(w http.ResponseWriter, r *http.Request) {
	username := mux.Vars(r)["username"]
	domain := r.Header.Get("Host")

	amount := r.URL.Query().Get("amount")
	msat, err := strconv.Atoi(amount)
	if err != nil {
		json.NewEncoder(w).Encode(lnurl.ErrorResponse("amount is not integer"))
		return
	}

	bolt11, err := makeInvoice(username, domain, msat)
	if err != nil {
		json.NewEncoder(w).Encode(
			lnurl.ErrorResponse("failed to create invoice: " + err.Error()))
		return
	}

	json.NewEncoder(w).Encode(lnurl.LNURLPayResponse2{
		LNURLResponse: lnurl.LNURLResponse{Status: "OK"},
		PR:            bolt11,
		Routes:        make([][]lnurl.RouteInfo, 0),
		Disposable:    lnurl.FALSE,
		SuccessAction: lnurl.Action("Payment received!", ""),
	})

	// handle comments by sending them to a specified webhook
	go func() {
		if comment := r.URL.Query().Get("comment"); comment != "" {
			if v, err := net.LookupTXT("_webhook." + domain); err == nil && len(v) > 0 {
				body, _ := sjson.Set("{}", "comment", comment)
				body, _ = sjson.Set(body, "pr", bolt11)
				body, _ = sjson.Set(body, "amount", msat)
				(&http.Client{Timeout: 5 * time.Second}).
					Post(v[0], "application/json", bytes.NewBufferString(body))
			}
		}
	}()
}
