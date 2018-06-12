#!/usr/bin/env bash

mkdir -p dist
cp install.sh dist

cd cmd
./ake build
cp dist/freedom-speedtest ../dist
cp misc/* ../dist
cd -

cd web
./ake build
cp misc/* ../dist
cd -

rsync -aP  dist/ root@box:/app/freedom-speedtest/
ssh root@box cd /app/freedom-speedtest && ./install.sh
