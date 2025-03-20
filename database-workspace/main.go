package main

import (
	"context"
	"dagger/database-workspace/internal/dagger"
	"database/sql"
	"fmt"
	"net/url"
	"strings"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/jackc/pgx/v5/stdlib" // Import the pgx driver
)

type DatabaseWorkspace struct {
	Conn *dagger.Secret // +private
}

func New(conn *dagger.Secret) *DatabaseWorkspace {
	return &DatabaseWorkspace{
		Conn: conn,
	}
}

func (m *DatabaseWorkspace) connect(ctx context.Context) (*sql.DB, string, string, error) {
	c, err := m.Conn.Plaintext(ctx)
	if err != nil {
		return nil, "", "", fmt.Errorf("error getting plaintext connection: %w", err)
	}

	var (
		db     *sql.DB
		dbType string
	)
	conn := strings.ToLower(c)
	switch {
	case strings.HasPrefix(conn, "postgres://"), strings.HasPrefix(conn, "postgresql://"), strings.Contains(conn, "user=") && strings.Contains(conn, "dbname="):
		d, err := sql.Open("pgx", c)
		if err != nil {
			return nil, "", "", fmt.Errorf("error opening database connection: %w", err)
		}
		db = d
		dbType = "postgres"
	case strings.HasPrefix(conn, "mysql://"), strings.Contains(conn, "@tcp("), strings.Contains(conn, "user:") && strings.Contains(conn, "@/"):
		d, err := sql.Open("mysql", c)
		if err != nil {
			return nil, "", "", fmt.Errorf("error opening database connection: %w", err)
		}
		db = d
		dbType = "mysql"
	default:
		return nil, "", "", fmt.Errorf("unable to determine database type from connection string: %s", c)
	}

	u, err := url.Parse(c)
	if err != nil {
		return nil, "", "", fmt.Errorf("error parsing connection string: %w", err)
	}

	return db, dbType, strings.TrimPrefix(u.Path, "/"), nil
}

// List the tables in a database in comma-separated format
func (m *DatabaseWorkspace) ListTables(ctx context.Context,
	// +optional
	// +default="public"
	schema string) (string, error) {
	db, dbType, database, err := m.connect(ctx)
	if err != nil {
		return "", fmt.Errorf("error opening database connection: %w", err)
	}
	defer db.Close()

	var query string
	switch dbType {
	case "postgres":
		query = fmt.Sprintf("SELECT table_name FROM information_schema.tables WHERE table_schema = '%s' AND table_catalog = '%s'", schema, database)
	default:
		query = "show tables"
	}

	rows, err := db.QueryContext(ctx, query)
	if err != nil {
		return "", fmt.Errorf("error querying tables: %w", err)
	}
	defer rows.Close()

	var tables []string
	for rows.Next() {
		var tableName string
		if err := rows.Scan(&tableName); err != nil {
			return "", fmt.Errorf("error scanning row: %w", err)
		}
		tables = append(tables, tableName)
	}

	if err := rows.Err(); err != nil {
		return "", fmt.Errorf("error iterating rows: %w", err)
	}

	if len(tables) == 0 {
		return "", fmt.Errorf("no tables found, you might be in the wrong database or schema based on the connection")
	}

	return strings.Join(tables, ","), nil
}

// List the columns in a table in comma-separated format
func (m *DatabaseWorkspace) ListColumns(ctx context.Context, table string) (string, error) {
	db, _, database, err := m.connect(ctx)
	if err != nil {
		return "", fmt.Errorf("error opening database connection: %w", err)
	}
	defer db.Close()

	rows, err := db.QueryContext(ctx, fmt.Sprintf("SELECT column_name FROM information_schema.columns WHERE table_name = '%s' AND table_catalog = '%s'", table, database))
	if err != nil {
		return "", fmt.Errorf("error querying columns: %w", err)
	}
	defer rows.Close()

	var columns []string
	for rows.Next() {
		var columnName string
		if err := rows.Scan(&columnName); err != nil {
			return "", fmt.Errorf("error scanning row: %w", err)
		}
		columns = append(columns, columnName)
	}

	if err := rows.Err(); err != nil {
		return "", fmt.Errorf("error iterating rows: %w", err)
	}

	if len(columns) == 0 {
		return "", fmt.Errorf("no columns found, you might be in the wrong database or table based on the connection")
	}

	return strings.Join(columns, ","), nil
}

// List details on a specific column in a table in comma-separated format with the name,type, and nullability
func (m *DatabaseWorkspace) ListColumnDetails(ctx context.Context, table, column string) (string, error) {
	db, dbType, database, err := m.connect(ctx)
	if err != nil {
		return "", fmt.Errorf("error opening database connection: %w", err)
	}
	defer db.Close()

	query := fmt.Sprintf("SELECT column_name, data_type, is_nullable FROM information_schema.columns WHERE table_name = '%s' AND table_catalog = '%s' AND column_name = '%s'", table, database, column)
	if dbType == "mysql" {
		query = fmt.Sprintf("SELECT column_name, data_type, is_nullable FROM information_schema.columns WHERE table_name = '%s' AND table_schema = '%s' AND column_name = '%s'", table, database, column)
	}

	rows, err := db.QueryContext(ctx, query)
	if err != nil {
		return "", fmt.Errorf("error querying columns: %w", err)
	}
	defer rows.Close()

	var columnDetails []string
	for rows.Next() {
		var columnName, dataType, isNullable string
		if err := rows.Scan(&columnName, &dataType, &isNullable); err != nil {
			return "", fmt.Errorf("error scanning row: %w", err)
		}

		columnDetails = append(columnDetails, fmt.Sprintf("%s,%s,%s", columnName, dataType, isNullable))
	}
	if err := rows.Err(); err != nil {
		return "", fmt.Errorf("error iterating rows: %w", err)
	}

	if len(columnDetails) == 0 {
		return "", fmt.Errorf("no columns found, you might be in the wrong database or table based on the connection")
	}

	return strings.Join(columnDetails, ","), nil
}

// Query the database and return the results in comma-separated format
func (m *DatabaseWorkspace) RunQuery(ctx context.Context, query string) (string, error) {
	db, _, _, err := m.connect(ctx)
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
		values := make([]any, len(columns))
		valuePtrs := make([]any, len(columns))
		for i := range columns {
			valuePtrs[i] = &values[i]
		}
		if err := rows.Scan(valuePtrs...); err != nil {
			return "", fmt.Errorf("error scanning row: %w", err)
		}
		var row []string
		for _, value := range values {
			row = append(row, fmt.Sprintf("%v", value))
		}
		results = append(results, strings.Join(row, ","))
	}

	if err := rows.Err(); err != nil {
		return "", fmt.Errorf("error iterating rows: %w", err)
	}

	if len(results) == 0 {
		return "", fmt.Errorf("no results found")
	}

	return strings.Join(results, "\n"), nil
}
