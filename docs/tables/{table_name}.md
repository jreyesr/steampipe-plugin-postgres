# Table: {table_name}

Query data from Postgres tables. A table is automatically created to represent each
table found in the remote PostgreSQL database.

Every table that is added here will have the same columns as its backing table, except for complex types such as user-defined types, spatial/geographical types, intervals, and so on. See [here](https://steampipe.io/docs/develop/writing-plugins#column-data-types) for a list of types that can be expressed. For example, `INTEGER` and `SERIAL` columns will become `INTEGER` on Steampipe.

The queries that can be performed, therefore, depend on the structure of the backing table. For instance, if you've configured this plugin to point to a Postgres database that contains a copy of [the Sakila example database](https://github.com/jOOQ/sakila):

```sql
select 
    actor_id, first_name, last_name, last_update 
from 
    postgres.actor
where first_name ILIKE 'a%'
limit 10
```

All columns will have data types that match their backing Postgres types. Some translations are applied, such as treating `DECIMAL`, `REAL`, `DOUBLE PRECISION` and `MONEY` as `DOUBLE PRECISION` columns. If any columns in the backing table have types that can't be cleanly translated to Steampipe (such as [`TSVECTOR`](https://www.postgresql.org/docs/current/datatype-textsearch.html) or arrays), the plugin will fail.

## Examples

The following examples all use the Sakila database. However, note that this plugin, by its nature, should work with any Postgres database.

### Inspect the table structure

```bash
> .inspect postgres.actor
+-------------+--------------------------+-------------------------------------------------------+
| column      | type                     | description                                           |
+-------------+--------------------------+-------------------------------------------------------+
| _ctx        | jsonb                    | Steampipe context in JSON form, e.g. connection_name. |
| actor_id    | bigint                   |                                                       |
| first_name  | text                     |                                                       |
| last_name   | text                     |                                                       |
| last_update | timestamp with time zone |                                                       |
+-------------+--------------------------+-------------------------------------------------------+
```

### List all actors

```sql
select 
  * 
from 
  postgres.actor 
```

```
+----------+-------------+--------------+---------------------------+--------------------------------+
| actor_id | first_name  | last_name    | last_update               | _ctx                           |
+----------+-------------+--------------+---------------------------+--------------------------------+
| 2        | NICK        | WAHLBERG     | 2006-02-14T23:34:33-05:00 | {"connection_name":"postgres"} |
| 21       | KIRSTEN     | PALTROW      | 2006-02-14T23:34:33-05:00 | {"connection_name":"postgres"} |
| 3        | ED          | CHASE        | 2006-02-14T23:34:33-05:00 | {"connection_name":"postgres"} |
| 4        | JENNIFER    | DAVIS        | 2006-02-14T23:34:33-05:00 | {"connection_name":"postgres"} |
| 5        | JOHNNY      | LOLLOBRIGIDA | 2006-02-14T23:34:33-05:00 | {"connection_name":"postgres"} |
| 6        | BETTE       | NICHOLSON    | 2006-02-14T23:34:33-05:00 | {"connection_name":"postgres"} |
| 7        | GRACE       | MOSTEL       | 2006-02-14T23:34:33-05:00 | {"connection_name":"postgres"} |
+----------+-------------+--------------+---------------------------+--------------------------------+
```

### Push filters down to the backing table

Simple filters (such as `=`, `!=`, comparison operators, `LIKE`, regex matches and `IS (NOT) NULL`) will be
forwarded to the backing database. Prefer such operators if possible, as it'll reduce the amount of data that will be
transferred from the data source to Steampipe.

```sql
select 
  * 
from 
  postgres.actor 
where
  first_name LIKE 'A%'
```

```
+----------+------------+-------------+---------------------------+--------------------------------+
| actor_id | first_name | last_name   | last_update               | _ctx                           |
+----------+------------+-------------+---------------------------+--------------------------------+
| 29       | ALEC       | WAYNE       | 2006-02-14T23:34:33-05:00 | {"connection_name":"postgres"} |
| 34       | AUDREY     | OLIVIER     | 2006-02-14T23:34:33-05:00 | {"connection_name":"postgres"} |
| 49       | ANNE       | CRONYN      | 2006-02-14T23:34:33-05:00 | {"connection_name":"postgres"} |
| 65       | ANGELA     | HUDSON      | 2006-02-14T23:34:33-05:00 | {"connection_name":"postgres"} |
| 71       | ADAM       | GRANT       | 2006-02-14T23:34:33-05:00 | {"connection_name":"postgres"} |
| 76       | ANGELINA   | ASTAIRE     | 2006-02-14T23:34:33-05:00 | {"connection_name":"postgres"} |
| 165      | AL         | GARLAND     | 2006-02-14T23:34:33-05:00 | {"connection_name":"postgres"} |
| 125      | ALBERT     | NOLTE       | 2006-02-14T23:34:33-05:00 | {"connection_name":"postgres"} |
| 173      | ALAN       | DREYFUSS    | 2006-02-14T23:34:33-05:00 | {"connection_name":"postgres"} |
| 190      | AUDREY     | BAILEY      | 2006-02-14T23:34:33-05:00 | {"connection_name":"postgres"} |
| 132      | ADAM       | HOPPER      | 2006-02-14T23:34:33-05:00 | {"connection_name":"postgres"} |
| 144      | ANGELA     | WITHERSPOON | 2006-02-14T23:34:33-05:00 | {"connection_name":"postgres"} |
| 146      | ALBERT     | JOHANSSON   | 2006-02-14T23:34:33-05:00 | {"connection_name":"postgres"} |
+----------+------------+-------------+---------------------------+--------------------------------+
```
