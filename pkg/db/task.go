package db

import (
	"database/sql"
	"fmt"
	"strconv"

	_ "modernc.org/sqlite"
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

	query := `INSERT INTO scheduler (date, title, comment, repeat) VALUES (:date, :title, :comment, :repeat)`
	res, err := db.Exec(query,
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
	res := []*Task{}

	rows, err := db.Query("SELECT id, date, title, comment, repeat from scheduler ORDER BY date DESC limit ?", limit)
	if err != nil {
		return res, err
	}

	for rows.Next() {
		s := Task{}
		err = rows.Scan(&s.ID, &s.Date, &s.Title, &s.Comment, &s.Repeat)

		if err != nil {
			return res, err
		}
		res = append(res, &s)

	}
	err = rows.Err()
	if err != nil {
		return res, err
	}

	return res, nil
}

func GetTask(id int) (*Task, error) {
	task := &Task{}

	raw := db.QueryRow("Select * from scheduler where id = :id", sql.Named("id", id))
	err := raw.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
	if err != nil {
		return task, err
	}

	err = raw.Err()
	if err != nil {
		return task, err
	}

	return task, nil
}

func UpdateTask(task *Task) error {
	// параметры пропущены, не забудьте указать WHERE

	id, err := strconv.ParseInt(task.ID, 10, 64)
	if err != nil {
		return err
	}
	if id < 0 {
		return fmt.Errorf("Id must be positive")
	}
	query := `UPDATE scheduler SET date = :date, title = :title, comment = :comment, repeat = :repeat where id = :id`
	res, err := db.Exec(query,
		sql.Named("date", task.Date),
		sql.Named("title", task.Title),
		sql.Named("comment", task.Comment),
		sql.Named("repeat", task.Repeat), sql.Named("id", id))
	if err != nil {
		return err
	}
	// метод RowsAffected() возвращает количество записей к которым
	// был применена SQL команда
	count, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if count == 0 {
		return fmt.Errorf(`incorrect id for updating task`)
	}
	return nil
}

func UpdateDate(nextDate string, id int) error {

	query := `UPDATE scheduler SET date = :date where id = :id`
	res, err := db.Exec(query, sql.Named("date", nextDate), sql.Named("id", id))
	if err != nil {
		return err
	}

	count, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if count == 0 {
		return fmt.Errorf(`incorrect id for updating task`)
	}
	return nil

}

func DeleteTask(id int) error {

	_, err := db.Exec("Delete from scheduler where id = :id", sql.Named("id", id))
	if err != nil {
		return err
	}

	return nil
}
