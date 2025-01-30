package services

import (
	"log"
	"os/exec"
)

const (
	resultsPath = "results/"
)

type MetasploitService struct {
	Neo4jService *Neo4jService
}

func NewMetasploitService(neo4jService *Neo4jService) *MetasploitService {
	return &MetasploitService{
		Neo4jService: neo4jService,
	}
}

func (ms *MetasploitService) SaveWPandPHPexploits(limit int) error {
	cmd := exec.Command("python3", "execute_single_payloads.py")

	// Capturar la salida y los errores
	_, err := cmd.Output()
	if err != nil {
		log.Printf("Error executing script: %s", err)
		return err
	}

	if err := ms.Neo4jService.LoadDirectoryToNeo4j(resultsPath, limit); err != nil {
		log.Printf("Error importing results to Neo4j: %s", err)
		return err
	}

	return nil
}
