.PHONY: tidy
tidy:
	go mod tidy

.PHONY: deps
deps: tidy
	go mod vendor
