TARGETS_NOVENDOR=$(shell glide novendor)

all: wbt

wbt: cmd/wbt/*.go server/*.go jsonrpc/*.go config/*.go wlog/*.go
	GO15VENDOREXPERIMENT=1 go build cmd/wbt/wbt.go

bundle:
	glide install

check:
	GO15VENDOREXPERIMENT=1 go test $(TARGETS_NOVENDOR)

fmt:
	@echo $(TARGETS_NOVENDOR) | xargs go fmt

clean:
	rm -rf wbt
