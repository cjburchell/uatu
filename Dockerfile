FROM node:10-alpine as uibuilder
COPY ui /ui
RUN cd /ui/yasls && npm install
RUN cd /ui/yasls && node_modules/@angular/cli/bin/ng build --prod

FROM golang:1.8.0-alpine as serverbuilder
WORKDIR /go/src/github.com/cjburchell/yasls
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o main

FROM scratch

COPY --from=uibuilder /ui/yasls/dist  /server/ui/yasls/dist
COPY --from=serverbuilder /go/src/github.com/cjburchell/yasls/main  /server/main

WORKDIR  /server

CMD ["./main"]

