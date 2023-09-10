
See below for an example that mixes data from a static DB (contact information for the teams that own Kubernetes namespaces) and data from Kubernetes (which namespaces have Failed pods). This may be part of an automated alerting system that runs periodically and sends emails.

![Alt text](docs/image.png)

# Postgres Plugin for Steampipe

Use SQL to query data from plain PostgreSQL databases.

This repo contains a [Steampipe](https://steampipe.io/) plugin that exposes plain PostgreSQL databases as Steampipe tables, much like [the CSV plugin](https://hub.steampipe.io/plugins/turbot/csv) does for CSV files, or like a reverse proxy does for HTTP. This can be used to join API data with semi-static data that is hosted on databases.


- **[Get started →](https://hub.steampipe.io/plugins/jreyesr/postgres)**
- Documentation: [Table definitions & examples](https://hub.steampipe.io/plugins/jreyesr/postgres/tables)
- Community: [Join #steampipe on Slack →](https://turbot.com/community/join)
- Get involved: [Issues](https://github.com/jreyesr/steampipe-plugin-postgres/issues)

## Quick start

Install the plugin with [Steampipe](https://steampipe.io):

```shell
steampipe plugin install jreyesr/postgres
```

Configure your [config file](https://hub.steampipe.io/plugins/jreyesr/postgres#configuration) to point to a Postgres database, and optionally specify the schema too.

Run steampipe:

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