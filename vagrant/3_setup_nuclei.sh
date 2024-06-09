#!/bin/bash

sudo docker pull projectdiscovery/nuclei:v3.2.6
#Install naabu
sudo apt install -y libpcap-dev
go install -v github.com/projectdiscovery/naabu/v2/cmd/naabu@v2.3.1
sudo cp ~/go/bin/naabu /usr/bin/naabu