package server

import (
	"net/http"
	"os"
	"strconv"

	"github.com/rostislav-uznv/go_final_project/pkg/api"
)

func Run() {
	api.Init()
}

func Start() {
	DEFAULT_TODO_PORT := "7540"
	TODO_PORT := os.Getenv("TODO_PORT")
	if _, err := strconv.Atoi(TODO_PORT); err != nil {
		TODO_PORT = DEFAULT_TODO_PORT
	}

	webDir := "./web"
	http.Handle("/", http.FileServer(http.Dir(webDir)))
	err := http.ListenAndServe(":"+TODO_PORT, nil)
	if err != nil {
		panic(err)
	}
}
