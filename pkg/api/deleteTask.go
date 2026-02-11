package api

import (
	"Practicum_final/pkg/db"
	"net/http"
	"strconv"
)

func deleteTaskHandler(w http.ResponseWriter, r *http.Request) {

	idUrl := r.URL.Query().Get("id")

	id, err := strconv.Atoi(idUrl)
	if err != nil {
		writeError(w, "Id invalid", http.StatusBadRequest)
		return
	}
	if id < 0 {
		writeError(w, "ID must be positive", http.StatusBadRequest)
		return
	}

	err = db.DeleteTask(id)
	if err != nil {
		writeError(w, "Delete error", http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("{}"))

}
