FROM golang:1.11-alpine3.7

RUN apk add --no-cache git mercurial \
    && apk add nodejs=8.9.3-r1 \
    && npm install -g yarn

WORKDIR /go/src/app
COPY . .

RUN go get -d -v ./... \
    && go install -v ./... \
    && apk del git mercurial \
    && cd /go/src/github.com/ekrem95/go-gin/app/ && yarn \
    && yarn build

CMD ["app"]