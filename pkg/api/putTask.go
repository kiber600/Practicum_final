package api

import (
	"Practicum_final/pkg/db"
	"encoding/json"
	"log"
	"net/http"
)

func putHandler(w http.ResponseWriter, r *http.Request) {
	var task db.Task

	defer r.Body.Close()

	err := json.NewDecoder(r.Body).Decode(&task)
	if err != nil {
		log.Printf("Invalid JSON: %v\n", err.Error())
		writeError(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	if task.Title == "" {
		log.Println("Title is empty")
		writeError(w, "Title is empty", http.StatusMisdirectedRequest)
		return
	}

	err = checkDate(&task)
	if err != nil {
		log.Printf("Error check date: %v\n", err.Error())
		writeError(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = db.UpdateTask(&task)
	if err != nil {
		log.Printf("Error update task: %v\n", err)
		writeError(w, "Task not found", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	writeJson(w, task)

}
