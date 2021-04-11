package storage

import (
	"database/sql"
	"fmt"
	"strconv"
	"strings"
)

// SchemaVersion is a simple int representation of the schema version.
type SchemaVersion int

// ToString converts the schema version into a string and returns that converted value.
func (s SchemaVersion) ToString() string {
	return strconv.Itoa(int(s))
}

type transaction interface {
	Exec(query string, args ...interface{}) (sql.Result, error)
}

type CreateTableIndexStmt struct {
	Index         string
	Table         string
	CommonStmt    string
	SpecificStmts map[string]string
}

func (ctis CreateTableIndexStmt) GetStmt(provider string) (stmt string) {
	if _, ok := ctis.SpecificStmts[provider]; ok {
		return ctis.SpecificStmts[provider]
	}

	switch provider {
	case "mysql":
		return ctis.mySQLStmt()
	case "postgres":
		return ctis.postgreStmt()
	default:
		return ctis.CommonStmt
	}
}

func (ctis CreateTableIndexStmt) mySQLStmt() (stmt string) {
	return strings.Replace(ctis.CommonStmt, "CREATE INDEX IF NOT EXISTS", "CREATE INDEX", 1)
}

func (ctis CreateTableIndexStmt) postgreStmt() (stmt string) {
	stmt = ctis.CommonStmt
	index := 1
	for _, rune := range stmt {
		if rune == '?' {
			stmt = strings.Replace(stmt, "?", fmt.Sprintf("$%d", index), 1)
			index = index + 1
		}
	}

	return stmt
}
