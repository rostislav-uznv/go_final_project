package api

import (
	"net/http"

	"github.com/rostislav-uznv/go_final_project/pkg/db"
)

func getTaskHandler(w http.ResponseWriter, r *http.Request) {
	var task *db.Task
	id := r.URL.Query().Get("id")
	task, err := db.GetTask(id)
	if err != nil {
		writeJson(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}
	writeJson(w, http.StatusOK, task)
}
