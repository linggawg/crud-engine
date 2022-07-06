## Prerequisite
- Docker Engine Latest
- Docker Compose Latest 
- Golang
- RDBMS (PostgreSQL, MariaDB)

## How to Run ?
- edit .env file to setting connection database, and place at root project path
```sh
DB_DIALECT=postgres
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=crud_go
DB_HOST=YOUR_IP_HOST/SERVER
DB_PORT=5432
```
- and run docker-compose up -d at root project path
```sh
docker-compose up -d
```
## SwaggerUI
- Generate Swagger Specification and SwaggerUI
```sh
$ swag init -g main.go --output docs
```
- Documentation API in path
```sh
{{base_url}}/swagger/index.html
```
<img src="https://resources.jetbrains.com/storage/products/company/brand/logos/GoLand.svg" alt="GoLand logo.">
