#!/usr/bin/env bash

NAME="freedom-speedtest-web"
set -x

sudo cp ${NAME}.service /etc/systemd/system
(cd programs/server && npm install --production)
sudo systemctl restart ${NAME}.service
