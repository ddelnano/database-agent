package main

import (
	"context"
	"dagger/database-agent/internal/dagger"
)

type DatabaseAgent struct {
	Connection *dagger.Secret // +private
}

func New(connection *dagger.Secret) *DatabaseAgent { return &DatabaseAgent{Connection: connection} }

func (m *DatabaseAgent) Ask(ctx context.Context, question string) error {
	env := dag.Env().
		WithSQLInput("sql", dag.SQL(m.Connection), "The SQL module to use to ask questions").
		WithStringInput("question", question, "The question being asked about the database")

	_, err := dag.LLM().
		WithEnv(env).
		WithPrompt(`You are an expert database administrator. You have been given
a SQL module that already has tools with credentials and the ability to connect to the database to run SQL queries.

$question

Always show the SQL query you used to get the result.
DO NOT STOP UNTIL YOU HAVE ANSWERED THE QUESTION COMPLETELY.`).LastReply(ctx)

	return err
}
