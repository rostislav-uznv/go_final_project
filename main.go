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
	db.Init(dbFile)
	fmt.Println("Запускаем сервер")

	server.Run()
	server.Start()
	fmt.Println("Завершаем работу")

}
