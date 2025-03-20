# Database Agent 🤖

## Take a database connection and ask the agent questions about your database!

### What is this?

This module is an example agent that uses an AI Agent to connect to a database and answer questions about the database using human terms. The agent supports both MySQL and PostgreSQL databases (determined by the connection string).

https://github.com/user-attachments/assets/c7b8d869-f0d5-4423-9a91-fa5fe978d591

### How do I try it?

Start a dev Dagger Engine with LLM support using:
https://docs.dagger.io/ai-agents#initial-setup

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
set DB_URL=postgres://postgres:postgres@host.docker.internal:5432/postgres?sslmode=disable
```

$ Enter the Dagger Shell:
```shell
dagger
```

⋈ Ask the agent a question about your database!
```shell
. env:DB_URL | ask "What tables do you have?"
```

⋈ Check out the answer.

*note: Increase verbosity to 2 or 3 and/or view in Dagger Cloud for best results*

#### Fun to try:
- Try asking who the most popular actor is.
- Ask information about a specific table, like the actor table.
- Ask for the most rented movie.
