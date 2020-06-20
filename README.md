# gin_server

gin_server is a web server implemented using [Gin](https://github.com/gin-gonic/gin). I will find a better name soon. It serves as a backend for [go project frontend](https://github.com/iamkahvi/go_project_frontend)


## Installation
1. Clone repository
2. Create database with mysql
3. Connect gorm to database with `gorm.Open()` [here](storage/main.go)

## Usage
In the project directory, run
```bash
go run main.go
```

## API Documentation
Located in [docs.md](docs.md)

## Built With
- [gin](https://github.com/gin-gonic/gin)
- [gorm](https://gorm.io/docs/)
- [mysql](https://dev.mysql.com/doc/refman/8.0/en/introduction.html)

## Config
- Locked to `http://google.com` and `http://localhost:3000` origins

## Inspiration

- https://github.com/burxtx/gin-microservice-boilerplate
- https://github.com/go-kit/kit
- https://github.com/demo-apps/go-gin-app
