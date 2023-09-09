package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"ariga.io/atlas/sql/postgres"
	"ariga.io/atlas/sql/schema"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
)

func connect(connectionString string) (*sql.DB, error) {
	return sql.Open("pgx", connectionString)
}

/*
GetAtlasSchemaForDBSchema gets the Atlas schema (as in, the metadata) for a Postgres schema (as in, the hierarchy below a database and above a table, such as `public`).
Must receive a connection string in the format expected by pgx (https://pkg.go.dev/github.com/jackc/pgx/v5#hdr-Establishing_a_Connection)
*/
func GetAtlasSchemaForDBSchema(ctx context.Context, connectionString, schema string) (*schema.Schema, error) {
	conn, err := connect(connectionString)
	if err != nil {
		return nil, fmt.Errorf("can't connect to DB: %w", err)
	}

	driver, err := postgres.Open(conn)
	if err != nil {
		return nil, fmt.Errorf("can't open Postgres driver: %w", err)
	}
	sch, err := driver.InspectSchema(ctx, schema, nil)
	if err != nil {
		return nil, fmt.Errorf("error inspecting schema: %w", err)
	}

	return sch, nil
}

/*
FindCommentOnAttrs tries to locate an Attr among the passed array that corresponds to a comment, and returns it if found.
Otherwise, returns an empty string.
This function can be used to identify the comment that is attached to a schema, table or column.
*/
func FindCommentOnAttrs(attrs []schema.Attr) string {
	var comment string
	for _, attr := range attrs {
		if _attr, ok := attr.(*schema.Comment); ok {
			comment = _attr.Text
		}
	}
	return comment
}

/*
PostgresColTypeToSteampipeColType converts an Atlas column type to a Steampipe column.
Atlas column types correspond almost one-to-one to actual SQL types, either standard SQL or Postgres extensions.
For example, DECIMAL, FLOAT and CURRENCY become DOUBLEs on Steampipe
*/
func PostgresColTypeToSteampipeColType(col *schema.Column) proto.ColumnType {
	var x proto.ColumnType

	switch t := col.Type.Type.(type) {
	case *schema.BinaryType, *postgres.BitType, *schema.EnumType, *schema.StringType, *schema.UUIDType:
		x = proto.ColumnType_STRING
	case *schema.BoolType:
		x = proto.ColumnType_BOOL
	case *schema.DecimalType, *schema.FloatType, *postgres.CurrencyType:
		x = proto.ColumnType_DOUBLE
	case *schema.IntegerType:
		x = proto.ColumnType_INT
	case *schema.JSONType:
		x = proto.ColumnType_JSON
	case *schema.TimeType, *postgres.IntervalType:
		x = proto.ColumnType_TIMESTAMP
	case *postgres.NetworkType:
		if t.T == "inet" {
			x = proto.ColumnType_INET
		} else if t.T == "cidr" {
			x = proto.ColumnType_CIDR
		} else {
			x = proto.ColumnType_UNKNOWN
		}
	default:
		// As of writing this, these are the types that fall here, AKA those that we don't know how to translate
		// *schema.SpatialType, *schema.UnsupportedType, *postgres.TextSearchType, *postgres.ArrayType, *postgres.SerialType, *postgres.OIDType, *postgres.RangeType, *postgres.UserDefinedType, *postgres.XMLType
		x = proto.ColumnType_UNKNOWN
	}

	return x
}

/*
Builds a slice to hold the columns of a single result row. Returns an array of pointers, that can be passed to DB.Scan()
*/
func prepareSliceForScanResults(columns []string) []any {
	arr := make([]any, len(columns))

	// Convert arr into an array of pointers, so we can save the results there
	for i := range arr {
		arr[i] = &arr[i]
	}

	return arr
}

func protoToPostgresValue(val *proto.QualValue) string {
	switch val.GetValue().(type) {
	case *proto.QualValue_BoolValue:
		return fmt.Sprintf("%t", val.GetBoolValue())
	case *proto.QualValue_DoubleValue:
		return fmt.Sprintf("%f", val.GetDoubleValue())
	case *proto.QualValue_InetValue:
		return fmt.Sprintf("'%s'", val.GetInetValue().GetCidr())
	case *proto.QualValue_Int64Value:
		return fmt.Sprintf("%d", val.GetInt64Value())
	case *proto.QualValue_JsonbValue:
		return fmt.Sprintf("'%s'", val.GetJsonbValue())
	case *proto.QualValue_StringValue:
		return fmt.Sprintf("'%s'", val.GetStringValue())
	case *proto.QualValue_TimestampValue:
		return fmt.Sprintf("'%s'", val.GetTimestampValue().AsTime().Format(time.RFC3339))
	default:
		return "<INVALID>" // this will probably cause an error on the query, which is OK
	}
}

/*
makeWhereConditions builds a string that contains the SQL conditions that match the passed quals
The string doesn't contain the WHERE keyword!
*/
func makeWhereConditions(quals plugin.KeyColumnQualMap) string {
	conds := make([]string, 0)

	for _, qualsForCol := range quals {
		for _, qual := range qualsForCol.Quals {
			if qual.Value.Value == nil {
				conds = append(conds, fmt.Sprintf("%s %s", qual.Column, qual.Operator))
			} else {
				conds = append(conds, fmt.Sprintf("%s %s %s", qual.Column, qual.Operator, protoToPostgresValue(qual.Value)))
			}
		}
	}

	return strings.Join(conds, " AND ")
}

/*
MakeSQLQuery sends a raw SQL query to a remote DB, and returns any results
*/
func MakeSQLQuery(ctx context.Context, connectionString, schema string, table string, quals plugin.KeyColumnQualMap) ([]map[string]any, error) {
	conn, err := connect(connectionString)
	if err != nil {
		return nil, fmt.Errorf("can't connect to DB: %w", err)
	}
	defer conn.Close()

	query := fmt.Sprintf("SELECT * FROM %s.%s", schema, table)
	whereConditions := makeWhereConditions(quals)
	if len(whereConditions) > 0 {
		query = query + " WHERE " + whereConditions
	}
	plugin.Logger(ctx).Debug("MakeSQLQuery.beforeExec", "query", query)
	rows, err := conn.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("error while making query \"%s\": %w", query, err)
	}
	defer rows.Close()

	colNames, err := rows.Columns()
	if err != nil {
		return nil, fmt.Errorf("error while reading column names: %w", err)
	}

	// The code here that stores results on a slice of map[string]any was inspired by https://lazypro.medium.com/make-sql-scan-result-be-map-in-golang-e04f0de5950f
	var results []map[string]any
	for rows.Next() {
		rowData := make(map[string]any)
		cols := prepareSliceForScanResults(colNames)

		if err := rows.Scan(cols...); err != nil {
			return nil, fmt.Errorf("error while reading columns: %w", err)
		}

		for i, v := range cols {
			rowData[colNames[i]] = v
		}
		plugin.Logger(ctx).Debug("Scan", "data", cols, "mapData", rowData)
		results = append(results, rowData)
	}

	// This must always be called after the for rows.Next() loop, since it may have terminated with an error
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error while reading columns: %w", err)
	}
	return results, nil
}
