FROM golang:alpine as builder

RUN apk update && apk add git
COPY . $GOPATH/src/blackrez/openweathermap_exporter/
WORKDIR $GOPATH/src/blackrez/openweathermap_exporter/
RUN go get -d -v
RUN CGO_ENABLED=0 go build -o /go/bin/openweathermap_exporter

FROM scratch
EXPOSE 2112
COPY --from=builder /go/bin/openweathermap_exporter /openweathermap_exporter
ENTRYPOINT ["/openweathermap_exporter"]
