package controllers

import (
	"fmt"
	"log"
	"metasploit-db/services"
	"net/http"
)

type MetasploitController struct {
	DBService         *services.Neo4jService
	MetasploitService *services.MetasploitService
}

func NewMetasploitController(dbService *services.Neo4jService, metasploitService *services.MetasploitService) *MetasploitController {
	return &MetasploitController{DBService: dbService, MetasploitService: metasploitService}
}

// StartExecution triggers the execution of all single payloads.
func (mc *MetasploitController) LoadWPandPHP(w http.ResponseWriter, r *http.Request) error {
	log.Println("Saving all exploits from wp and php...")

	err := mc.MetasploitService.SaveWPandPHPexploits()
	if err != nil {
		http.Error(w, "Failed executing exploits", http.StatusInternalServerError)
		return err
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, "All exploits are correctly saved")
	return nil
}
