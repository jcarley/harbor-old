#!/usr/bin/env bash

curl -sL https://deb.nodesource.com/setup_0.12 | sudo bash - \
  && sudo apt-get update \
  && sudo apt-get install -y --no-install-recommends nodejs \
  && sudo apt-get clean \
  && sudo rm -rf /var/lib/apt/lists/* /tmp/* /var/tmp/* \
  && sudo npm install npm -g

sudo npm install -g yo

