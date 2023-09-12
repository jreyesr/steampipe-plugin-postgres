package postgres

import (
	"context"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

func tableRawQuery(ctx context.Context, connection *plugin.Connection) *plugin.Table {
	return &plugin.Table{
		Name:        "raw_query",
		Description: "Makes a raw SQL query (as a string) and returns any results as a single JSONB column. Use for more complex queries",
		List: &plugin.ListConfig{
			Hydrate:    ListRaw,
			KeyColumns: plugin.SingleColumn("query"),
		},
		Columns: []*plugin.Column{
			{Name: "query", Description: "The query that will be forwarded to the Postgres DB", Type: proto.ColumnType_STRING, Transform: transform.FromQual("query")},
			{Name: "data", Description: "The resultset, all wrapped in a JSONB column", Type: proto.ColumnType_JSON, Transform: transform.FromValue()},
		},
	}
}

func ListRaw(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	config := GetConfig(d.Connection)
	connectionString, err := config.GetConnectionString()
	if err != nil {
		return nil, err
	}
	schemaName := config.GetSchema()

	plugin.Logger(ctx).Debug("raw.ListRaw", "equalsQuals", d.EqualsQuals)
	plugin.Logger(ctx).Debug("raw.ListRaw", "schema", schemaName)

	results, err := MakeRawSQLQuery(ctx, connectionString, schemaName, d.Table.Name, d.EqualsQuals["query"].GetStringValue())
	if err != nil {
		return nil, err
	}

	for _, result := range results {
		d.StreamListItem(ctx, result)
	}

	return nil, nil
}
