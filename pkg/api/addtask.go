package api

import (
	"Practicum_final/pkg/db"
	"encoding/json"
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
	}

}

func addTaskHandler(w http.ResponseWriter, r *http.Request) {
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

	if task.Date == "" {
		task.Date = time.Now().Format(dateFormat)
	}
	err = checkDate(&task)
	if err != nil {
		writeError(w, err.Error(), http.StatusBadRequest)
		return
	}

	var resp = make(map[string]int64)
	resp["id"], err = db.AddTask(&task)

	if err != nil {
		writeJson(w, err.Error())
		return
	}

	respJson, err := json.Marshal(resp)
	if err != nil {
		writeJson(w, err)
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
		return err
	}

	if now.After(t) && task.Repeat != "" {
		nextD, err := nextDate(now, task.Date, task.Repeat)
		if err != nil {
			return err
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
