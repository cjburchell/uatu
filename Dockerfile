FROM golang:1.14 as serverbuilder
WORKDIR /uatu
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o main

FROM scratch

COPY --from=serverbuilder /uatu/main  /server/main

WORKDIR  /server

CMD ["./main"]

