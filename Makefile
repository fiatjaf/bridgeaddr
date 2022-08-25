bridgeaddr: $(shell find . -name "*.go")
	CC=$$(which musl-gcc) go build -ldflags='-s -w -linkmode external -extldflags "-static"' -o ./bridgeaddr

deploy: bridgeaddr
	ssh root@turgot 'systemctl stop bridgeaddr'
	scp bridgeaddr turgot:bridgeaddr/bridgeaddr
	ssh root@turgot 'systemctl start bridgeaddr'
