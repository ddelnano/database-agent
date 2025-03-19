package main

import (
	"context"
	"dagger/database-agent/internal/dagger"
	"errors"
	"fmt"
)

type ConnectionType string

const (
	MYSQL    ConnectionType = "mysql"
	POSTGRES ConnectionType = "postgres"
)

var (
	Prompt = `You are an expert database administrator. You have been given
a workspace with the ability to connect to a database and run SQL queries and you have access to the following tools:

- list-tables
- list-columns
- list-coulmn-details
- run-query

Use ONLY the tools provided in the workspace to answer the question.

<question>
$question
</question>

Always show the SQL query you used to get the result.
DO NOT STOP UNTIL YOU HAVE ANSWERED THE QUESTION COMPLETELY.`
)

type DatabaseAgent struct {
	// +private
	ConnectionType ConnectionType
	// +private
	// +default="env:PG_URL"
	Connection *dagger.Secret
}

func New(connectionType ConnectionType, connection *dagger.Secret) *DatabaseAgent {
	return &DatabaseAgent{
		ConnectionType: connectionType,
		Connection:     connection,
	}
}

func (m *DatabaseAgent) Run(ctx context.Context, question string) (*dagger.Container, error) {
	switch m.ConnectionType {
	case POSTGRES:
		ws := dag.PostgresWorkspace(m.Connection)

		return dag.Llm().
			SetPostgresWorkspace("workspace", ws).
			WithPromptVar("question", question).
			WithPrompt(Prompt).
			Container(), nil
	case MYSQL:
		return nil, fmt.Errorf("mysql is not implemented")
	}

	return nil, errors.New("unknown connection type")
}
