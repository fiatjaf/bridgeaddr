module github.com/fiatjaf/lnurl-tip

go 1.13

require (
	github.com/elazarl/go-bindata-assetfs v1.0.0
	github.com/fiatjaf/go-lnurl v1.0.0
	github.com/fiatjaf/lightningd-gjson-rpc v0.1.1-0.20191204225807-4e73275bc053
	github.com/fiatjaf/ln-decodepay v0.0.0-20191204194730-0355a2d6e26e
	github.com/go-bindata/go-bindata v3.1.2+incompatible // indirect
	github.com/gorilla/mux v1.7.3
	github.com/jmoiron/sqlx v1.2.0
	github.com/kelseyhightower/envconfig v1.4.0
	github.com/lib/pq v1.2.0
	github.com/lightningnetwork/lnd v0.8.0-beta-rc3.0.20191206010316-d230cf89b94d // indirect
	github.com/lucsky/cuid v1.0.2
	github.com/rs/zerolog v1.17.2
	github.com/stretchr/testify v1.4.0 // indirect
	github.com/tidwall/gjson v1.6.0
	github.com/tidwall/pretty v1.0.1 // indirect
	github.com/tidwall/sjson v1.0.4
	golang.org/x/crypto v0.0.0-20191205180655-e7c4368fe9dd // indirect
)

replace github.com/fiatjaf/go-lnurl => /home/fiatjaf/comp/go-lnurl
