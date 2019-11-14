IMAGE_NAME := remmelt/openweathermap_exporter:0.0.1

.PHONY: build
build: openweathermap_exporter

openweathermap_exporter: get
	@go build

.PHONY: clean
clean:
	@go clean

.PHONY: get
get:
	@go get -d -v

.PHONY: docker-image
docker-image:
	@docker build -t ${IMAGE_NAME} .

.PHONY: push-image
push-image:
	@docker push ${IMAGE_NAME}
