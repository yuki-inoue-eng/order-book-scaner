.DEFAULT_GOAL := help

.PHONY: count-go
count-go: ## Count number of lines of all go codes.
	find . -name "*.go" -type f | xargs wc -l | tail -n 1

.PHONY: setup
setup: ## Resolve library dependencies with Go Modules
	go mod download

.PHONY: build
build: setup  ## build application
	go build -o build/order-book-searcher .

# See "Self-Documented Makefile" article
# https://marmelab.com/blog/2016/02/29/auto-documented-makefile.html
.PHONY: help
help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'
