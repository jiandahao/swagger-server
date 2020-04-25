#!/usr/bin/env sh
set -x
for file in `find ./swagger -type f -name '*.swagger.json' | xargs echo`
do
	swagger2openapi "$file" -o "$file"
done