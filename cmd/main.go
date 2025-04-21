package main

import (
	"log"
	"net/http"
	"os"

	"metasploit-db/controllers"
	"metasploit-db/routes"
	"metasploit-db/services"

	"github.com/gorilla/mux"
)

func main() {
	router := mux.NewRouter()
	port := os.Getenv("SERVER_PORT")
	databaseURL := os.Getenv("NEO4J_DB")
	user := os.Getenv("NEO4J_USER")
	password := os.Getenv("NEO4J_PASSWORD")

	neo4jService, err := services.NewNeo4jService(databaseURL, user, password)
	if err != nil {
		log.Fatalf("Failed to connect to Neo4j: %v", err)
	}
	defer neo4jService.Close()

	metasploitService := services.NewMetasploitService(neo4jService)
	metasploitController := controllers.NewMetasploitController(neo4jService, metasploitService)
	routes.RegisterRoutes(router, metasploitController)

	log.Printf("Server running on %s", port)
	log.Fatal(http.ListenAndServe(":"+port, router))
}
