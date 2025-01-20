.DEFAULT_GOAL = help

## Display this help message
help:
	@cat $(MAKEFILE_LIST) | awk '							\
	BEGIN { FS = ":"; printf "Usage: make \033[36m<target>\033[0m\n\n" }		\
	/^## / {									\
		sub(/##[ \t]*/, "", $$0) ;						\
		help = help (help ? "\n\t" : "") $$0 ;					\
		next ;									\
	}										\
	/^[^ \t]+:/ {									\
		if (help) printf ("  \033[36m" "%-25s" "\033[0m-- " help "\n", $$1 )	\
	}										\
	/.*/ { help = "" }'

## Run go fmt
fmt:
	goimports -w -local github.com/doom $$(find . -type f -name '*.go' -not -path "./vendor/*")
	go fmt ./...

## Run golangci-lint
lint:
	golangci-lint run

## Run tests
test:
	go test -v -race -cover -coverprofile coverage.txt ./...

.PHONY: help fmt lint test
