export GOPATH := /home/nox73/.go:$(PWD)

run:
	@( go run *.go )

deps:
	@( \
		go get github.com/NOX73/go-neural; \
		go get github.com/NOX73/go-neural/persist; \
		go get github.com/NOX73/go-neural/lern; \
	)
