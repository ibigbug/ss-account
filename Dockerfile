FROM golang:1.9

WORKDIR $HOME/go/src/github.com/ibigbug/ss-account

COPY . .

RUN go build -o app main.go

CMD ["app"]
