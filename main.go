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
	_, err := dag.LLM().
		WithSQL(dag.SQL(m.Connection)).
		WithPromptVar("question", question).
		WithPrompt(`You are an expert database administrator. You have been given
a SQL module with the ability to connect to a database and run SQL queries and you have access to the following tools:

- list-tables
- list-columns
- list-column-details
- run-query

Use the tools provided by the SQL module to answer the question.

<question>
$question
</question>

Always show the SQL query you used to get the result.
DO NOT STOP UNTIL YOU HAVE ANSWERED THE QUESTION COMPLETELY.`).LastReply(ctx)

	return err
}
