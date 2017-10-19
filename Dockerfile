FROM golang:1.9 as gobuilder
WORKDIR /go/src/github.com/ibigbug/ss-account
COPY . .
RUN go build -o app main.go

FROM node:8 as nodebuilder
WORKDIR /frontend
COPY ./frontend/ /frontend
RUN npm i -g yarn
RUN yarn
RUN ./node_modules/.bin/ng build --prod --base-href /dashboard/

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=gobuilder /go/src/github.com/ibigbug/ss-account/app .
COPY --from=nodebuilder /frontend/dist/ ./public
CMD [ "./app" ]
