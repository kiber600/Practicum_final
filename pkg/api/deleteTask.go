package api

import (
	"Practicum_final/pkg/db"
	"log"
	"net/http"
	"strconv"
)

func deleteTaskHandler(w http.ResponseWriter, r *http.Request) {

	idUrl := r.URL.Query().Get("id")

	id, err := strconv.Atoi(idUrl)
	if err != nil {
		log.Printf("Convert id: %v\n", err)
		writeError(w, "Id invalid", http.StatusBadRequest)
		return
	}
	if id < 0 {
		log.Println("ID must be positive")
		writeError(w, "ID must be positive", http.StatusBadRequest)
		return
	}

	err = db.DeleteTask(id)
	if err != nil {
		log.Printf("Delete task error: %v\n", err)
		writeError(w, "Delete error", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("{}"))

}
