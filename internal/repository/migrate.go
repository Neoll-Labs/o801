/*
 license x
*/

package repository

import (
	"database/sql"
	"fmt"
	"github.com/nelsonstr/o801/internal/model"
	strings2 "github.com/nelsonstr/o801/strings"
	"log"
	"reflect"
	"strings"
)

type migrate struct {
	db      DBInterface
	Query   string
	queries []string
}

// DBInterface is a custom interface that matches the methods used from *sql.DB.
type DBInterface interface {
	Begin() (*sql.Tx, error)
	Exec(query string, args ...interface{}) (sql.Result, error)
	Prepare(query string) (*sql.Stmt, error)
	Query(query string, args ...interface{}) (*sql.Rows, error)
	QueryRow(query string, args ...interface{}) *sql.Row
	Close() error
}

func MigrateDB(db DBInterface) error {
	m := &migrate{db: db}

	log.Printf("start migration.")

	defer log.Printf("end migration.")

	m.createTableScript(model.User{})

	return m.executeTablesScripts()
}

func (m *migrate) executeTablesScripts() error {
	tx, err := m.db.Begin()
	if err != nil {
		return fmt.Errorf("error starting tx: %w", err)
	}

	defer func() { _ = tx.Rollback() }()

	for _, query := range m.queries {
		_, err := tx.Exec(query)
		if err != nil {
			return fmt.Errorf("error executing sql to create table creation: %w", err)
		}
	}

	_ = tx.Commit()

	return nil
}

type Column struct {
	Name  string
	Types string
}

type Table struct {
	Name    string
	Values  map[string]reflect.Value
	Columns []Column
}

func readTags(tags string) map[string][]string {
	attributes := strings.Split(tags, ";")

	v := make(map[string][]string)

	for i := 0; i < len(attributes); i++ {
		pre := strings.Split(attributes[i], ":")
		v[pre[0]] = strings.Split(pre[1], ",")
	}

	return v
}

func (m *migrate) createTableScript(model any) {
	mod := processStruct(model)

	var cols []string
	for i := 0; i < len(mod.Columns); i++ {
		cols = append(cols, fmt.Sprintf("%s %s  ", mod.Columns[i].Name, mod.Columns[i].Types))
	}

	m.Query = fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s (%s);", mod.Name, strings.Join(cols, ", "))

	m.queries = append(m.queries, m.Query)
}

func processStruct(model any) *Table {
	tbl := &Table{
		Values:  make(map[string]reflect.Value),
		Columns: make([]Column, 0),
	}

	typ := reflect.TypeOf(model)
	value := reflect.ValueOf(model)

	tbl.Name = getTableName(typ)

	for i := 0; i < typ.NumField(); i++ {
		const key = "sql"

		var col Column
		col.Name = typ.Field(i).Name

		if typ.Field(i).Tag.Get(key) == "" {
			continue
		}

		attr := readTags(typ.Field(i).Tag.Get(key))

		tbl.Values[col.Name] = value.Field(i)

		if slice, typeOk := attr["type"]; typeOk {
			col.Types = strings.Join(slice, " ")
		}

		// TODO relationships and other types of columns

		tbl.Columns = append(tbl.Columns, col)
	}

	return tbl
}

func getTableName(typ reflect.Type) string {
	if typ == nil {
		return ""
	}

	n := strings.Split(typ.String(), ".")

	return strings.ToLower(strings2.Plural(n[len(n)-1]))
}
