# Table: {table_name}

Query data from Postgres tables. A table is automatically created to represent each
table found in the remote PostgreSQL database.

Every table that is added here will have the same columns as its backing table, except for complex types such as user-defined types, spatial/geographical types, intervals, and so on. See [here](https://steampipe.io/docs/develop/writing-plugins#column-data-types) for a list of types that can be expressed. For example, `INTEGER` and `SERIAL` columns will become `INTEGER` on Steampipe.

The queries that can be performed, therefore, depend on the structure of the backing table. For instance, if you've configured this plugin to point to a Postgres database that contains a copy of [the Sakila example database](https://github.com/jOOQ/sakila):

```sql
select
  actor_id,
  first_name,
  last_name,
  last_update 
from
  postgres.actor 
where
  first_name ILIKE 'a%' limit 10;
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
  actor_id,
  first_name,
  last_name,
  last_update 
from
  postgres.actor;
```

### Push filters down to the backing table

Simple filters (such as `=`, `!=`, comparison operators, `LIKE`, regex matches and `IS (NOT) NULL`) will be
forwarded to the backing database. Prefer such operators if possible, as it'll reduce the amount of data that will be
transferred from the data source to Steampipe. For example, here the `first_name LIKE 'A%'` will be pushed to the
remote DB, so only rows where the first name starts with A will be streamed back to Steampipe.

```sql
select
  actor_id,
  first_name,
  last_name,
  last_update 
from
  postgres.actor 
where
  first_name like 'A%';
```

If you use several conditions joined by `AND`, they'll be forwarded too:

```sql
select
  film_id,
  title,
  description,
  release_year,
  lang 
from
  postgres.film 
where
  release_year < 2000 AND description ILIKE '% epic %';
```

### Aggregations

Aggregations are run client-side (i.e. on Steampipe), so the remote DB will just see a query to list all the records:

```sql
select
  release_year,
  count(*) 
from
  postgres.film 
group by
  release_year;
```

The remote DB will only see `SELECT * FROM public.film`. The grouing and counting will be done by Steampipe. Keep this in mind if 
your query causes a lot of data to be fetched.

### Other operators

Similarly, other operators (such as JSONB operators, string operations or datetime operations) can't be forwarded to the remote DB:

```sql
select
  * 
from
  postgres.film 
where
  char_length(title) > 20;
```

In general, the following operations can be forwarded:

* For all column types: `=` and `!=`, `IS NULL` and `IS NOT NULL`
* For numbers: simple comparison operators, such as `>`, `<=` and its family
* For strings: `LIKE`, `ILIKE`, the regex matching operators such as `~`
* For JSONB: the `@>`, `<@`, `?`, `?|`, `?&`, `@?` and `@@` operators. See the "Additional JSONB operators" Table 9.46 [here](https://www.postgresql.org/docs/current/functions-json.html) for more information