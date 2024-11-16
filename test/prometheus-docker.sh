docker run -d \
  --name prometheus \
  --net host \
  -p 9091:9091 \
  -v /Users/justin.mason/repos/otel-collector/test/prometheus.yml:/etc/prometheus/prometheus.yml \
  prom/prometheus