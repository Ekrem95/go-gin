package main

import (
	"log"
	"os"

	"github.com/ekrem95/go-gin/db"
	"github.com/ekrem95/go-gin/router"
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	if err := db.TestSQLConnection(); err != nil {
		log.Fatal(err)
	}

	if _, err := os.Stat(router.UploadPath); os.IsNotExist(err) {
		os.Mkdir(router.UploadPath, 0700)
	}

	r := router.Default()

	r.Run(":8080")
}
