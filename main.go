package main

import (
	"context"
	"database/sql"
	"log"

	"ariga.io/atlas/sql/postgres"
	"ariga.io/atlas/sql/schema"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
)

// This is on my machine, won't work for anyone else
// const CONN_STRING = "postgres://postgres:postgres@localhost/steampipe_host"

// This is from https://stepzen.com/blog/join-data-postgresql-declarative-graphql-without-sql
const CONN_STRING = "postgres://testUserIntrospection:HurricaneStartingSample1934@postgresql.introspection.stepzen.net/introspection"

const SCHEMA = "public"

func main() {
	ctx := context.Background()

	conn, err := sql.Open("pgx", CONN_STRING)
	if err != nil {
		log.Fatal(err)
	}

	driver, err := postgres.Open(conn)
	if err != nil {
		log.Fatal(err)
	}
	sch, err := driver.InspectSchema(ctx, SCHEMA, nil)
	if err != nil {
		log.Fatalf("failed inspecting schema: %s", err)
	}

	for _, table := range sch.Tables {
		var comment string
		for _, attr := range table.Attrs {
			if _attr, ok := attr.(*schema.Comment); ok {
				comment = _attr.Text
			}
		}

		log.Println("===", table.Name, "//", comment)

		for _, col := range table.Columns {
			var x proto.ColumnType
			switch col.Type.Type.(type) {
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
				x = proto.ColumnType_CIDR
			default:
				x = proto.ColumnType_UNKNOWN
			}
			log.Println("    -", col.Name, "->", col.Type.Raw, "=", x)
		}
	}
}
