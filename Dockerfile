FROM golang:1.16.5-buster

RUN apt-get update
RUN apt-get install -y build-essential

RUN curl -fsSL https://deb.nodesource.com/setup_lts.x | bash -
RUN apt-get install -y nodejs

RUN npm install -g nodemon

WORKDIR /go/src/lazer
COPY ./src .

COPY ./nodemon.json /

COPY ./scripts/wait-for-it.sh /
RUN chmod +x /wait-for-it.sh