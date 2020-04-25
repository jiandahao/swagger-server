docker run --rm -v $(pwd):$(pwd) -w $(pwd) nanxi/protoc:v1.0.0 --go_out=plugins=grpc:. \
--grpc-gateway_out=logtostderr=false:.  -I. *.proto