FROM golang:latest AS build_step
ENV GO111MODULE=on
WORKDIR  /go/src
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /go/build/forum /go/src/cmd/server/main.go

FROM alpine
WORKDIR /app

COPY --from=build_step /go/build/forum .
RUN chmod +x forum

ENTRYPOINT ["forum"]

EXPOSE 8000/tcp
