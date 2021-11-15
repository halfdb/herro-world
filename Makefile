BIN = $(PWD)/bin
PKG = $(PWD)/internal/pkg
APP = $(PWD)/internal/app
CMD = $(PWD)/cmd
SCRIPTS = $(PWD)/scripts
SERVER_BIN = $(BIN)/herro-world
SQLBOILER_TOML = configs/sqlboiler.toml

install-tools:
	$(SCRIPTS)/install-tools.sh

models:
	sqlboiler --wipe --add-global-variants --add-soft-deletes --no-tests mysql -c $(SQLBOILER_TOML) -o $(PKG)/models

db-up:
	goose -v -dir $(PWD)/deployments/database mysql $$DB_STRING up

db-down:
	goose -v -dir $(PWD)/deployments/database mysql $$DB_STRING down

server:
	go build -v -o $(SERVER_BIN) $(CMD)/server