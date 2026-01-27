#!/usr/bin/env bash
set -eo pipefail

echo "Generating gogo proto code"

cd proto
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
rm -rf backend/pkg/pb && mkdir -p backend/pkg/pb
cp -r github.com/wutdahack1337/introductory-capstone/* backend/pkg/pb
rm -rf github.com
