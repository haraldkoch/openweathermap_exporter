FROM golang:alpine as builder

ARG TARGETPLATFORM
ARG BUILDPLATFORM
ARG VERSION

ENV CGO_ENABLED=0 \
    GOPATH=/go \
    GOBIN=/go/bin \
    GO111MODULE=on

WORKDIR /workspace

COPY . .

RUN \
  export GOOS \
  && GOOS=$(echo ${TARGETPLATFORM} | cut -d / -f1) \
  && export GOARCH \
  && GOARCH=$(echo ${TARGETPLATFORM} | cut -d / -f2) \
  && export GOARM \
  && GOARM=$(echo ${TARGETPLATFORM} | cut -d / -f3 | cut -c2-) \
  && go get -d -v \
  && go build -o /bin/prometheus-dnssec-exporter -ldflags="-w -s"

FROM scratch

ARG BUILD_DATE
ARG VCS_REF

EXPOSE 2112
COPY --from=builder /go/bin/openweathermap_exporter /openweathermap_exporter
ENTRYPOINT ["/openweathermap_exporter"]

LABEL maintainer="Harald Koch <harald.koch@gmail.com>" \
      org.opencontainers.image.created=${BUILD_DATE} \
      org.opencontainers.image.revision=${VCS_REF} \
      org.opencontainers.image.source="https://github.com/haraldkoch/openweathermap_exporter" \
      org.opencontainers.image.title="openweathermap_exporter" \
      org.opencontainers.image.version="${VERSION}"
