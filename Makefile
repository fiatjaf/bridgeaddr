lightningaddr: $(shell find . -name "*.go")
	go build -ldflags="-s -w" -o ./lightningaddr

deploy: lightningaddr
	ssh root@hulsmann 'systemctl stop lightningaddr'
	scp lightningaddr hulsmann:lightningaddr/lightningaddr
	ssh root@hulsmann 'systemctl start lightningaddr'
