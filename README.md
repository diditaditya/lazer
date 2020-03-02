# Lazer

## Overview

Too lazy to write overview lol

## Development

Make sure [docker](https://docs.docker.com/install/) and [docker-compose](https://docs.docker.com/compose/install/) are installed.

To start development, clone the repo and start the docker.

```bash
$ git clone https://github.com/diditaditya/lazer
$ cd lazer
$ docker-compose up -d
```

This will create 2 containers, the go and the mysql. The app will not automatically run. You must start it manually from inside the go container.

```bash
$ docker exec -it go bash
$ cd src/lazer
$ go run main.go
```

Which will `panic` because you need `.env` file in your root. The required variables are:

```bash
DB=mysql
DB_HOST=thedbhost
DB_USER=thedbuser
DB_PASSWORD=thedbpassword
DB_NAME=thedbname
```