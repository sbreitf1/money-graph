package moneydb

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"sync"
	"time"
)

type Database struct {
	localDir string
	Name     string  `json:"name"`
	Groups   []Group `json:"groups"`
	mutex    sync.Mutex
}

type Group []struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type chunkDescriptor struct {
	IBAN  IBAN
	Year  int
	Month time.Month
}

func (desc chunkDescriptor) String() string {
	return fmt.Sprintf("%04d-%02d [%s]", desc.Year, desc.Month, desc.IBAN.String())
}

type chunk struct {
	desc    chunkDescriptor
	Entries []Entry `json:"entries"`
}

type Entry struct {
	Hash      string
	GroupID   string
	IBAN      IBAN
	Date      time.Time
	Type      string
	Message   string
	Amount    Money
	OtherName string
	OtherIBAN IBAN
}

func (e Entry) ChunkDescriptor() chunkDescriptor {
	return chunkDescriptor{IBAN: e.IBAN, Year: e.Date.Year(), Month: e.Date.Month()}
}

func ExistsInFolder(localDir string) (bool, error) {
	fi, err := os.Stat(filepath.Join(localDir, "moneydb.json"))
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, err
	}
	return !fi.IsDir(), nil
}

func CreateInFolder(localDir, name string) (*Database, error) {
	if exists, err := ExistsInFolder(localDir); err != nil {
		return nil, fmt.Errorf("check for existing db: %s", err.Error())
	} else if exists {
		return nil, fmt.Errorf("%q already contains a money-db", localDir)
	}

	db := &Database{
		localDir: localDir,
		Name:     name,
		Groups:   make([]Group, 0),
	}

	if err := db.Save(); err != nil {
		return nil, err
	}

	return db, nil
}

func OpenFromFolder(localDir string) (*Database, error) {
	data, err := os.ReadFile(filepath.Join(localDir, "moneydb.json"))
	if err != nil {
		return nil, fmt.Errorf("read db json: %s", err.Error())
	}

	var db *Database
	if err := json.Unmarshal(data, &db); err != nil {
		return nil, fmt.Errorf("unmarshal db json: %s", err.Error())
	}

	db.localDir = localDir
	return db, nil
}

func OpenOrCreateInFolder(localDir, name string) (*Database, error) {
	exists, err := ExistsInFolder(localDir)
	if err != nil {
		return nil, err
	}

	if exists {
		db, err := OpenFromFolder(localDir)
		if err != nil {
			return nil, err
		}
		if db.Name != name {
			db.Name = name
			if err := db.Save(); err != nil {
				return nil, err
			}
		}
		return db, nil
	}
	return CreateInFolder(localDir, name)
}

func (db *Database) Save() error {
	data, err := json.MarshalIndent(&db, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal db json: %s", err.Error())
	}
	if err := os.MkdirAll(db.localDir, os.ModePerm); err != nil {
		return fmt.Errorf("create db directory: %s", err.Error())
	}
	if err := os.WriteFile(filepath.Join(db.localDir, "moneydb.json"), data, os.ModePerm); err != nil {
		return fmt.Errorf("export db json: %s", err.Error())
	}
	return nil
}

func (db *Database) getChunkPath(desc chunkDescriptor) string {
	return filepath.Join(db.localDir, "data", desc.IBAN.String(), "chunks", fmt.Sprintf("%04d-%02d.json", desc.Year, int(desc.Month)))
}

func (db *Database) loadChunk(desc chunkDescriptor) (*chunk, error) {
	path := db.getChunkPath(desc)

	rawData, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return &chunk{
				desc:    desc,
				Entries: []Entry{},
			}, nil
		}
	}

	var chunk chunk
	if err := json.Unmarshal(rawData, &chunk); err != nil {
		return nil, err
	}

	chunk.desc = desc
	return &chunk, nil
}

func (db *Database) saveChunk(c *chunk) error {
	path := db.getChunkPath(c.desc)

	rawData, err := json.MarshalIndent(&c, "", "  ")
	if err != nil {
		return err
	}

	if err := os.MkdirAll(filepath.Dir(path), os.ModePerm); err != nil {
		return err
	}
	return os.WriteFile(path, rawData, os.ModePerm)
}

func (c *chunk) AddEntry(newEntry Entry) (bool, error) {
	if c.desc != newEntry.ChunkDescriptor() {
		return false, fmt.Errorf("given entry does not belong to chunk %v", c.desc)
	}

	for _, e := range c.Entries {
		if e.Hash == newEntry.Hash {
			// entry already in list
			return false, nil
		}
	}

	c.Entries = append(c.Entries, newEntry)
	sort.SliceStable(c.Entries, func(i, j int) bool {
		return c.Entries[i].Date.Before(c.Entries[j].Date)
	})

	return true, nil
}
