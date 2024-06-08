#!/bin/bash
#mkdir -p /app/shieldsup/infra/

cd /app/build/panel/
sudo make up

mkdir -p /app/shieldsup/panel/frontend/external


sudo cp -r /app/build/panel/panelfiles/frontend/external/* /app/shieldsup/panel/frontend/external


sudo chown -R vagrant:vagrant /app/shieldsup/