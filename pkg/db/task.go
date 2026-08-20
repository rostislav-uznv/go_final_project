package db

import (
	"database/sql"
	"fmt"
)

type Task struct {
	ID      string `json:"id"`
	Date    string `json:"date"`
	Title   string `json:"title"`
	Comment string `json:"comment"`
	Repeat  string `json:"repeat"`
}

func AddTask(task *Task) (int64, error) {
	var id int64
	// определите запрос
	query := `INSERT INTO scheduler (date, title, comment, repeat) VALUES (:date, :title, :comment, :repeat)`
	res, err := db.Exec(query,
		// sql.Named("id", task.ID),
		sql.Named("date", task.Date),
		sql.Named("title", task.Title),
		sql.Named("comment", task.Comment),
		sql.Named("repeat", task.Repeat))
	if err == nil {
		id, err = res.LastInsertId()
	}
	return id, err
}

func Tasks(limit int) ([]*Task, error) {
	tasks := []*Task{}

	rows, err := db.Query("SELECT id, date, title, comment, repeat FROM scheduler ORDER BY date LIMIT :limit", sql.Named("limit", limit))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var tsk Task
		err = rows.Scan(&tsk.ID, &tsk.Date, &tsk.Title, &tsk.Comment, &tsk.Repeat)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, &tsk)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	// var TaskResp
	return tasks, nil
}

func GetTask(id string) (*Task, error) {
	var tsk Task
	row := db.QueryRow("SELECT id, date, title, comment, repeat FROM scheduler WHERE id = :id", sql.Named("id", id))
	err := row.Scan(&tsk.ID, &tsk.Date, &tsk.Title, &tsk.Comment, &tsk.Repeat)
	if err != nil {
		return nil, err
	}

	return &tsk, nil
}

func UpdateTask(task *Task) error {
	query := `UPDATE scheduler SET date = :date, title = :title, comment = :comment, repeat = :repeat WHERE id = :id`
	res, err := db.Exec(query,
		sql.Named("id", task.ID),
		sql.Named("date", task.Date),
		sql.Named("title", task.Title),
		sql.Named("comment", task.Comment),
		sql.Named("repeat", task.Repeat))
	if err != nil {
		return err
	}
	// метод RowsAffected() возвращает количество записей к которым
	// была применена SQL команда
	count, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if count == 0 {
		return fmt.Errorf(`incorrect id for updating task`)
	}
	return nil
}
