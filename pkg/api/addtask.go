package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/rostislav-uznv/go_final_project/pkg/db"
)

func addTaskHandler(w http.ResponseWriter, r *http.Request) {
	var task db.Task
	var buf bytes.Buffer
	_, err := buf.ReadFrom(r.Body)
	if err != nil {
		writeJson(w, map[string]string{"error": err.Error()})
		return
	}

	if err = json.Unmarshal(buf.Bytes(), &task); err != nil {
		writeJson(w, map[string]string{"error": err.Error()})
		return
	}

	if task.Title == "" {
		writeJson(w, map[string]string{"error": "empty title"})
		return
	}
	// fmt.Println("Read task from form, comment:", task.Comment)
	// fmt.Println("Read task from form, date:", task.Date)
	// fmt.Println("Read task from form, repeat:", task.Repeat)
	// fmt.Println("Read task from form, title:", task.Title)
	// fmt.Println()

	if err = checkDate(&task); err != nil {
		writeJson(w, map[string]string{"error": err.Error()})
		return
	}

	// _, err = NextDate(time.Now(), task.Date, task.Repeat)
	// if err != nil {
	// 	http.Error(w, err.Error(), http.StatusBadRequest)
	// 	return
	// }

	id, err := db.AddTask(&task)
	if err != nil {
		writeJson(w, map[string]string{"error": err.Error()})
		return
	}

	fmt.Println("id:", id)
	fmt.Println()

	writeJson(w, map[string]string{"id": fmt.Sprintf("%d", id)})
}

func checkDate(task *db.Task) error {

	var next string
	var t time.Time
	var err error

	now := time.Now()
	if task.Date == "" {
		task.Date = now.Format(layout)
	}

	if t, err = time.Parse(layout, task.Date); err != nil {
		return err
	}

	if task.Repeat != "" {
		next, err = NextDate(now, task.Date, task.Repeat)
		if err != nil {
			return err
		}
	}

	if afterNow(now, t) {
		if len(task.Repeat) == 0 {
			task.Date = now.Format(layout)
		} else {
			task.Date = next
		}
	}

	return nil
}

func writeJson(w http.ResponseWriter, data any) {
	resp, err := json.Marshal(data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	w.Write(resp)
}
