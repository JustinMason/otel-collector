FROM golang:1.22-bookworm

RUN apt-get update && apt-get clean


COPY . /go/src/github.com/justinmason/otel-collector
WORKDIR /go/src/github.com/justinmason/otel-collector

RUN go build && chmod +x otel-collector

FROM golang:1.22-bookworm
WORKDIR /app

COPY --from=0 /go/src/github.com/justinmason/otel-collector/otel-collector otelcol

ENTRYPOINT ["./otelcol"]