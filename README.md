# crud-engine
Crud Engine using Go Lang

## Prerequisite
- Docker Engine Latest
- Docker Compose Latest 
- Golang
- RDBMS (PostgreSQL, MySQL)

## How to Run ?
- edit .env file to setting connection database, and place at root project path
```sh
# for postgresql connection
DB_DIALECT=postgres
DB_USER=postgres
DB_PASSWORD=root
DB_NAME=crud_go
DB_HOST=YOUR_IP_HOST/SERVER
DB_PORT=5432

# for mysql connection
DB_DIALECT=mysql
DB_USER=root
DB_PASSWORD=
DB_NAME=chat_app
DB_HOST=YOUR_IP_HOST/SERVER
DB_PORT=3306
```
- and run docker-compose up -d at root project path
```sh
docker-compose up -d
```
## SwaggerUI?
- Generate Swagger Specification and SwaggerUI
```sh
$ swag init -g main.go --output docs
```
- Documentation API in path
```sh
{{base_url}}/swagger/index.html
```
