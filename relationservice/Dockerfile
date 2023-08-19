FROM golang:1.20.7-alpine3.17

WORKDIR /app

COPY . /app

RUN go env -w GO111MODULE=on
RUN go env -w GOPROXY=https://proxy.golang.com.cn,direct
RUN go mod download

EXPOSE 8080

CMD ["go","run","/app/cmd/relation.go"]