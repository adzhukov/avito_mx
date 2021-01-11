.DEFAULT_GOAL = help

PROJECT = avito_mx

.PHONY: build
build: clone ## Build docker images
	docker-compose build

.PHONY: up
up: clone ## Run containers
	docker-compose up

.PHONY: pg_shell
shell_pg: ## Connect to a shell in a postgres container
	@docker exec -it pg_$(PROJECT) bash

.PHONY: shell_app
shell_app: ## Connect to a shell in a app container
	@docker exec -it $(PROJECT) bash

.PHONY: psql
psql: ## Connect to psql in a postgres container
	@docker exec -it pg_$(PROJECT) psql -U postgres -d $(PROJECT)

.PHONY: down
down: ## Stop containers
	docker-compose down

.PHONY: wipe_db
wipe_db: ## Wipe and restart db
	-docker-compose down postgres
	-rm -rf .postgres
	docker-compose up postgres

.PHONY: help
help: ## Display this
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) \
	  | sort \
	  | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[0;32m%-30s\033[0m %s\n", $$1, $$2}'
