.PHONY: build
build: openweathermap-exporter

openweathermap-exporter: get
	@go build

.PHONY: clean
clean:
	@go clean

.PHONY: get
get:
	@go get -d -v
