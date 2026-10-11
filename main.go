package main

import (
	"fmt"
	"os"

	"github.com/rostislav-uznv/go_final_project/pkg/db"
	"github.com/rostislav-uznv/go_final_project/pkg/server"
)

func main() {

	// TODO_PASSWORD := os.Getenv("TODO_PASSWORD")

	// TODO_PASSWORD = "111"
	dbFile := "scheduler.db"
	TODO_DBFILE := os.Getenv("TODO_DBFILE")
	if TODO_DBFILE != "" {
		dbFile = TODO_DBFILE
	}
	database, err := db.Init(dbFile)
	if err != nil {
		fmt.Println("Ошибка инициализации базы данных")
		return
	}
	defer database.Close()

	server.Run()
	err = server.Start()
	if err != nil {
		fmt.Println("Ошибка запуска базы данных")
		return
	}

}
