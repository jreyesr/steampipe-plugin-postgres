package postgres

import (
	"context"
	"fmt"

	"ariga.io/atlas/sql/schema"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

func tablePostgres(ctx context.Context, connection *plugin.Connection) (*plugin.Table, error) {
	tableAtlas := ctx.Value(keyTable).(*schema.Table)
	name := tableAtlas.Name

	return &plugin.Table{
		Name:        name,
		Description: FindCommentOnAttrs(tableAtlas.Attrs),
		List: &plugin.ListConfig{
			Hydrate:    ListTableRecords,
			KeyColumns: makeKeyColumns(ctx, tableAtlas),
		},
		Columns: makeColumns(ctx, tableAtlas),
	}, nil
}

func getMapKey(ctx context.Context, d *transform.TransformData) (interface{}, error) {
	asMap, ok := d.HydrateItem.(map[string]any)
	if !ok {
		plugin.Logger(ctx).Error("postgres.hydrate.getMapKey", "err", "invalid type", "expected", "map[string]any", "actual", fmt.Sprintf("%T", d.HydrateItem))
		return nil, fmt.Errorf("can't convert hydrate item %v to map", d.HydrateItem)
	}

	key := d.Param.(string)
	return asMap[key], nil
}

func makeColumns(ctx context.Context, tableAtlas *schema.Table) []*plugin.Column {
	columns := make([]*plugin.Column, 0, len(tableAtlas.Columns))

	for _, col := range tableAtlas.Columns {
		postgresType := PostgresColTypeToSteampipeColType(ctx, col)
		if postgresType == proto.ColumnType_UNKNOWN {
			plugin.Logger(ctx).Warn("postgres.makeColumns", "msg", "unknown type, skipping column!", "column", col.Name, "type", col.Type.Raw)
			continue
		}
		columns = append(columns, &plugin.Column{
			Name:        col.Name,
			Type:        postgresType,
			Description: FindCommentOnAttrs(col.Attrs),
			Transform:   transform.FromP(getMapKey, col.Name),
		})
	}

	return columns
}

func makeKeyColumns(ctx context.Context, tableAtlas *schema.Table) plugin.KeyColumnSlice {
	var all = make([]*plugin.KeyColumn, 0, len(tableAtlas.Columns))

	for _, col := range tableAtlas.Columns {
		postgresType := PostgresColTypeToSteampipeColType(ctx, col)
		if postgresType == proto.ColumnType_UNKNOWN {
			plugin.Logger(ctx).Warn("postgres.makeColumns", "msg", "unknown type, skipping column!", "column", col.Name, "type", col.Type.Raw)
			continue
		}
		all = append(all, &plugin.KeyColumn{
			Name:      col.Name,
			Operators: plugin.GetValidOperators(), // Everything is valid! Just reuse Steampipe's own "list of all operators that can be handled"
			Require:   plugin.Optional,
		})
	}

	plugin.Logger(ctx).Info("makeKeyColumns.done", "val", all)
	return all
}

func ListTableRecords(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	config := GetConfig(d.Connection)
	connectionString, err := config.GetConnectionString()
	if err != nil {
		return nil, err
	}
	schemaName := config.GetSchema()

	plugin.Logger(ctx).Debug("postgres.ListTable", "quals", d.Quals)
	plugin.Logger(ctx).Debug("postgres.ListTable", "schema", schemaName)

	results, err := MakeSQLQuery(ctx, connectionString, schemaName, d.Table.Name, d.Quals)
	if err != nil {
		return nil, err
	}

	for _, result := range results {
		d.StreamListItem(ctx, result)
	}

	return nil, nil
}
