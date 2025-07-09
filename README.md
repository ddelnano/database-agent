# Database Agent 🤖

**Ask your database questions and get answers!**

## Overview

This is an AI agent that connects to an existing database allows the user to ask plain langugae questions to explore and get information from a database. Supports both MySQL and PostgreSQL databases (determined by the connection string).

Built with [Dagger](https://dagger.io), the open platform for agentic software.

## Demo/Video

[![Build an AI Agent That Queries Your Database with Dagger](https://img.youtube.com/vi/LzHE0QTkQsM/maxresdefault.jpg)](https://www.youtube.com/watch?v=LzHE0QTkQsM)

## Features

- **Natural Language Interface**: Ask questions in plain language and get answers.
- **Database Support**: Connects to MySQL and PostgreSQL databases.
- **Interactive Exploration**: Ask questions about your database schema and data.
- **Developer experience**: Easy to use and integrate with existing databases.

## Installation

Install [Dagger](https://docs.dagger.io/install)

$ Clone the repo and enter into Dagger Shell:
```shell
git clone git@github.com:jasonmccallister/database-agent.git
```
```shell
cd database-agent
```
```shell
cp .env.example .env
# make sure you update the values in .env
```

> Note: The repo has a compose.yaml file that will start a MySQL and PostgreSQL database for you. You can use this to test the agent.

$ Start the example databases using Docker Compose:
```shell
docker compose up -d
```

$ Load the example dvdrental database into Postgres:
```shell
pg_restore -U postgres -h 127.0.0.1 -d postgres example-data/dvdrental.tar
```

> Note: when prompted for a password, use `postgres`.

$ Set your DB_URL environment variable to your database connection string.

```shell
export DB_URL=postgres://postgres:postgres@host.docker.internal:5432/postgres?sslmode=disable
```

$ Enter the Dagger Shell:
```shell
dagger
```

⋈ Ask the agent a question about your database!
```shell
ask "DB_URL=postgres://postgres:postgres@host.docker.internal:5432/postgres?sslmode=disable" "What tables do you have?"
```

⋈ Check out the answer.

*note: Increase verbosity to 2 or 3 and/or view in Dagger Cloud for best results*

### Fun to try:
- Try asking who the most popular actor is.
- Ask information about a specific table, like the actor table.
- Ask for the most rented movie.

Sample MCP commands

{"jsonrpc":"2.0","id":3,"method":"tools/call","params":{"name":"select_methods","arguments":{"methods":["DatabaseAgent.ask"]}}}

{"jsonrpc":"2.0","id":4,"method":"tools/call","params":{"name":"call_method","arguments":{"method":"DatabaseAgent.ask","self":"DatabaseAgent#1","args":{"dbUrl":"postgres://postgres:postgres@10.129.0.8:5432/postgres?sslmode=disable","question":"What tables are in this database?","uuid":"1"}}}}

### Local testing

```bash
# port 6060 is for pprof. Other one is for MCP
docker run -p 6060:6060 -p 8081:8080 -v /var/run/docker.sock:/var/run/docker.sock -e USE_LOCAL_DAGGER=true -e OPENAI_MODEL='gpt-4.1-nano-2025-04-14' -e OPENAI_API_KEY=<openai api key> -e DAGGER_CLOUD_TOKEN=<dagger cloud token> -it dagger-mcp-server
```

```bash
# port 6060 is for pprof. Other one is for MCP
docker run -p 6060:6060 -p 8081:8080 -v /var/run/docker.sock:/var/run/docker.sock -e USE_LOCAL_DAGGER=true -e OPENAI_MODEL='gpt-4.1-nano-2025-04-14' -e OPENAI_API_KEY=<openai api key> -e DAGGER_CLOUD_TOKEN=<dagger cloud token> -it dagger-mcp-server
```
### Local run, k8s dagger engine

```bash
DAGGER_ENGINE_POD_NAME="$(kubectl get pod \
    --selector=name=dagger-dagger-helm-engine --namespace=dagger \
    --output=jsonpath='{.items[0].metadata.name}')"
export DAGGER_ENGINE_POD_NAME

docker run -p 6060:6060 -p 8081:8080  -e OPENAI_MODEL='gpt-4.1-nano-2025-04-14' -e OPENAI_API_KEY=<open ai key> -e DAGGER_CLOUD_TOKEN=<cloud token> -it dagger-mcp-server
docker run -p 6060:6060 -p 8081:8080 -e DAGGER_ENGINE_POD_NAME=${DAGGER_ENGINE_POD_NAME} -e OPENAI_MODEL='gpt-4.1-nano-2025-04-14' -e OPENAI_API_KEY=<openai key> -e DAGGER_CLOUD_TOKEN=<dagger cloud token> -v $HOME/.kube:/root/.kube:ro -v $HOME/.aws:/root/.aws:ro  -it dagger-mcp-server

curl -X POST http://localhost:8081 -d '{"jsonrpc":"2.0","id":3,"method":"tools/call","params":{"name":"select_methods","arguments":{"methods":["DatabaseAgent.ask"]}}}'

curl -X POST http://localhost:8081 -d '{"jsonrpc":"2.0","id":4,"method":"tools/call","params":{"name":"call_method","arguments":{"method":"DatabaseAgent.ask","self":"DatabaseAgent#1","args":{"dbUrl":"postgres://postgres:postgres@postgres.default.svc.cluster.local:5432/postgres?sslmode=disable","question":"What tables are in this database?"}}}}'
```

### Fully k8s
```bash

curl -X POST http://localhost:8081 -d '{"jsonrpc":"2.0","id":3,"method":"tools/call","params":{"name":"select_methods","arguments":{"methods":["DatabaseAgent.ask"]}}}'

curl -X POST http://localhost:8081 -d '{"jsonrpc":"2.0","id":4,"method":"tools/call","params":{"name":"call_method","arguments":{"method":"DatabaseAgent.ask","self":"DatabaseAgent#1","args":{"dbUrl":"postgres://postgres:postgres@postgres.default.svc.cluster.local:5432/postgres?sslmode=disable","question":"What tables are in this database?"}}}}'
```
