package services

import (
	"context"
	"fmt"
	"io/fs"
	"log"
	services "metasploit-db/services/model"
	"os"
	"path/filepath"
	"strings"

	"github.com/gocarina/gocsv"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

type Neo4jService struct {
	Driver neo4j.DriverWithContext
}

func NewNeo4jService(uri, username, password string) *Neo4jService {
	driver, err := neo4j.NewDriverWithContext(uri, neo4j.BasicAuth(username, password, ""))
	if err != nil {
		log.Fatalf("Failed to create Neo4j driver: %v", err)
	}
	return &Neo4jService{Driver: driver}
}

func (s *Neo4jService) Close() {
	s.Driver.Close(context.Background())
}

func (s *Neo4jService) LoadDirectoryToNeo4j(directoryPath string) error {

	// Walk through all files in the directory
	err := filepath.WalkDir(directoryPath, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return fmt.Errorf("error accessing path %s: %w", path, err)
		}

		if !d.IsDir() && filepath.Ext(path) == ".csv" {
			log.Printf("Processing file: %s\n", path)
			if err := s.importCSVToNeo4j(path); err != nil {
				log.Printf("Error processing file %s: %v", path, err)
			}
		}
		return nil
	})

	if err != nil {
		return fmt.Errorf("error walking through directory: %w", err)
	}

	log.Println("All files processed successfully.")
	return nil
}

func (s *Neo4jService) importCSVToNeo4j(filePath string) error {
	ctx := context.Background()
	// Open the CSV file
	file, err := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE, os.ModePerm)
	if err != nil {
		return fmt.Errorf("error opening CSV file: %w", err)
	}
	defer file.Close()

	var records []services.CSVRecord

	if err := gocsv.Unmarshal(file, &records); err != nil {
		return fmt.Errorf("error reading CSV file: %w", err)
	}

	// Process each record and insert it into Neo4j
	session := s.Driver.NewSession(ctx, neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close(ctx)

	fileName := filepath.Base(filePath)
	parts := strings.SplitN(fileName, "+", 2)
	var payload, exploit string
	if len(parts) == 2 {
		payload, exploit = parts[0], parts[1]
	} else {
		payload = fileName
		exploit = "unknown"
	}

	if err := s.createExploit(ctx, session, exploit, payload, fileName); err != nil {
		return err
	}

	if err := s.createPacketsInExploit(ctx, session, records, exploit, payload, filePath); err != nil {
		return err
	}

	fmt.Println("Data successfully imported into Neo4j.")
	return nil
}

func (s *Neo4jService) createExploit(ctx context.Context, session neo4j.SessionWithContext, exploit, payload, fileName string) error {
	_, err := session.Run(ctx, `
		CREATE (e:Exploit {name: $fileName, payload: $payload})
	`, map[string]interface{}{"fileName": exploit, "payload": payload})
	if err != nil {
		return fmt.Errorf("error creating exploit for file %s: %w", fileName, err)
	}

	return nil
}

func (s *Neo4jService) createPacketsInExploit(ctx context.Context, session neo4j.SessionWithContext,
	records []services.CSVRecord, exploit, payload, filePath string) error {
	for _, record := range records {

		query := `
		MATCH (e:Exploit {name: $fileName, payload: $payload})
		CREATE (p:Packet {
			seq: $seq,
			size: $size,
			protocol: $protocol,
			request: $request,
			body: $body
		})-[:BELONGS_TO]->(e)
	`

		// Execute the query
		_, err := session.Run(ctx, query, map[string]interface{}{
			"fileName": exploit,
			"payload":  payload,
			"seq":      record.Seq,
			"size":     record.Size,
			"protocol": record.Protocol,
			"request":  record.Request,
			"body":     record.Body,
		})
		if err != nil {
			log.Printf("Error inserting record from file %s with seq %s: %v", filePath, record.Seq, err)
		}
	}

	return nil
}
