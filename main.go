package main

import (
	"context"
	"dagger/database-agent/internal/dagger"
)

type DatabaseAgent struct{}

// Ask the database agent a question and get a response
func (m *DatabaseAgent) Ask(
	ctx context.Context,
	// The database connection URL to use
	dbURL *dagger.Secret,
	// The question to ask the database agent
	question string,
) (string, error) {
	env := dag.Env().
		WithStringInput("question", question, "The question about the database being asked").
		WithSQLInput("sql", dag.SQL(dbURL), "The SQL module to use to inspect the database")

	return dag.LLM().
		WithEnv(env).
		WithPrompt(`You are an expert database administrator. You have been given
a SQL module that already has tools with credentials and the ability to connect to the database to run SQL queries.
Always show the SQL query you used to get the result.

The question is: $question

DO NOT STOP UNTIL YOU HAVE ANSWERED THE QUESTION COMPLETELY.`).
		LastReply(ctx)
}
