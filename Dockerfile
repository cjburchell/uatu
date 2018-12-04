FROM golang:1.9.0-alpine as serverbuilder
WORKDIR /go/src/github.com/cjburchell/uatu
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o main

FROM scratch

COPY --from=serverbuilder /go/src/github.com/cjburchell/uatu/main  /server/main

WORKDIR  /server

CMD ["./main"]

