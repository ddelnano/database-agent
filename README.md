# Database Agent 🤖

**Ask your database questions and get answers!**

## Overview

This is an AI agent that connects to an existing database allows the user to ask plain langugae questions to explore and get information from a database. Supports both MySQL and PostgreSQL databases (determined by the connection string).

Built with [Dagger](https://dagger.io), the open platform for agentic software.

## Demo

https://github.com/user-attachments/assets/d966a4a2-001c-4a5c-9e51-2c1e06ba9b95

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
ask env:DB_URL "What tables do you have?"
```

⋈ Check out the answer.

*note: Increase verbosity to 2 or 3 and/or view in Dagger Cloud for best results*

### Fun to try:
- Try asking who the most popular actor is.
- Ask information about a specific table, like the actor table.
- Ask for the most rented movie.
