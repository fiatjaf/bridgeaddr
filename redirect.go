package main

import (
	"bytes"
	"errors"
	"fmt"
	"net"
	"net/http"
	"regexp"
	"strings"
)

func redirect(w http.ResponseWriter, r *http.Request) {
	domain := r.Host
	if v, err := net.LookupTXT("_redirect." + domain); err == nil && len(v) > 0 {
		if redirect, err := getRedirect(v, r.URL.String()); err == nil {
			http.Redirect(w, r, redirect.Location, redirect.Status)
		}
	}

	fmt.Fprintf(w, "hosted by "+s.ServiceURL)
}

type redirectConfig struct {
	From          string
	To            string
	RedirectState string
}

type redirectValue struct {
	Location string
	Status   int
}

var configRE = regexp.MustCompile(`Redirects?(\s+.*)`)
var fromRE = regexp.MustCompile(`\s+from\s+(/\S*)`)
var toRE = regexp.MustCompile(`\s+to\s+(https?\://\S+|/\S*)`)
var stateRE = regexp.MustCompile(`\s+(permanently|temporarily)|\s+with\s+(301|302|307|308)`)

func getRedirect(txt []string, url string) (*redirectValue, error) {
	var catchAlls []*redirectConfig
	for _, record := range txt {
		config := parseRedirect(record)
		if config.From == "" {
			catchAlls = append(catchAlls, config)
			continue
		}
		redirect := translateRedirect(url, config)
		if redirect != nil {
			return redirect, nil
		}
	}

	var config *redirectConfig
	for _, config = range catchAlls {
		redirect := translateRedirect(url, config)
		if redirect != nil {
			return redirect, nil
		}
	}

	return nil, errors.New("No paths matched")
}

func parseRedirect(record string) *redirectConfig {
	configMatches := configRE.FindStringSubmatch(record)
	if len(configMatches) == 0 {
		return nil
	}

	fromMatches := fromRE.FindStringSubmatch(configMatches[1])
	toMatches := toRE.FindStringSubmatch(configMatches[1])
	stateMatches := stateRE.FindStringSubmatch(configMatches[1])

	config := new(redirectConfig)
	if len(fromMatches) > 0 {
		config.From = fromMatches[1]
	}
	if len(toMatches) > 0 {
		config.To = toMatches[1]
	}
	if len(stateMatches) > 0 {
		config.RedirectState = stateMatches[1]
		if config.RedirectState == "" {
			config.RedirectState = stateMatches[2]
		}
	}

	return config
}

func translateRedirect(uri string, config *redirectConfig) *redirectValue {
	if uri == "" {
		return nil
	}
	if config == nil {
		return nil
	}
	if config.To == "" {
		return nil
	}

	redirect := &redirectValue{Location: config.To}

	switch config.RedirectState {
	case "301", "permanently":
		redirect.Status = 301
	case "302", "temporarily":
		redirect.Status = 302
	case "307":
		redirect.Status = 307
	case "308":
		redirect.Status = 308
	default:
		redirect.Status = 302
	}

	// no `From` assumes catch-all, so redirect immediately to `Location`
	if config.From == "" {
		return redirect
	}

	count := strings.Count(config.From, `*`)

	var exp bytes.Buffer
	exp.WriteString(`^`)
	exp.WriteString(strings.Replace(regexp.QuoteMeta(config.From), `\*`, `(.*)`, 1))
	exp.WriteString(`$`)

	fromRE := regexp.MustCompile(exp.String())

	// if we can't find the pattern, return to continue to next record
	if !fromRE.MatchString(uri) {
		return nil
	}

	// wildcard replacement of `uri` if there's a wildcard in our `From` path
	if count > 0 {
		redirect.Location = fromRE.ReplaceAllString(uri, strings.Replace(redirect.Location, `*`, `${1}`, 1))
	}

	return redirect
}
