#!/bin/bash
cd /app/build/manager/
make build
mkdir -p /app/shieldsup/manager/
mkdir -p /app/shieldsup/manager/local_temp/

cp /app/build/manager/bin/manager /app/shieldsup/manager/manager
chmod +x /app/shieldsup/manager/manager

cp /app/build/manager/.env /app/shieldsup/manager/.env


cp /app/build/manager/shieldsup-manager.service /etc/systemd/system/shieldsup-manager.service
systemctl daemon-reload
systemctl enable shieldsup-manager.service
systemctl start shieldsup-manager.service

