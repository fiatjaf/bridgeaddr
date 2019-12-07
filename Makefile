lnurl-tip: $(shell find . -name "*.go") bindata.go
	go build

public/bundle.js: $(shell find ./client)
	./node_modules/.bin/rollup -c rollup.config.js

bindata.go: public/bundle.js public/donate.html public/index.html public/global.css
	go-bindata -o bindata.go public/...
