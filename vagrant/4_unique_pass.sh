#!/bin/bash
mkdir -p /app/build/panel/
mkdir -p /app/build/scanner/
mkdir -p /app/shieldsup/reporter/config
mkdir -p /app/build/scanner-api/
mkdir -p /app/build/manager/


cp -a /vagrant/infra/panel/. /app/build/panel/
cp -a /vagrant/code/scanner/. /app/build/scanner/
cp -r /vagrant/code/reporter/config/. /app/shieldsup/reporter/config
cp -a /vagrant/code/scanner-api/. /app/build/scanner-api/
cp -a /vagrant/code/manager/. /app/build/manager/


#Edit and make SSL panel
cd /app/build/panel/

rootpass=`(sudo head /dev/urandom | tr -dc 'A-Za-z0-9' | head -c 20)`
normalpass=`(sudo head /dev/urandom | tr -dc 'A-Za-z0-9' | head -c 20)`

sed -i "s/\[ROOT_PASS_TO_REPLACE\]/$rootpass/g" /app/build/panel/docker-compose.yml
sed -i "s/\[PASSWORD_TO_REPLACE\]/$normalpass/g" /app/build/panel/docker-compose.yml


sudo make


randomapikey=`(sudo head /dev/urandom | tr -dc 'A-Za-z0-9' | head -c 40)`
sed -i "s/\[ROOT_PASS_TO_REPLACE\]/$rootpass/g" /app/build/scanner/.env
sed -i "s/\[REPLACE_API_KEY\]/$randomapikey/g" /app/build/scanner/.env



#Edit reporter config
mkdir -p /app/shieldsup/reporter/config
sed -i "s/\[ROOT_PASS_TO_REPLACE\]/$rootpass/g" /app/shieldsup/reporter/config/app.conf

#Edit scanner api config
sed -i "s/\[REPLACE_API_KEY\]/$randomapikey/g" /app/build/scanner-api/.env
#Edit manager config
sed -i "s/\[ROOT_PASS_TO_REPLACE\]/$rootpass/g" /app/build/manager/.env
#chown back

sudo chown -R vagrant:vagrant /app/shieldsup/