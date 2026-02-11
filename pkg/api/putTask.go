package api

import (
	"Practicum_final/pkg/db"
	"encoding/json"
	"net/http"
)

func putHandler(w http.ResponseWriter, r *http.Request) {
	var task db.Task

	defer r.Body.Close()

	err := json.NewDecoder(r.Body).Decode(&task)
	if err != nil {
		writeError(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	if task.Title == "" {
		writeError(w, "Title is empty", http.StatusBadRequest)
		return
	}

	err = checkDate(&task)
	if err != nil {
		writeError(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = db.UpdateTask(&task)
	if err != nil {
		writeError(w, "Task not found", http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	writeJson(w, task)

}
