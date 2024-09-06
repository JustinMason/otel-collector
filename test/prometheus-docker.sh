docker run -d \
  --name prometheus \
  -p 9090:9090 \
  -v /Users/justinmason/repos/justinmason/otel-collector/test/prometheus.yml:/etc/prometheus/prometheus.yml \
  prom/prometheus