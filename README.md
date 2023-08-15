# Test Go Gin Gonic Project

This project is intended to accomplish [Depoguna Bangunan Online](https://dbo.id/id) test.

## Tech Stack
- Golang v1.21
- Gin Gonic v1.9.1
- PostgreSQL

## Tools
- Docker
- Docker Compose
- Git
- Postman

## How to Run the Project
### Prerequisites
1. Docker has been installed.
2. Docker Compose has been installed.
3. Golang has been installed.
4. Git has been installed.

### Steps
1. Go to the project folder
```shell
cd go/to/your/project/path
```
2. Clone this project
```shell
git clone github.com/willywartono14/test-go-gin-gonic.git
```
3. Change directory to `test-go-gin-gonic`
```shell
cd test-go-gin-gonic
```
4. Execute `docker-compose` and wait for both `postgres` and `api` containers are successfully ran
```shell
docker-compose up --build
```
5. Execute migration
```shell
go run migration/main.go migrate-up
```

## Features
### Customer Management
1. **GET** `/api/customers`
2. **GET** `/api/customers/:id`
3. **PUT** `/api/customers/:id`
4. **DELETE** `/api/customers/:id`

### Order Management
1. **GET** `/api/orders`
2. **GET** `/api/orders/:id`
3. **POST** `/api/orders`
4. **PUT** `/api/orders/:id`
5. **DELETE** `/api/orders/:id`

### Item Management
1. **GET** `/api/items`

### Authentication
1. **POST** `/api/login`
2. **POST** `/api/register`