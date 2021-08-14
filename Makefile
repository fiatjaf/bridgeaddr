bridgeaddr: $(shell find . -name "*.go")
	go build -ldflags="-s -w" -o ./bridgeaddr

deploy: bridgeaddr
	ssh root@hulsmann 'systemctl stop bridgeaddr'
	scp bridgeaddr hulsmann:bridgeaddr/bridgeaddr
	ssh root@hulsmann 'systemctl start bridgeaddr'
