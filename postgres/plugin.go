package postgres

import (
	"context"

	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

func Plugin(ctx context.Context) *plugin.Plugin {
	p := &plugin.Plugin{
		Name:             "steampipe-plugin-postgres",
		DefaultTransform: transform.FromGo().NullIfZero(),
		ConnectionConfigSchema: &plugin.ConnectionConfigSchema{
			NewInstance: ConfigInstance,
			Schema:      ConfigSchema,
		},
		SchemaMode:   plugin.SchemaModeDynamic,
		TableMapFunc: PluginTables,
	}
	return p
}

type key string

const (
	keyTable key = "table"
)

func PluginTables(ctx context.Context, d *plugin.TableMapData) (map[string]*plugin.Table, error) {
	tables := map[string]*plugin.Table{}

	config := GetConfig(d.Connection)
	schemaName := *config.Schema
	if schemaName == "" {
		schemaName = "public"
	}

	schema, err := GetAtlasSchemaForDBSchema(ctx, *config.ConnectionString, schemaName)
	if err != nil {
		plugin.Logger(ctx).Error("postgres.PluginTables", "get_schema_error", err)
		return nil, err
	}

	plugin.Logger(ctx).Debug("postgres.PluginTables", "tables", schema.Tables)
	for _, tableAtlas := range schema.Tables {
		// Pass the actual *schema.Table as a context key, as the CSV plugin does
		tableCtx := context.WithValue(ctx, keyTable, tableAtlas)

		tableSteampipe, err := tablePostgres(tableCtx, d.Connection)
		if err != nil {
			plugin.Logger(ctx).Error("postgres.PluginTables", "create_table_error", err, "tableName", tableAtlas.Name)
			return nil, err
		}

		plugin.Logger(ctx).Debug("postgres.PluginTables.makeTables", "table", tableSteampipe)
		tables[tableAtlas.Name] = tableSteampipe
	}
	plugin.Logger(ctx).Debug("tfbridge.PluginTables.makeTables", "tables", tables)

	return tables, nil
}
