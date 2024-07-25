#!/bin/bash

sudo docker pull projectdiscovery/nuclei:v3.3.0
#Install naabu
sudo apt install -y libpcap-dev
go install -v github.com/projectdiscovery/naabu/v2/cmd/naabu@v2.3.1
sudo cp ~/go/bin/naabu /usr/bin/naabu
#Install katana
go install github.com/projectdiscovery/katana/cmd/katana@v1.1.0
sudo cp ~/go/bin/katana /usr/bin/katana