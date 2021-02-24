# Golang Microservices Boilerplate
[![issues](https://img.shields.io/github/issues/gbrayhan/microservices-go)](https://github.com/gbrayhan/microservices-go/tree/master/.github/ISSUE_TEMPLATE)
[![forks](https://img.shields.io/github/forks/gbrayhan/microservices-go)](https://github.com/gbrayhan/microservices-go/network/members)
[![stars](https://img.shields.io/github/stars/gbrayhan/microservices-go)](https://github.com/gbrayhan/microservices-go/stargazers)
[![license](https://img.shields.io/github/license/gbrayhan/microservices-go)](https://github.com/gbrayhan/microservices-go/tree/master/LICENSE)
[![CodeFactor](https://www.codefactor.io/repository/github/gbrayhan/microservices-go/badge/master)](https://www.codefactor.io/repository/github/gbrayhan/microservices-go/overview/master)

Example structure to start a microservices project with golang. Using a MySQL database.


# Build image docker development
```bash
docker build -t ${name_image} --force-rm .
```

# Swagger Implementation
```bash
swag init -g routes/ApplicationV1.go
```


# Unit test command
```bash
# run recursive test
go test  ./test/unitgo/...
# clean go test results in cache
go clean -testcache
```



