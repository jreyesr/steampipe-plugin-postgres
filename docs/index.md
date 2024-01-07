---
organization: jreyesr
category: ["software development"]
icon_url: "/images/plugins/jreyesr/postgres.svg"
brand_color: "#336791"
display_name: Postgres
name: postgres
description: "Steampipe plugin for proxying queries to plain Postgres databases."
og_description: "Query any Postgres table from Steampipe with SQL! Open source CLI. No DB required."
og_image: "/images/plugins/jreyesr/postgres-social-graphic.png"
---

# Postgres + Steampipe

[PostgreSQL](https://www.postgresql.org/) is an open-source relational database, on which Steampipe is based.

[Steampipe](https://steampipe.io) is an open source CLI to instantly query cloud APIs using SQL.

The Postgres plugin for Steampipe lets you access arbitrary PostgreSQL databases from your Steampipe instance, much like a reverse proxy does
for HTTP. This lets you bring in static-ish data that resides on normal DBs, so you can join it with data pulled from APIs or other sources (which is Steampipe's specialty). For example, let's say that you have AWS EC2 instances, each of which has a tag that denotes the team that is responsible for it. You also have an internal DB where you have the contact details for each team, perhaps as some sort of on-call rotation system. With this plugin, you can `JOIN` across those datasources, detect some sort of noncompliance via AWS and automatically page the team that is responsible for said instance.

Steampipe currently has no such functionality, with an alternative being exporting a copy of the Postgres DB as a CSV and then using the [CSV plugin](https://hub.steampipe.io/plugins/turbot/csv). However, the CSV file may be out of date, and you're responsible for keeping it updated. This plugin, instead, will always have up-to-date results, since it queries the backing Postgres DB whenever a query comes in to Steampipe.

This plugin forwards all conditions that are supported by Steampipe to the remote DB. For example, a `WHERE col=1` condition _will_ be forwarded, so the remote DB can optimize its searches. More complex operators (such as JSONB operations) can't be forwarded and will thus result in a full table scan on the remote DB, and the filtering will be applied by Steampipe.

For example (using [the Sakila example database](https://github.com/jOOQ/sakila)):

```sql
select
  actor_id,
  first_name,
  last_name,
  last_update 
from
  postgres.actor limit 10;
```

## Documentation

- **[Table definitions & examples →](/plugins/jreyesr/postgres/tables)**

## Get started

### Install

Download and install the latest Postgres plugin:

```bash
steampipe plugin install jreyesr/postgres
```

### Credentials

You must provide a connection string, in the format [expected by `pgx`](https://pkg.go.dev/github.com/jackc/pgx/v5#hdr-Establishing_a_Connection). See [here](https://www.postgresql.org/docs/current/libpq-connect.html#LIBPQ-CONNSTRING)) for the official docs. For example, this is a valid connection string:

```
postgresql://user:pass@localhost/otherdb?connect_timeout=10&application_name=myapp
```

### Configuration

Installing the latest Postgres plugin will create a config file (`~/.steampipe/config/postgres.spc`) with a single connection named `postgres`:

```hcl
connection "postgres" {
  plugin = "jreyesr/postgres"

  # A connection string (https://pkg.go.dev/github.com/jackc/pgx/v5#hdr-Establishing_a_Connection), in the form that is 
  # expected by the pgx package. Required. 
  # Can also be set with the `DATABASE_URL` environment variable.
  # connection_string = "postgres://username:password@localhost:5432/database_name"

  # The remote DB's schema that this plugin will expose. If you leave this unset, it'll default to `public`.
  # schema = "public"

  # List of tables that will be exposed from the remote DB. No dynamic tables will be created if this arg is empty or not set.
  # Wildcard based searches are supported.
  # For example:
  #  - "*" will expose every table in the remote DB
  #  - "auth-*" will expose tables whose names start with "auth-"
  #  - "users" will only expose the specific table "users"
  # You can have several items (for example, ["auth-*", "users"] will expose 
  # all the tables that start with "auth-", PLUS the table "users")
  # Defaults to all tables
  # tables_to_expose = ["*"]
}
```

Alternatively, you can also use the following environment variable to obtain credentials **only if the other argument (`connection_string`)** is not specified in the connection:

```bash
export DATABASE_URL=postgres://username:password@localhost:5432/database_name
```

Uncomment and edit the `connection_string` parameter as described in the previous section. Alternatively, provide the `DATABASE_URL` envvar.

If the tables that you wish to expose don't live in the `public` schema on the remote DB, also uncomment and edit the `schema` parameter. If you don't provide it, it'll default to `public`.

## Get involved

- Open source: https://github.com/jreyesr/steampipe-plugin-postgres
- Community: [Join #steampipe on Slack →](https://turbot.com/community/join)
