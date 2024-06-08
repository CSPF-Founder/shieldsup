#!/bin/bash

mkdir -p /app/shieldsup/scanner_api/
mkdir -p /app/shieldsup/scanner_api/local_temp/

cd /app/build/scanner-api/
make build


cp /app/build/scanner-api/bin/api /app/shieldsup/scanner_api/api
chmod +x /app/shieldsup/scanner_api/api

cp /app/build/scanner-api/.env /app/shieldsup/scanner_api/.env

sudo cp /app/build/scanner-api/scanner-api.service /etc/systemd/system/scanner-api.service
sudo systemctl daemon-reload
sudo systemctl enable scanner-api.service
sudo systemctl start scanner-api.service