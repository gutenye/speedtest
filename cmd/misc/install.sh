#!/usr/bin/env bash

set -x

NAME="freedom-speedtest"

sudo ln -sf `pwd`/$NAME /usr/bin
sudo cp *.{service,timer} /etc/systemd/system
