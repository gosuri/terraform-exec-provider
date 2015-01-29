PROGRAM = terraform-provider-exec

test: install
	TF_ACC=1 TF_LOG=1 go test -v ""

build:
	go build -o bin/$(PROGRAM)

install: build
	cp bin/$(PROGRAM) $(GOPATH)/bin

