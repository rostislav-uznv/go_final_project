package api

import (
	"net/http"
	"time"
)

const layout = "20060102"

func Init() {
	http.HandleFunc("/api/nextdate", nextDayHandler)
	http.HandleFunc("/api/task", auth(taskHandler))
	http.HandleFunc("/api/task/done", auth(taskDoneHandler))
	http.HandleFunc("/api/tasks", auth(tasksHandler))
	http.HandleFunc("/api/signin", signinHandler)
}

func nextDayHandler(w http.ResponseWriter, r *http.Request) {

	var timeNow time.Time
	var err error

	now := r.FormValue("now")
	date := r.FormValue("date")
	repeat := r.FormValue("repeat")

	if now == "" {
		timeNow = time.Now()
	} else {
		timeNow, err = time.Parse(layout, now)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	}
	res, err := NextDate(timeNow, date, repeat)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	_, err = w.Write([]byte(res))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func taskHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		addTaskHandler(w, r)
	case http.MethodGet:
		getTaskHandler(w, r)
	case http.MethodPut:
		updateTaskHandler(w, r)
	case http.MethodDelete:
		deleteTaskHandler(w, r)
	default:
		writeJson(w, http.StatusMethodNotAllowed, map[string]string{"error": "Method not allowed"})
	}
}
