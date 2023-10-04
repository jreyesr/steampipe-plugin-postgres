# Postgres Plugin for Steampipe

Use SQL to query data from plain PostgreSQL databases.

This repo contains a [Steampipe](https://steampipe.io/) plugin that exposes plain PostgreSQL databases as Steampipe tables, much like [the CSV plugin](https://hub.steampipe.io/plugins/turbot/csv) does for CSV files, or like a reverse proxy does for HTTP. This can be used to join API data with semi-static data that is hosted on databases.

See below for an example that mixes data from a static DB (contact information for the teams that own Kubernetes namespaces) and data from Kubernetes (which namespaces have Failed pods). This may be part of an automated alerting system that runs periodically and sends emails.

![Alt text](docs/image.png)

Steampipe currently has no such functionality, with an alternative being exporting a copy of the Postgres DB as a CSV and then using the [CSV plugin](https://hub.steampipe.io/plugins/turbot/csv). However, the CSV file may be out of date, and you're responsible for keeping it updated. This plugin, instead, will always have up-to-date results, since it queries the backing Postgres DB whenever a query comes in to Steampipe. Another alternative is to manually install and configure [the `postgres-fdw` module on Steampipe](https://www.postgresql.org/docs/current/postgres-fdw.html), which requires connecting to the DB as root, is not documented and doesn't play too well with Dockerized deployments.

- **[Get started →](https://hub.steampipe.io/plugins/jreyesr/postgres)**
- Documentation: [Table definitions & examples](https://hub.steampipe.io/plugins/jreyesr/postgres/tables)
- Community: [Join #steampipe on Slack →](https://turbot.com/community/join)
- Get involved: [Issues](https://github.com/jreyesr/steampipe-plugin-postgres/issues)

## Quick start

Download and install the latest Postgres plugin:

```shell
steampipe plugin install jreyesr/postgres
```

Installing the latest Postgres plugin will create a config file (`~/.steampipe/config/postgres.spc`) with a single connection named `postgres`:

```hcl
connection "postgres" {
  plugin = "jreyesr/postgres"

  # A connection string (https://pkg.go.dev/github.com/jackc/pgx/v5#hdr-Establishing_a_Connection), in the form that is 
  expected by the pgx package. Required. 
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

Run Steampipe:

```shell
steampipe query
```

Run a query for whatever table the Postgres DB has:

```sql
select
  attr1,
  attr2
from
  postgres.some_table;
```

This plugin forwards all conditions that are supported by Steampipe to the remote DB. For example, a `WHERE col=1` condition _will_ be forwarded, so the remote DB can optimize its searches. More complex operators (such as JSONB operations) can't be forwarded and will thus result in a full table scan on the remote DB. In this case, the filtering will be applied by Steampipe.

## Developing

Prerequisites:

- [Steampipe](https://steampipe.io/downloads)
- [Golang](https://golang.org/doc/install)

Clone:

```sh
git clone https://github.com/jreyesr/steampipe-plugin-postgres.git
cd steampipe-plugin-postgres
```

Build, which automatically installs the new version to your `~/.steampipe/plugins` directory:

```
make
```

Configure the plugin:

```
cp config/* ~/.steampipe/config
vi ~/.steampipe/config/postgres.spc
```

Try it!

```
steampipe query
> .inspect postgres
```

Further reading:

- [Writing plugins](https://steampipe.io/docs/develop/writing-plugins)
- [Writing your first table](https://steampipe.io/docs/develop/writing-your-first-table)

## Contributing

Please see the [contribution guidelines](https://github.com/turbot/steampipe/blob/main/CONTRIBUTING.md) and our [code of conduct](https://github.com/turbot/steampipe/blob/main/CODE_OF_CONDUCT.md). All contributions are subject to the [Apache 2.0 open source license](https://github.com/jreyesr/steampipe-plugin-postgres/blob/master/LICENSE.md).

`help wanted` issues:

- [Steampipe](https://github.com/turbot/steampipe/labels/help%20wanted)
- [Postgres Plugin](https://github.com/jreyesr/steampipe-plugin-postgres/labels/help%20wanted)
