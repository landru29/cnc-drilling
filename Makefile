.PHONY: cnc-drilling

cnc-drilling:
	go build -o $@ ./cmd/...
