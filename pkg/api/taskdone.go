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
		writeError(w, "ID is required", http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(idUrl)
	if err != nil {
		writeError(w, "Id invalid", http.StatusBadRequest)
		return
	}
	if id < 0 {
		writeError(w, "ID must be positive", http.StatusBadRequest)
		return
	}

	now := time.Now()
	var nextD string

	task, err := db.GetTask(id)
	log.Println(task)
	if err != nil {
		writeError(w, "No data with this id was found", http.StatusNotFound)
		return
	}
	if task.Repeat == "" {
		err = db.DeleteTask(id)
		if err != nil {
			writeError(w, "Error delete task", http.StatusBadRequest)
			return
		}
		w.Write([]byte("{}"))

		return
	}

	nextD, err = nextDate(now, task.Date, task.Repeat)
	if err != nil {
		writeError(w, err.Error(), http.StatusNotFound)
		return
	}

	err = db.UpdateDate(nextD, id)
	if err != nil {
		writeError(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte("{}"))

}
