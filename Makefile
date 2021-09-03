bridgeaddr: $(shell find . -name "*.go")
	go build -ldflags="-s -w" -o ./bridgeaddr

deploy: bridgeaddr
	ssh root@turgot 'systemctl stop bridgeaddr'
	scp bridgeaddr turgot:bridgeaddr/bridgeaddr
	ssh root@turgot 'systemctl start bridgeaddr'
