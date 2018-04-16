test:
		go test $$(go list ./... | grep -v /vendor/)

fmt:
		go fmt $$(go list ./... | grep -v /vendor/)

.PHONY: test fmt