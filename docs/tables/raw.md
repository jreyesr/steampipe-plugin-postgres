# Table: `raw`

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

### Make a `JOIN`

```sql
select 
    * 
from 
    postgres.raw 
where 
    query='SELECT * FROM film JOIN language USING (language_id) LIMIT 3'
```

```
+--------------------------------------------------------------+--------------------------------------------------------------------------------------------------------------->
| query                                                        | data                                                                                                          >
+--------------------------------------------------------------+--------------------------------------------------------------------------------------------------------------->
| SELECT * FROM film JOIN language USING (language_id) LIMIT 3 | {"description":"A Astounding Reflection of a Lumberjack And a Car who must Sink a Lumberjack in A Baloon Facto>
| SELECT * FROM film JOIN language USING (language_id) LIMIT 3 | {"description":"A Astounding Epistle of a Database Administrator And a Explorer who must Find a Car in Ancient>
| SELECT * FROM film JOIN language USING (language_id) LIMIT 3 | {"description":"A Epic Drama of a Feminist And a Mad Scientist who must Battle a Teacher in The Canadian Rocki>
+--------------------------------------------------------------+--------------------------------------------------------------------------------------------------------------->
```
