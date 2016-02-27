all: wbt

wbt: cmd/wbt/*.go server/*.go jsonrpc/*.go config/*.go wlog/*.go
	gom build github.com/mercari/widebullet/cmd/wbt

gom:
	go get -u github.com/mattn/gom

bundle:
	go get -u golang.org/x/tools/cmd/goimports
	gom install

check:
	gom test ./...

fmt:
	go fmt ./...

imports:
	goimports -w .

clean:
	rm -rf wbt
