package main

import (
	_ "embed"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/fiatjaf/makeinvoice"
	"github.com/gorilla/mux"
	"github.com/kelseyhightower/envconfig"
	"github.com/rs/zerolog"
)

const NOTICE = "NOTICE: bridgeaddr is shutting down at the end of 2024, please move your address to somewhere else."

type Settings struct {
	Host             string `envconfig:"HOST" default:"0.0.0.0"`
	Port             string `envconfig:"PORT" required:"true"`
	ServiceURL       string `envconfig:"SERVICE_URL" required:"true"`
	SafeDomainSuffix string `envconfig:"SAFE_DOMAIN_SUFFIX"`
}

var (
	err    error
	s      Settings
	router = mux.NewRouter()
	log    = zerolog.New(os.Stderr).Output(zerolog.ConsoleWriter{Out: os.Stderr})
)

//go:embed README.md
var readme string

func main() {
	err = envconfig.Process("", &s)
	if err != nil {
		log.Fatal().Err(err).Msg("couldn't process envconfig.")
	}

	// increase default makeinvoice client timeout because people are using tor
	makeinvoice.Client = &http.Client{Timeout: 25 * time.Second}

	// this is here so caddy can call it and validate new certificate requests
	router.Path("/domain-validate").Methods("GET").
		HandlerFunc(handleDomainValidate)

		// the core endpoint that handles lightning address calls
	router.Path("/.well-known/lnurlp/{username}").Methods("GET").
		HandlerFunc(handleLNURL)

	router.Host(strings.Split(s.ServiceURL, "://")[1]).Path("/").HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("content-type", "text/html")
			fmt.Fprintf(w, readme+`
<style class="fallback">body{visibility:hidden;white-space:pre;font-family:monospace}</style></script><script src="https://casual-effects.com/markdeep/latest/markdeep.min.js" charset="utf-8"></script><script>window.alreadyProcessedMarkdeep||(document.body.style.visibility="visible")</script>

<script>
setTimeout(() => {
  document.body.innerHTML += `+"`"+`
    <div style="position: fixed;top:20px;right:20px;width:300px;padding:10px;font-size:1.2rem;border:2px solid;background: #f3d797;filter:drop-shadow(11px 12px 8px #efdede);"
    >`+NOTICE+`</div>
  `+"`"+`
}, 3000)

</script>
`)
		},
	)

	router.PathPrefix("/").Methods("GET").HandlerFunc(redirect)

	srv := &http.Server{
		Handler:      router,
		Addr:         s.Host + ":" + s.Port,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	log.Debug().Str("addr", srv.Addr).Msg("listening")
	srv.ListenAndServe()
}
