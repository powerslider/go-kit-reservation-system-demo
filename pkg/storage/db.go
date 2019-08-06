package storage

import (
	"database/sql"
	"fmt"
	"github.com/doug-martin/goqu/v7"
	_ "github.com/doug-martin/goqu/v7/dialect/sqlite3"
	"github.com/doug-martin/goqu/v7/exec"
	_ "github.com/mattn/go-sqlite3"
	"io/ioutil"
	"os"
	errors "reservations/pkg/error"
	"strings"
)

type Persistence struct {
	DB *goqu.Database
}

type QueryOptions struct {
	Limit  uint
	Offset uint
}

type Transaction func(tx *goqu.TxDatabase) exec.QueryExecutor

func NewDB(dbName string) (*Persistence, error) {
	storageFile := fmt.Sprintf("%s.db", dbName)

	db, err := sql.Open("sqlite3", storageFile)
	if err != nil {
		return nil, errors.DBError.Wrapf(err, "error initializing %s database", dbName)
	}

	if !fileExists(storageFile) {
		if err := createSchema(db, dbName); err != nil {
			return nil, err
		}
	}

	goquDB := goqu.New("sqlite3", db)

	return &Persistence{
		DB: goquDB,
	}, nil
}

func (p *Persistence) Tx(txFunc Transaction) (res sql.Result, err error) {
	tx, err := p.DB.Begin()
	if err != nil {
		return nil, err
	}

	err = tx.Wrap(func() error {
		r, e := txFunc(tx).Exec()
		if e != nil {
			return e
		}
		res = r
		return nil
	})

	return res, nil
}

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

func createSchema(db *sql.DB, dbName string) error {
	file, err := ioutil.ReadFile("sql/reservations.sql")
	if err != nil {
		return errors.DBError.Wrapf(err, "error initializing %s database model", dbName)
	}

	queries := strings.Split(string(file), ";\n")
	for _, q := range queries {
		_, err := db.Exec(q)
		if err != nil {
			return errors.DBError.Wrapf(err, "error executing DDL query: %s", q)
		}
	}

	return nil
}
