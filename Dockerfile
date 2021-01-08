
FROM golang:latest AS build

WORKDIR  /go/src
COPY . .

RUN GO111MODULE=on \
  CGO_ENABLED=0 \
  GOOS=linux \
  GOARCH=amd64 \
  go build -o forum cmd/forum-api/main.go

FROM ubuntu:20.04

RUN apt-get -y update && apt-get install -y tzdata

ENV TZ=Russia/Moscow
RUN ln -snf /usr/share/zoneinfo/$TZ /etc/localtime && echo $TZ > /etc/timezone

ENV PGVER 12
RUN apt-get -y update && apt-get install -y postgresql-$PGVER

USER postgres

RUN /etc/init.d/postgresql start &&\
    psql --command "CREATE USER forum WITH SUPERUSER PASSWORD 'postgres';" &&\
    createdb -O forum forum &&\
    /etc/init.d/postgresql stop

VOLUME  ["/etc/postgresql", "/var/log/postgresql", "/var/lib/postgresql"]

USER root

WORKDIR /usr/src/app

COPY . .
COPY --from=build /go/src/forum .

ENV DB_HOST=127.0.0.1 \
    DB_PORT=5432 \
    DB_NAME=forum \
    DB_USER=forum \
    DB_PASSWORD=postgres \
    PGPASSWORD=postgres

COPY scripts/init.sql .

CMD service postgresql start && psql -h localhost -d forum -U forum -p 5432 -aq -f init.sql && ./forum -p 5000

EXPOSE 5432
EXPOSE 5000
