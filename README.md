# Requirements
### Glide
https://github.com/Masterminds/glide

### Go
https://golang.org/doc/install

# Setup
go get github.com/autonomousdotai/handshake-dispatcher

glide install

# Migrate db
go run migrate.go

# Run server
go run main.go
