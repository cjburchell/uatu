FROM golang:1.8.0-alpine as serverbuilder
WORKDIR /go/src/github.com/cjburchell/yasls
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o main

FROM scratch

COPY --from=serverbuilder /go/src/github.com/cjburchell/yasls/main  /server

WORKDIR  /server

CMD ["./main"]

