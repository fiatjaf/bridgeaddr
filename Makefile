bridgeaddr: $(shell find . -name "*.go")
	CC=$$(which musl-gcc) go build -ldflags='-s -w -linkmode external -extldflags "-static"' -o ./bridgeaddr

deploy: bridgeaddr
	ssh root@cantillon 'systemctl stop bridgeaddr'
	scp bridgeaddr cantillon:bridgeaddr/bridgeaddr
	ssh root@cantillon 'systemctl start bridgeaddr'
