FROM golang:1.27-trixie

RUN apt-get update
RUN apt-get install -y build-essential

WORKDIR /go/src/lazer
COPY . .

RUN go install github.com/air-verse/air@latest

COPY ./scripts/wait-for-it.sh /
RUN chmod +x /wait-for-it.sh
