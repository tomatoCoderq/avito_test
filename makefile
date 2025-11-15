# Makefile to control project via docker-compose

COMPOSE ?= docker-compose
COMPOSE_FILE ?= docker-compose.yml
SERVICE ?= web

.PHONY: up restart recreate clean lint-fix

# Start services (detached)
up:
	$(COMPOSE) -f $(COMPOSE_FILE) up -d

# Restart all services  
restart:
	$(COMPOSE) -f $(COMPOSE_FILE) restart

# Recreate service (SERVICE=... optional)
recreate:
	$(COMPOSE) -f $(COMPOSE_FILE) up -d --build --force-recreate $(SERVICE)

# Remove containers, images and volumes (destructive)
clean:
	$(COMPOSE) -f $(COMPOSE_FILE) down --rmi all --volumes --remove-orphans

# Run linter with auto-fix
lint-fix:
	@which golangci-lint > /dev/null || (echo "golangci-lint not found. Install: brew install golangci-lint" && exit 1)
	golangci-lint run --fix ./...