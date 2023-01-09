package main

import (
	"net/http"
	"strings"
	"time"

	"golang.org/x/net/publicsuffix"
)

var stupidCache = make(map[string]time.Time)

func handleDomainValidate(w http.ResponseWriter, r *http.Request) {
	domain := r.URL.Query().Get("domain")
	log := log.With().Str("domain", domain).Logger()

	if s.SafeDomainSuffix != "" && strings.HasSuffix(domain, s.SafeDomainSuffix) {
		w.WriteHeader(200)
		return
	}

	effective, err := publicsuffix.EffectiveTLDPlusOne(domain)
	if err != nil {
		log.Warn().Err(err).Msg("failed to parse effective tld")
		w.WriteHeader(400)
		return
	}
	if when, ok := stupidCache[effective]; ok {
		if when.AddDate(0, 0, 3).After(time.Now()) {
			log.Warn().Msg("domain rejected because effective tld+1 was already used in the past 3 days")
			w.WriteHeader(400)
			return
		}
	}

	if len(domain) > 30 {
		log.Warn().Msg("domain rejected for being too large")
		w.WriteHeader(400)
		return
	}

	parts := strings.Split(domain, ".")
	if len(parts) > 3 {
		log.Warn().Msg("domain rejected for having too many parts")
		w.WriteHeader(400)
		return
	}

	stupidCache[effective] = time.Now()
}
