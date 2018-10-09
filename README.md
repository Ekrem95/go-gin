## go-gin

A React app with Go backend.

#### Prerequisites

```
Go 1.6 or higher
A working golang environment
Docker/Docker Compose
Node.js v8+
Yarn (optional)
```

#### Installing

```
go get github.com/ekrem95/go-gin
```

Go to React app directory
```
cd $GOPATH/src/github.com/ekrem95/go-gin/app
```

Install dependincies and build the app
```
yarn && yarn build
```

Go to project directory and run `docker-compose.yml` file to start databases
```
cd $GOPATH/src/github.com/ekrem95/go-gin && docker-compose up -d
```

Run the tests
```
go test -v ./...
```

Run the app with `go run main.go` or `go-gin`
