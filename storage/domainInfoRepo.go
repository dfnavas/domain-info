package storage

import (
	"database/sql"
	"encoding/json"
	"log"
	"time"
)

// DomainInfoRepo :
type DomainInfoRepo interface {
	Get(url string) (*DomainInfo, error)
	Upsert(seed DomainInfo) error
	GetAll() []*DomainHistory
}

type domainInfoRepoDeps struct {
	db *sql.DB
}

// DomainInfo :
type DomainInfo struct {
	Host        string       `json:"host"`
	IsDown      bool         `json:"is_down"`
	Severs      []ServerInfo `json:"servers"`
	SSLGrade    string       `json:"ssl_grade"`
	Title       string       `json:"title"`
	Logo        string       `json:"logo"`
	LastUpdated int64        `json:"last_updated"`
}

// ServerInfo :
type ServerInfo struct {
	Address  string `json:"address"`
	SSLGrade string `json:"ssl_grade"`
	Country  string `json:"country"`
	Owner    string `json:"owner"`
}

// DomainHistory :
type DomainHistory struct {
	Host     string       `json:"host"`
	Severs   []ServerInfo `json:"servers"`
	SSLGrade string       `json:"ssl_grade"`
	Title    string       `json:"title"`
	Logo     string       `json:"logo"`
}

// NewDomainInfoRepo :
func NewDomainInfoRepo(dbURL string) (DomainInfoRepo, error) {
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		return nil, err
	}
	initDB(db)
	return &domainInfoRepoDeps{db}, nil
}

func initDB(db *sql.DB) error {
	if _, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS domains (
			host VARCHAR NOT NULL PRIMARY KEY,
			data JSON NOT NULL,
			last_updated BIGINT NOT NULL
		);
	`); err != nil {
		log.Fatalln("Error:", err)
		return err
	}
	return nil
}

// Get :
func (dir *domainInfoRepoDeps) Get(url string) (*DomainInfo, error) {
	rows, err := dir.db.Query(`SELECT d.host,d.data,d.last_updated FROM domains d WHERE d.host = $1	LIMIT 1`, url)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var host string
		var data string
		var lastUpdate int64
		rows.Scan(&host, &data, &lastUpdate)
		var domRecord DomainInfo
		json.Unmarshal([]byte(data), &domRecord)
		domRecord.LastUpdated = lastUpdate
		return &domRecord, nil
	}
	return nil, nil
}

// Upsert :
func (dir *domainInfoRepoDeps) Upsert(seed DomainInfo) error {
	seedJSON, err := json.Marshal(seed)
	if err != nil {
		log.Println("Error:", err)
		return err
	} else if _, err := dir.db.Exec(
		`UPSERT INTO domains (host, data, last_updated) VALUES ($1, $2, $3)`,
		seed.Host, seedJSON, time.Now().UnixNano()); err != nil {
		log.Fatalln("Error:", err)
		return err
	}
	return nil
}

// GetAll :
func (dir *domainInfoRepoDeps) GetAll() []*DomainHistory {
	records := make([]*DomainHistory, 0)
	rows, err := dir.db.Query(`SELECT d.host,d.data FROM domains d`)
	if err != nil {
		log.Println("Error:", err)
		return records
	}
	defer rows.Close()
	for rows.Next() {
		var host string
		var data string
		rows.Scan(&host, &data)
		var domRecord DomainHistory
		json.Unmarshal([]byte(data), &domRecord)
		records = append(records, &domRecord)
	}
	return records
}
