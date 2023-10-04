# Table: raw

Query data from Postgres tables, by sending a raw (string) SQL query to the backing database. Use this table
when you need to perform more complex operations (such as `JOIN`s across tables), when you can't query the separate
tables for each backing table.

Every result of the query will appear in a separate row.

## Examples

The following examples all use the Sakila database. However, note that this plugin, by its nature, should work with any Postgres database.

### Inspect the table structure

```bash
> .inspect postgres.raw
+--------+-------+-------------------------------------------------------+
| column | type  | description                                           |
+--------+-------+-------------------------------------------------------+
| _ctx   | jsonb | Steampipe context in JSON form, e.g. connection_name. |
| data   | jsonb | The resultset, all wrapped in a JSONB column          |
| query  | text  | The query that will be forwarded to the Postgres DB   |
+--------+-------+-------------------------------------------------------+
```

### Simple query

Send a raw query to the remote DB by passing the SQL statement(which should be a `SELECT`) on the `query` column:

```sql
select
  query,
  data 
from
  postgres.raw 
where
  query = 'SELECT * FROM film JOIN language USING (language_id) LIMIT 3' limit 2;
```

Every record that the remote table returns will be packed into the column `data`, as a JSONB record, and provided by Steampipe. 
You may want to use [JSONB functions](https://www.postgresql.org/docs/current/functions-json.html) to unpack the data back into 
SQL columns.

### Make a `JOIN`

```sql
select
  * 
from
  postgres.raw 
where
  query = 'SELECT * FROM film JOIN language USING (language_id) LIMIT 3';
```

The JOIN will be performed on the remote DB.

### Aggregations

The `raw` table is capable of pushing down aggregations to the remote datasource, thus reducing the amount of data
that has to be transferred back to Steampipe:

```sql
select
  * 
from
  postgres.raw 
where
  query = 'SELECT release_year, count(*) FROM film group by release_year';
```

This will precompute the number of films released per year in the remote DB, and will only transfer 
the counts per year from the DB to Steampipe. Consider using this pattern when operating on very 
large amounts of data, as otherwise you'd have to transfer all the data to Steampipe for aggregation.

### Extracting data from the results

The `raw` table returns its results in a single JSONB column, called `data`. To perform further operations, you may want to 
extract the fields of that JSONB object into columns using standard Postgres operations:

```sql
select
  data ->> 'title' as title,
  (
    data -> 'character_length'
  )
  ::int as title_length,
  data ->> 'description' as description 
from
  postgres.raw 
where
  query = 'SELECT title, character_length(title), description FROM film WHERE description ILIKE ''% epic %''';
```

This will return a table with three columns, `title`, `title_length` and `description`, instead of a single column `data` with a JSONB object.