package main

import (
	"log"

	"github.com/ekrem95/go-gin/db"
	"github.com/ekrem95/go-gin/router"
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	if err := db.TestSQLConnection(); err != nil {
		log.Fatal(err)
	}

	r := router.Default()

	r.Run(":8080")
}
