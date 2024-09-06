# Justin's Experimental Open Telemetry Collector

Custom OTel Collector component, see [OTel Collector Builder](https://github.com/open-telemetry/opentelemetry-collector/blob/main/cmd/builder/README.md)
for more information.

## Receivers

This collector provides examples of a prometheus remote write endpoint and
json log endpoint.

## Debug

```bash
go install github.com/go-delve/delve/cmd/dlv@latest
~/go/bin/dlv --listen=:2345 --headless=true --api-version=2 --accept-multiclient --log exec otel-collector -- --config otelcol.yaml
```
