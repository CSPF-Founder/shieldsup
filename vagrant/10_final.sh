#!/bin/bash

(sudo crontab -l 2>/dev/null; echo "0 * * * * /usr/bin/find /app/shieldsup/logs/scans/ -mindepth 1 -type f -mtime +10 -delete") | sudo crontab -


sudo rm -rf /app/build/scanner/
sudo rm -rf /app/build/scanner-api/
sudo rm -rf /app/build/manager/
sudo chown -R vagrant:vagrant /app/shieldsup/