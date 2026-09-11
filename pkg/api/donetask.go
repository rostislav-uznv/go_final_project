package api

import (
	"net/http"
	"time"

	"github.com/rostislav-uznv/go_final_project/pkg/db"
)

func taskDoneHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeJson(w, http.StatusMethodNotAllowed, map[string]string{"error": "Method not allowed"})
		return
	}
	var task *db.Task
	id := r.URL.Query().Get("id")
	task, err := db.GetTask(id)
	if err != nil {
		writeJson(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	if task.Repeat == "" {
		err := db.DeleteTask(id)
		if err != nil {
			writeJson(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
			return
		}
		writeJson(w, http.StatusOK, map[string]string{})
		return
	}

	task.Date, err = NextDate(time.Now(), task.Date, task.Repeat)
	if err != nil {
		writeJson(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}
	err = db.UpdateDate(task)
	if err != nil {
		writeJson(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	writeJson(w, http.StatusOK, map[string]string{})

}
