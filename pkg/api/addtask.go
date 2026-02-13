package api

import (
	"Practicum_final/pkg/db"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

func taskHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		getHandler(w, r)
	case http.MethodPost:
		addTaskHandler(w, r)
	case http.MethodDelete:
		deleteTaskHandler(w, r)
	case http.MethodPut:
		putHandler(w, r)
	default:
		log.Println("Method not allowed")
		writeError(w, "Method not allowed", http.StatusMethodNotAllowed)
	}

}

func addTaskHandler(w http.ResponseWriter, r *http.Request) {
	var task db.Task

	defer r.Body.Close()

	err := json.NewDecoder(r.Body).Decode(&task)
	if err != nil {
		log.Printf("Invalid JSON: %w\n", err.Error())
		writeError(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	if task.Title == "" {
		log.Println("Title is empty")
		writeError(w, "Title is empty", http.StatusMisdirectedRequest)
		return
	}

	if task.Date == "" {
		task.Date = time.Now().Format(dateFormat)
	}
	err = checkDate(&task)
	if err != nil {
		log.Printf("Error check date: %v\n", err)
		writeError(w, err.Error(), http.StatusMisdirectedRequest)
		return
	}

	var resp = make(map[string]int64)
	resp["id"], err = db.AddTask(&task)

	if err != nil {
		log.Printf("Error: %v\n", err)
		writeError(w, err.Error(), http.StatusNotImplemented)
		return
	}

	respJson, err := json.Marshal(resp)
	if err != nil {
		log.Printf("Error Marshal JSON: %v\n", err)
		writeError(w, err.Error(), http.StatusNotImplemented)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	w.Write(respJson)
}

func checkDate(task *db.Task) error {
	now := time.Now()

	if task.Date == now.Format(dateFormat) {
		return nil
	}

	if task.Date == "" {
		task.Date = time.Now().Format(dateFormat)
	}

	t, err := time.Parse(dateFormat, task.Date)
	if err != nil {
		return fmt.Errorf("Error parsing date: %w\n", err)
	}

	if now.After(t) && task.Repeat != "" {
		nextD, err := nextDate(now, task.Date, task.Repeat)
		if err != nil {
			return fmt.Errorf("Error getting next date: %w\n", err)
		}
		task.Date = nextD
	}

	if now.After(t) && task.Repeat == "" {
		nextD := now
		task.Date = nextD.Format(dateFormat)
		return nil
	}

	return nil

}

func writeJson(w http.ResponseWriter, data any) {

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	jsonData, err := json.Marshal(data)
	if err != nil {
		log.Println("Failed to marshal JSON")
		w.WriteHeader(http.StatusInternalServerError)
		errorResponse := map[string]string{"error": "Failed to marshal JSON"}
		json.NewEncoder(w).Encode(errorResponse)
		return
	}

	w.Write(jsonData)
}

func writeError(w http.ResponseWriter, message string, statusCode int) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(map[string]string{
		"error": message,
	})
}
