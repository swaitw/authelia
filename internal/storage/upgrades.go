package storage

import (
	"fmt"
	"sort"

	"github.com/authelia/authelia/internal/utils"
)

func (p *SQLProvider) upgradeCreateTableStatements(tx transaction, statements map[string]string, existingTables []string) (err error) {
	keys := make([]string, 0, len(statements))
	for k := range statements {
		keys = append(keys, k)
	}

	sort.Strings(keys)

	for _, table := range keys {
		if !utils.IsStringInSlice(table, existingTables) {
			_, err := tx.Exec(fmt.Sprintf(statements[table], table))
			if err != nil {
				return fmt.Errorf("Unable to create table %s: %v", table, err)
			}
		}
	}

	return nil
}

func (p *SQLProvider) upgradeCreateTableIndexStatements(tx transaction, statements []CreateTableIndexStmt) (err error) {
	for _, statement := range statements {
		if p.name == "mysql" {
			err = p.upgradeCreateTableIndexStatementsMySqlIfNotExists(tx, statement)
			if err != nil {
				return err
			}

			continue
		}

		_, err = tx.Exec(statement.GetStmt(p.name), statement.Table, statement.Index)
		if err != nil {
			return err
		}
	}

	return nil
}

func (p *SQLProvider) upgradeCreateTableIndexStatementsMySqlIfNotExists(tx transaction, statement CreateTableIndexStmt) (err error) {
	var exists bool
	err = p.db.QueryRow(checkIndexExistsMySQLStmt, statement.Table, statement.Index).Scan(&exists)

	if err != nil {
		return err
	}

	if exists {
		return nil
	}

	_, err = tx.Exec(statement.GetStmt(p.name), statement.Table, statement.Index)

	return err
}

// upgradeFinalize sets the schema version and logs a message, as well as any other future finalization tasks.
func (p *SQLProvider) upgradeFinalize(tx transaction, version SchemaVersion) (err error) {
	_, err = tx.Exec(p.sqlConfigSetValue, "schema", "version", version.ToString())
	if err != nil {
		return err
	}

	p.log.Debugf("%s%d", storageSchemaUpgradeMessage, version)

	return nil
}

func (p *SQLProvider) upgradeStandard(tx transaction, tables []string, version SchemaVersion) (err error) {
	if _, ok := p.sqlUpgradesCreateTableStatements[version]; ok {
		err = p.upgradeCreateTableStatements(tx, p.sqlUpgradesCreateTableStatements[version], tables)

		if err != nil {
			return err
		}
	}

	if _, ok := p.sqlUpgradesCreateTableIndexesStatements[version]; ok {
		err = p.upgradeCreateTableIndexStatements(tx, p.sqlUpgradesCreateTableIndexesStatements[version])

		if err != nil {
			return err
		}
	}

	return p.upgradeFinalize(tx, version)
}

// upgradeSchemaToVersion001 upgrades the schema to version 1.
func (p *SQLProvider) upgradeSchemaToVersion001(tx transaction, tables []string) (err error) {
	version := SchemaVersion(1)

	return p.upgradeStandard(tx, tables, version)
}
