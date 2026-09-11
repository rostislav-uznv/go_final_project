package main

import (
	"fmt"
	"os"

	"github.com/rostislav-uznv/go_final_project/pkg/db"
	"github.com/rostislav-uznv/go_final_project/pkg/server"
)

func main() {

	dbFile := "scheduler.db"
	TODO_DBFILE := os.Getenv("TODO_DBFILE")
	if TODO_DBFILE != "" {
		dbFile = TODO_DBFILE
	}
	err := db.Init(dbFile)
	if err != nil {
		fmt.Println("Ошибка инициализации базы данных")
		return
	}

	server.Run()
	server.Start()

}
