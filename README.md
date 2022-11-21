# Go RESTful API Starter Kit (Boilerplate) Inspired by [go-rest-api](https://github.com/qiangxue/go-rest-api)

An idiomatic Go REST API starter kit (boilerplate) following the SOLID principles and Clean Architecture

This starter kit is designed to get you up and running with a project structure optimized for developing
RESTful API services in Go. It promotes the best practices that follow the [SOLID principles](https://en.wikipedia.org/wiki/SOLID)
and [clean architecture](https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html). 
It encourages writing clean and idiomatic Go code. 

The kit provides the following features right out of the box:

* RESTful endpoints in the widely accepted format
* Standard CRUD operations of a database table
* JWT-based authentication
* Environment dependent application configuration management
* Structured logging with contextual information
* Error handling with proper error response generation
* Database migration
* Data validation
* Full test coverage
* Live reloading during development
 
The kit uses the following Go packages which can be easily replaced with your own favorite ones
since their usages are mostly localized and abstracted. 

* Routing: [echo](https://github.com/labstack/echo)
* Database access: [sqlc](https://github.com/kyleconroy/sqlc)
* Database migration: [golang-migrate](https://github.com/golang-migrate/migrate)
* Data validation: [validator](https://github.com/go-playground/validator)
* Logging: [zap](https://github.com/uber-go/zap)
* JWT: [jwt-go](https://github.com/golang-jwt/jwt)
