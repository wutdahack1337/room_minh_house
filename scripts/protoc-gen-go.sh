#!/usr/bin/env bash
set -eo pipefail

echo "Generating gogo proto code"

cd backend/proto
buf mod update
proto_dirs=$(find ./ -path -prune -o -name '*.proto' -print0 | xargs -0 -n1 dirname | sort | uniq)
for dir in $proto_dirs; do
  for file in $(find "${dir}" -maxdepth 1 -name '*.proto'); do
    if grep "option go_package" $file &> /dev/null ; then
      buf generate --template buf.gen.gogo.yml $file
    fi
  done
done

cd ..

# move proto files to the right places
rm -rf types/proto && mkdir -p types/proto
cp -r github.com/wutdahack1337/introductory-capstone/* types/proto
rm -rf github.com
