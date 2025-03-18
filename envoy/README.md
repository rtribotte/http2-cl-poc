# Run with docker
# Unfortunatelly, this does not work

docker run --rm -p 8080:8080 \
-v ./envoy.yaml:/etc/envoy/envoy.yaml \
-v ./../certs:/etc/certs \
envoyproxy/envoy:v1.28.0 --log-level debug --config-path /etc/envoy/envoy.yaml

# Run with docker-compose