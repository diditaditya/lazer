# Lazer

## Overview

This is mere a fun project of utter laziness. Essentially the app directly *reads* the mysql database table and it can then be requested through http. Nothing fancy. Just so the *backend* ready without writing any structs for each table, then go back lazing around wholeheartedly.

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

Start the app again, and you can access the app from your browser or whatever at `localhost:3500`. If you want to change the port just change the `docker-compose.yml` which maps default `gin` port at 8080 to 3500. Or just experiment with [traefik](https://docs.traefik.io/) to access it from your localhost with subdomain, which is cool and helpful, so you don't need to care about clashing ports.

## Features

Currently the features are developed really slowly. No test whatsoever, hence expect bugs, and when you find them please do fix them and create pull request :smiley:

The features are:
1. Automatically creates routes based on tables in the database. If you go to your browser and hit the root, you'll get the list of the tables which can be used as the routes or paths. For example if you have tables named `users` and `books`, they can be accessed through `http://localhost:3500/users` and `http://localhost:3500/books` assuming default setting for host and port.
2. The routes can filter the data by using querystring. For example if your `books` table has columns `id`, `title`, `author`, and `year`, you can filter the books by author like so `http://localhost:3500/books?author=nobody` or maybe the year also `http://localhost:3500/books?author=nobody&year=2000`. If you want multiple authors, use it like `author=nobody&author=somebody`. Please note this feature is still very basic and simplistic, manage your expectation.
3. The routes can also be used to `POST` data to create new entry. For example if you want to add new book, `POST` to `http://localhost:3500/books` with JSON body like `{"title": "abc", "author": 'me', "year": 1890}`, which returns nothing :satisfied: