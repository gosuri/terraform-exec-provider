PROGRAM = terraform-provider-exec

test: build
	TF_ACC=1 TF_LOG=1 go test -v ""

updatedeps:
	@go get -u golang.org/x/tools/cmd/stringer
	@go get -f -u -v ./...

build:
	go build -o bin/$(PROGRAM)

install: updatedeps build
	cp bin/$(PROGRAM) $(GOPATH)/bin

.PHONEY: test updatedeps build install
