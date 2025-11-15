# Makefile to control project via docker-compose

COMPOSE ?= docker-compose
SERVICE ?= api

.PHONY: up restart recreate clean lint-fix

# Start services (detached)
up:
	$(COMPOSE) up -d

# Restart all services  
restart:
	$(COMPOSE) restart

# Recreate service (SERVICE=... optional)
recreate:
	$(COMPOSE) up -d --build --force-recreate $(SERVICE)

# Remove containers, images and volumes (destructive)
clean:
	$(COMPOSE) down --rmi all --volumes --remove-orphans

# Run linter with auto-fix
lint-fix:
	@which golangci-lint > /dev/null || (echo "golangci-lint not found. Install: brew install golangci-lint" && exit 1)
	golangci-lint run --fix ./...