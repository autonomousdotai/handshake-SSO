# Requirements
### Glide
https://github.com/Masterminds/glide

### Go
https://golang.org/doc/install

# Setup
```bash
go get github.com/autonomousdotai/handshake-dispatcher
cd /path/to/handshake-dispatcher
glide install
```

# Configure
```bash
cd /path/to/handshake-dispatcher
cp config/conf.yaml.default config/conf.yaml
```

Edit `config/conf.yaml` to fix your config

# Migrate db
Mysql. create database if not exists

`CREATE DATABASE database CHARACTER SET utf8 COLLATE utf8_general_ci;`

```bash
go run migrate.go
```

# Run server
```bash
go run main.go
```
