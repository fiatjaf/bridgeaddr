package main

import (
	"net/http"
	"os"
	"time"

	assetfs "github.com/elazarl/go-bindata-assetfs"
	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
	"github.com/kelseyhightower/envconfig"
	_ "github.com/lib/pq"
	"github.com/rs/zerolog"
)

type Settings struct {
	Port        string `envconfig:"PORT" required:"true"`
	ServiceURL  string `envconfig:"SERVICE_URL" required:"true"`
	PostgresURL string `envconfig:"DATABASE_URL" required:"true"`
}

var err error
var s Settings
var pg *sqlx.DB
var router = mux.NewRouter()
var log = zerolog.New(os.Stderr).Output(zerolog.ConsoleWriter{Out: os.Stderr})

func main() {
	err = envconfig.Process("", &s)
	if err != nil {
		log.Fatal().Err(err).Msg("couldn't process envconfig.")
	}

	// postgres connection
	pg, err = sqlx.Connect("postgres", s.PostgresURL)
	if err != nil {
		log.Fatal().Err(err).Msg("couldn't connect to postgres")
	}

	// files
	assets := &assetfs.AssetFS{Asset: Asset, AssetDir: AssetDir, Prefix: "/public/"}
	indexhtml := MustAsset("public/index.html")
	donatehtml := MustAsset("public/donate.html")

	router.Path("/set/{kind}").Methods("PUT").HandlerFunc(setUser)
	router.Path("/lnurl/{id}/params").Methods("Get").HandlerFunc(lnurlParams)
	router.Path("/lnurl/{id}/values").Methods("Get").HandlerFunc(lnurlValues)
	router.PathPrefix("/public/").Handler(http.StripPrefix("/public/", http.FileServer(assets)))
	router.Path("/").Methods("GET").HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			w.Header().Add("Content-Type", "text/html")
			w.Write(indexhtml)
		},
	)
	router.Path("/{id}").Methods("GET").HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			w.Header().Add("Content-Type", "text/html")
			w.Write(donatehtml)
		},
	)

	srv := &http.Server{
		Handler:      router,
		Addr:         "0.0.0.0:" + s.Port,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	log.Debug().Str("addr", srv.Addr).Msg("listening")
	srv.ListenAndServe()
}
