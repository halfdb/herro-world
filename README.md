# Preparation

## Set up database access

1. Set up the `DB_STRING` environment variable. It will be used in the program and in goose.
```shell
export DB_STRING=user:pass@tcp\(localhost:3306\)/db?parseTime=true
```

2. Set up the `sqlboiler.toml`. See `configs/sqlboiler.example` for example.

## Install tools

* sqlboiler v4.6.0
* goose v3.3.1
