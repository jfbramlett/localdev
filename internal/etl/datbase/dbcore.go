package datbase

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	_ "github.com/go-sql-driver/mysql"
	"github.com/pkg/errors"
)

type QueryResult struct {
	Cols []string
	Vals [][]interface{}
}

type Database interface {
	DisableForeignKeys(ctx context.Context) error
	EnableForeignKeys(ctx context.Context) error
	DeleteAll(ctx context.Context, table string) error
	Query(ctx context.Context, stmt string, args ...interface{}) (*QueryResult, error)
	QueryAll(ctx context.Context, table string) (*QueryResult, error)
	QueryIDs(ctx context.Context, query string) ([]interface{}, error)
	Insert(ctx context.Context, table string, cols []string, values [][]interface{}) error
	Exec(ctx context.Context, stmt string, args ...interface{}) error
}

func NewDatabase(dsn string) (Database, error) {
	conn, err := connect(dsn)
	if err != nil {
		return nil, err
	}

	return &database{conn: conn}, nil
}

func connect(dsn string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return db, nil
}

type database struct {
	conn *sql.DB
}

func (d *database) DisableForeignKeys(_ context.Context) error {
	_, err := d.conn.Exec("SET FOREIGN_KEY_CHECKS=0")
	if err != nil {
		return err
	}
	return nil
}
func (d *database) EnableForeignKeys(_ context.Context) error {
	_, err := d.conn.Exec("SET FOREIGN_KEY_CHECKS=1")
	if err != nil {
		return err
	}
	return nil
}

func (d *database) QueryAll(ctx context.Context, table string) (*QueryResult, error) {
	return d.Query(ctx, fmt.Sprintf("SELECT * from %s", table))
}

func (d *database) Query(_ context.Context, stmt string, args ...interface{}) (*QueryResult, error) {
	var rows *sql.Rows
	var err error
	if len(args) > 0 {
		rows, err = d.conn.Query(stmt, args...)
	} else {
		rows, err = d.conn.Query(stmt)
	}
	if err != nil {
		return nil, errors.WithStack(err)
	}

	result := &QueryResult{}

	cols, err := rows.Columns()
	if err != nil {
		return nil, errors.WithStack(err)
	}
	result.Cols = cols

	for rows.Next() {
		// Create a slice of interface{}'s to represent each column,
		// and a second slice to contain pointers to each item in the columns slice.
		columns := make([]interface{}, len(cols))
		columnPointers := make([]interface{}, len(cols))
		for i := 0; i < len(columns); i++ {
			columnPointers[i] = &columns[i]
		}

		// Scan the result into the column pointers...
		if err := rows.Scan(columnPointers...); err != nil {
			return nil, errors.WithStack(err)
		}

		// Create our map, and retrieve the value for each column from the pointers slice,
		// storing it in the map with the name of the column as the key.
		m := make([]interface{}, len(cols))
		for i := 0; i < len(cols); i++ {
			val := columnPointers[i].(*interface{})
			m[i] = val
		}
		result.Vals = append(result.Vals, m)
	}

	return result, nil
}

func (d *database) QueryIDs(ctx context.Context, query string) ([]interface{}, error) {
	queryResult, err := d.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	result := make([]interface{}, len(queryResult.Vals))
	for i, v := range queryResult.Vals {
		result[i] = v[0]
	}
	return result, nil
}

func (d *database) Insert(_ context.Context, table string, cols []string, vals [][]interface{}) error {
	if len(vals) == 0 {
		return nil
	}
	params := ValuesToParameters(vals[0])
	stmt := fmt.Sprintf("INSERT IGNORE INTO %s (%s) VALUES (%s)", table, strings.Join(cols, ","), params)
	for _, val := range vals {
		_, err := d.conn.Exec(stmt, val...)
		if err != nil {
			return errors.WithStack(err)
		}
	}

	return nil
}

func (d *database) Exec(_ context.Context, stmt string, args ...interface{}) error {
	if len(args) > 0 {
		_, err := d.conn.Exec(stmt, args...)
		return err
	}

	_, err := d.conn.Exec(stmt)
	return err
}

func (d *database) DeleteAll(_ context.Context, table string) error {
	stmt := fmt.Sprintf("DELETE FROM %s", table)
	_, err := d.conn.Exec(stmt)
	return err
}

func ValuesToParameters(vals []interface{}) string {
	args := make([]string, len(vals))
	for i := 0; i < len(vals); i++ {
		args[i] = "?"
	}
	return strings.Join(args, ",")
}
