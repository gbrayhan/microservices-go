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

# Run container development

```bash
docker run --name microservice \
-v "$(pwd)":/app/microservices \
-e HOST_MONGO=mongo_host \
-e DATABASE_MONGO=database_mongo \ 
-e USER_DB=user_database_mysql \
-e PASSWORD_DB=password_mysql \
-e NAME_DB=name_db_mysql \
-e HOST_DB=host_mysql \
-e PORT_DB=port_mysql \
-d -p 8000:8080 \
${name_image}
```

