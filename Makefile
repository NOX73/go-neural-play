export GOPATH := $(GOPATH):$(PWD)

run:
	@( go run main.go )

deps:
	@( \
		go get github.com/NOX73/go-neural; \
		go get github.com/NOX73/go-neural/persist; \
		go get github.com/NOX73/go-neural/lern; \
	)

vim: 
	vim .
