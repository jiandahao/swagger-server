#!/bin/bash
set -e
cd $(dirname $0)

mkdir -p ./swagger
export OUTPUT=./swagger

find ./ -type f -name '*.proto'  -exec \
docker run --rm -u $(id -u ${USER}):$(id -g ${USER}) \
-v $(pwd):$(pwd) -w $(pwd) nanxi/protoc:v1.0.0 \
--swagger_out=logtostderr=true:./swagger -I. {} \;

docker run --rm -v $(pwd):$(pwd) -w $(pwd) mermade/swagger2openapi:latest sh swagger2openapi.sh

export HOST="http://47.244.202.189:8080"
#export HOST="http://127.0.0.1:8088"
for file in `find ./swagger -type f -name '*.swagger.json' | xargs echo`
do
  curl -X POST \
  -H "Content-Type: multipart/form-data" \
  -F "file=@${file}" \
 ${HOST}/new
#  docker run --rm -v $(pwd):$(pwd) -w $(pwd) jiandahao/swagger-gen-cli:v1.0.0 --namespace=xlppc --server-host=127.0.0.1:8088
done