package services

type AttackL7Neo4j struct {
	ID      string        `json:"id"`
	Payload string        `json:"payload"`
	Packets []interface{} `json:"packets"`
	Action  string        `json:"action"`
}
