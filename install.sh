#!/usr/bin/env bash

cp *.{service,timer} /etc/systemd/system
cp freedom-speedtest /usr/bin

(cd bundle/programs/server && npm install --production)
