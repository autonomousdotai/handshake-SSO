# Setup
glide install

# Migrate db
go run migrate.go

# Run server
go run main.go
