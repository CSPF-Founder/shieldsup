#!/bin/bash
export DEBIAN_FRONTEND=noninteractive
sudo apt-get update -y 
sudo apt install ntp resolvconf -y
sudo systemctl start ntp
sleep 5
mkdir -p /etc/resolvconf/resolv.conf.d/
sudo touch /etc/resolvconf/resolv.conf.d/base
sudo echo "nameserver 8.8.8.8" >> /etc/resolvconf/resolv.conf.d/base
sudo echo "nameserver 8.8.4.4" >> /etc/resolvconf/resolv.conf.d/base
resolvconf -u
sudo resolvconf -u
sudo apt-get install --fix-missing
sleep 5
sudo apt install make net-tools sudo wget nano telnet acl python3 python3-poetry python3-venv -y
sleep 5

#Makedirs

sudo mkdir -p /app/shieldsup/bin/
sudo mkdir -p /app/shieldsup/panel/

sudo mkdir -p /app/shieldsup/scanner_api/templates/ 
sudo mkdir -p /app/shieldsup/scanner_api/temp_templates/ 

sudo chown -R vagrant:vagrant /app/shieldsup/

#Install Golang
cd ~
sleep 5
wget https://go.dev/dl/go1.22.3.linux-amd64.tar.gz
sudo rm -rf /usr/local/go && tar -C /usr/local -xzf go1.22.3.linux-amd64.tar.gz
sudo echo 'export PATH=$PATH:/usr/local/go/bin' >> /etc/profile
source /etc/profile