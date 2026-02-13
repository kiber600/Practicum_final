package api

import (
	"Practicum_final/pkg/db"
	"log"
	"net/http"
	"strconv"
	"time"
)

func taskDoneHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		writeError(w, "Required POST method", http.StatusMethodNotAllowed)
		return
	}

	idUrl := r.URL.Query().Get("id")

	if idUrl == "" {
		log.Println("ID is required")
		writeError(w, "ID is required", http.StatusMisdirectedRequest)
		return
	}

	id, err := strconv.Atoi(idUrl)
	if err != nil {
		log.Printf("Convert id to int: %v\n", err)
		writeError(w, "Id invalid", http.StatusConflict)
		return
	}
	if id < 0 {
		log.Println("ID must be positive")
		writeError(w, "ID must be positive", http.StatusConflict)
		return
	}

	now := time.Now()
	var nextD string

	task, err := db.GetTask(id)
	if err != nil {
		log.Printf("Error get task: %v\n", err)
		writeError(w, "No data with this id was found", http.StatusNotFound)
		return
	}
	if task.Repeat == "" {
		err = db.DeleteTask(id)
		if err != nil {
			log.Printf("Error delete task: %v\n", err)
			writeError(w, "Error delete task", http.StatusMisdirectedRequest)
			return
		}
		w.Write([]byte("{}"))

		return
	}

	nextD, err = nextDate(now, task.Date, task.Repeat)
	if err != nil {
		log.Printf("Error getting next date: %v\n", err)
		writeError(w, err.Error(), http.StatusNotFound)
		return
	}

	err = db.UpdateDate(nextD, id)
	if err != nil {
		log.Printf("Error update task: %v\n", err)
		writeError(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte("{}"))

}
