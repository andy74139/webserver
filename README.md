Web Server
======

An implementation of a web server with user and authentication management.

## Usage

Prerequisite: please install following software to build and run backend system

* git: https://git-scm.com/downloads
* taskfile: https://taskfile.dev/installation/
* docker: https://docs.docker.com/engine/install/

```bash
# Download, find a folder to put source code
git clone https://github.com/andy74139/webserver
cd webserver

# IMPORTANT: Please ensure docker daemon is running
# start server
task start

# stop server
task stop

# stop and remove server, it cleans DB data.
task remove

## Development and document usages

# rebuild server, for development testing
task rebuild

# start/stop API document server (read it on localhost:8090)
task start-docs
task stop-docs

# update document, it can run when API document server is running
task update-docs
```

## Technical Stack

* Task and build management tool: [taskfile](https://taskfile.dev/)
* Infra-container
    * [Docker](https://www.docker.com/)
    * Kubernetes (in progress)
* Backend: [Go](https://go.dev/)
    * HTTP web server: [gin](https://github.com/gin-gonic/gin)
        * RESTful API
        * Authentication: SSO (Single Sign-On)
        * Authorization: JWT (JSON Web Token)
    * Database ORM: [bun](https://bun.uptrace.dev/)
    * Testing (in progress)
        * simple test: [testify](https://pkg.go.dev/github.com/stretchr/testify)
        * structured BDD test: [ginkgo/gomega](https://onsi.github.io/ginkgo/)
        * e2e test: [httpexpect](https://github.com/gavv/httpexpect)
    * Logging: [zap](https://github.com/uber-go/zap/)
    * Documentation: [swagger](https://swagger.io/), [go-swagger3](https://github.com/parvez3019/go-swagger3), [swagger-ui](https://github.com/swagger-api/swagger-ui/blob/HEAD/docs/usage/installation.md)
    * Code Structure: refer to [DDD](https://github.com/sklinkert/go-ddd)
* Database: Postgresql
    * pgbouncer  (in progress)
* Cache: Redis

