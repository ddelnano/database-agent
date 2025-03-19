package main

import (
	"context"
	"dagger/mysql-workspace/internal/dagger"
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

type MysqlWorkspace struct {
	// +private
	Conn *dagger.Secret
}

func New(conn *dagger.Secret) *MysqlWorkspace {
	return &MysqlWorkspace{
		Conn: conn,
	}
}

func (m *MysqlWorkspace) connect(ctx context.Context) (*sql.DB, string, error) {
	c, err := m.Conn.Plaintext(ctx)
	if err != nil {
		return nil, "", fmt.Errorf("error getting plaintext connection: %w", err)
	}

	db, err := sql.Open("mysql", c)
	if err != nil {
		return nil, "", fmt.Errorf("error opening database connection: %w", err)
	}

	return db, "", nil
}

// List the tables in a database in comma-separated format
func (m *MysqlWorkspace) ListTables(ctx context.Context) (string, error) {
	db, _, err := m.connect(ctx)
	if err != nil {
		return "", fmt.Errorf("error opening database connection: %w", err)
	}
	defer db.Close()

	rows, err := db.QueryContext(ctx, "show tables")
	if err != nil {
		return "", fmt.Errorf("error querying database: %w", err)
	}
	defer rows.Close()

	var tables []string
	for rows.Next() {
		var table string
		if err := rows.Scan(&table); err != nil {
			return "", fmt.Errorf("error scanning row: %w", err)
		}
		tables = append(tables, table)
	}

	return fmt.Sprintf("%s", tables), nil
}

// List columns in a table in comma-separated format
func (m *MysqlWorkspace) ListColumns(ctx context.Context, table string) (string, error) {
	db, _, err := m.connect(ctx)
	if err != nil {
		return "", fmt.Errorf("error opening database connection: %w", err)
	}
	defer db.Close()

	rows, err := db.QueryContext(ctx, fmt.Sprintf("show columns from %s", table))
	if err != nil {
		return "", fmt.Errorf("error querying database: %w", err)
	}
	defer rows.Close()

	var columns []string
	for rows.Next() {
		var column string
		if err := rows.Scan(&column); err != nil {
			return "", fmt.Errorf("error scanning row: %w", err)
		}
		columns = append(columns, column)
	}

	return fmt.Sprintf("%s", columns), nil
}

// List details on a specific column in a table in comma-separated format with the name,type, and nullability
func (m *MysqlWorkspace) ListColumnDetails(ctx context.Context, table, column string) (string, error) {
	db, _, err := m.connect(ctx)
	if err != nil {
		return "", fmt.Errorf("error opening database connection: %w", err)
	}
	defer db.Close()

	rows, err := db.QueryContext(ctx, fmt.Sprintf("show columns from %s where field = '%s'", table, column))
	if err != nil {
		return "", fmt.Errorf("error querying database: %w", err)
	}
	defer rows.Close()

	var name, typ, nullable string
	if rows.Next() {
		if err := rows.Scan(&name, &typ, &nullable); err != nil {
			return "", fmt.Errorf("error scanning row: %w", err)
		}
	}

	return fmt.Sprintf("%s,%s,%s", name, typ, nullable), nil
}

// Query the database with a custom query and return the results in comma-separated format
func (m *MysqlWorkspace) Query(ctx context.Context, query string) (string, error) {
	db, _, err := m.connect(ctx)
	if err != nil {
		return "", fmt.Errorf("error opening database connection: %w", err)
	}
	defer db.Close()

	rows, err := db.QueryContext(ctx, query)
	if err != nil {
		return "", fmt.Errorf("error querying database: %w", err)
	}
	defer rows.Close()

	columns, err := rows.Columns()
	if err != nil {
		return "", fmt.Errorf("error getting columns: %w", err)
	}

	var results []string
	for rows.Next() {
		columns := make([]any, len(columns))
		columnPointers := make([]any, len(columns))
		for i := range columns {
			columnPointers[i] = &columns[i]
		}
		if err := rows.Scan(columnPointers...); err != nil {
			return "", fmt.Errorf("error scanning row: %w", err)
		}
		results = append(results, fmt.Sprintf("%s", columns))
	}

	return fmt.Sprintf("%s", results), nil
}
