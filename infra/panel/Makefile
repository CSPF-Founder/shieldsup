DOCKER_STACK = shieldsup

setup: create_dirs cert
	@echo "First time setup done(dont call this agin)!"

create_dirs:
	@echo "First time setup..."
	mkdir -p /app/shieldsup/panel/certs/
	mkdir -p /app/shieldsup/panel/frontend/external/
	mkdir -p /app/shieldsup/data/temp_uploads/
	mkdir -p /app/shieldsup/data/work_dir/
	mkdir -p /app/shieldsup/bin/
	mkdir -p /app/shieldsup/manager/local_temp/
	mkdir -p /app/shieldsup/scanner/local_temp/
	mkdir -p /app/shieldsup/scanner_api/local_temp/
	mkdir -p /app/shieldsup/reporter/src
	mkdir -p /app/shieldsup/reporter/config
	mkdir -p /app/shieldsup/scanner_api/templates/
	mkdir -p /app/shieldsup/scanner_api/temp_templates/
	mkdir -p /app/shieldsup/data/reports/
	mkdir -p /app/shieldsup/logs/scans/


cert:
	@echo "Generating SSL Certificate..."
	openssl req -x509 -sha256 -nodes -days 4096 -newkey rsa:2048 -keyout /app/shieldsup/panel/certs/panel.key -out /app/shieldsup/panel/certs/panel.crt -subj "/C=US/ST=Cysecurity/L=Cysecurity/O=Dis/CN=shieldsup"
	@echo "SSL Generated"

up:
	@echo "Starting containers..."
	docker compose -p ${DOCKER_STACK} up -d
	@echo "Containers started!"

#up_build: starts all containers in the background and recreates them
up_build:
	@echo "Stopping containers..."
	docker compose -p ${DOCKER_STACK} down
	@echo "Building (when required) and Starting containers..."
	docker compose -p ${DOCKER_STACK} up -d --build
	@echo "Docker images build and containers started!"

# down: stops all containers
down:
	@echo "Stopping containers..."
	docker compose -p ${DOCKER_STACK} down
	@echo "Containers stopped!"
