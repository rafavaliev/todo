# Todo list with search feature

This is a simple todo list with search feature.

# Project structure

* cmd - main package that initilized the app and starts http server
* internal - internal packages of the project
* task - task package with task model and task repository, service
* search - search package with search model and search service
* user - user package with user model and user repository
* handler - handlers for http requests

# How to run

* docker-compose up -d
* Run http_example.http file in your IDE
    * Or run some curl commands:

```bash
    curl -X POST -d '{"username": "rafael5","password":"test"}' "localhost:80/v1/signup"
    curl -u rafael5:test localhost:80/v1/tasks 
    curl -X POST -d '{"title": "test","description":"test"}' -u rafael5:test "localhost:80/v1/tasks"
    curl -X POST -d '{"title": "test2","description":"test2"}' -u rafael5:test "localhost:80/v1/tasks"
    curl -u rafael5:test "localhost:80/v1/search?query=test"
```

# How to run tests

* make tests

# What's done

1. Simple RESTful API for managing tasks and search over them
2. Basic authentication(only signup and basic auth)
3. Search over tasks using inverted index
4. Dockerfile and docker-compose
5. Some unit tests and integration tests for http handler
6. Prometheus metrics

# What can be improved?

1. HTTPS :)
2. tests
4. SQL queries instead of ORM
5. logging
6. error handling