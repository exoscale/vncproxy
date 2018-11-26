.PHONY: proxy
proxy:
	@go build -o bin/proxy ./proxy/cmd/

clean:
	@rm -rf bin/
