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