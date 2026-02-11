package api

import (
	"Practicum_final/pkg/db"
	"net/http"
	"strconv"
)

func getHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		writeError(w, err.Error(), http.StatusBadRequest)
		return
	}

	if id < 0 {
		writeError(w, "ID must be positive", http.StatusBadRequest)
		return
	}

	t, err := db.GetTask(id)

	if err != nil {
		writeError(w, "Task not found", http.StatusBadRequest)
		return
	}
	writeJson(w, t)
}
