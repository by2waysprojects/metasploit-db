package routes

import (
	"metasploit-db/controllers"
	"net/http"

	"github.com/gorilla/mux"
)

// InitializeRouter sets up application routes.
func RegisterRoutes(router *mux.Router, metasploitController *controllers.MetasploitController) {
	router.HandleFunc("/save-wp-php", func(w http.ResponseWriter, r *http.Request) {
		err := metasploitController.LoadWPandPHP(w, r)
		if err != nil {
			http.Error(w, "Failed to execute payloads", http.StatusInternalServerError)
		}
	}).Methods("GET")
}
