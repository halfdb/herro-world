BIN = $(PWD)/bin
PKG = $(PWD)/internal/pkg
APP = $(PWD)/internal/app
CMD = $(PWD)/cmd
SERVER_BIN = $(BIN)/herro-world
SQLBOILER_TOML = configs/sqlboiler.toml

models:
	sqlboiler --wipe --add-global-variants --add-soft-deletes --no-tests mysql -c $(SQLBOILER_TOML) -o $(PKG)/models

server: $(SERVER_BIN)

$(SERVER_BIN): $(APP)/server $(PKG) $(CMD)/server
	go build -v -o $(SERVER_BIN) $(CMD)/server