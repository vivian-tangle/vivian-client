package tools

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/dgraph-io/badger"
)

// DBExists check if the badger DB exists
func DBExists(path string) bool {
	if _, err := os.Stat(path + "/MANIFEST"); os.IsNotExist(err) {
		return false
	}

	return true
}

// OpenDB opens the badger database
func OpenDB(dir string, opts badger.Options) (*badger.DB, error) {
	var db *badger.DB
	var err error

	if db, err = badger.Open(opts); err != nil {
		if strings.Contains(err.Error(), "LOCK") {
			if db, err := retry(dir, opts); err == nil {
				log.Println("database unlocked, value log truncated")
				return db, nil
			}
			log.Println("could not unlock database: ", err)
		}
		return nil, err
	}

	return db, nil

}

func retry(dir string, originalOpts badger.Options) (*badger.DB, error) {
	lockPath := filepath.Join(dir, "LOCK")
	if err := os.Remove(lockPath); err != nil {
		return nil, fmt.Errorf(`removing "LOCK": %s`, err)
	}
	retryOpts := originalOpts
	retryOpts.Truncate = true
	db, err := badger.Open(retryOpts)

	return db, err
}

// Serialize PedersonCommit structure into bytes
func (pc *PedersonCommit) Serialize() []byte {
	var res bytes.Buffer
	encoder := gob.NewEncoder(&res)

	err := encoder.Encode(pc)

	HandleErr(err)

	return res.Bytes()
}

// Deserialize bytes into PedersonCommit structure
func (pc *PedersonCommit) Deserialize(data []byte) *PedersonCommit {
	var pedersonCommit PedersonCommit

	decoder := gob.NewDecoder(bytes.NewReader(data))

	err := decoder.Decode(&pedersonCommit)

	HandleErr(err)

	return &pedersonCommit
}

// PendingDomainName is the data structure for storing the information of a pending domain name
type PendingDomainName struct {
	Name      string
	Value     string
	RegTxHash string
}

// Serialize PendingDomainName structure into bytes
func (pdn *PendingDomainName) Serialize() []byte {
	var res bytes.Buffer
	encoder := gob.NewEncoder(&res)

	err := encoder.Encode(pdn)

	HandleErr(err)

	return res.Bytes()
}

// Deserialize bytes into PendingDomainName structure
func (pdn *PendingDomainName) Deserialize(data []byte) *PendingDomainName {
	var pendingDomainName PendingDomainName

	decoder := gob.NewDecoder(bytes.NewReader(data))

	err := decoder.Decode(&pendingDomainName)

	HandleErr(err)

	return &pendingDomainName
}
