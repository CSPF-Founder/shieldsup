#!/bin/bash

mkdir -p /app/shieldsup/reporter/src

cp -r /vagrant/code/reporter/src/. /app/shieldsup/reporter/src


cd /app/shieldsup/reporter/src
python3 -m venv venv
source venv/bin/activate
pip install poetry
poetry install --no-interaction --no-root
deactivate

echo '#!/bin/bash' > /app/shieldsup/bin/reporter && echo 'cd /app/shieldsup/reporter/src && source venv/bin/activate && python3 cli.py $@' >> /app/shieldsup/bin/reporter

chmod +x /app/shieldsup/bin/reporter