#!/bin/bash

cd /app/build/scanner/

make build

cp /app/build/scanner/bin/scanner /app/shieldsup/scanner/scanner

chmod +x /app/shieldsup/scanner/scanner

cp /app/build/scanner/.env /app/shieldsup/scanner/.env

echo '#!/bin/bash' > /app/shieldsup/bin/scanner && echo 'cd /app/shieldsup/scanner/ && ./scanner $@' >> /app/shieldsup/bin/scanner

chmod +x /app/shieldsup/bin/scanner