package postgres

import (
	"context"
	"path"

	"github.com/turbot/go-kit/helpers"
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
	connectionString, err := config.GetConnectionString()
	if err != nil {
		return nil, err
	}
	schemaName := config.GetSchema()

	schema, err := GetAtlasSchemaForDBSchema(ctx, connectionString, schemaName)
	if err != nil {
		plugin.Logger(ctx).Error("postgres.PluginTables", "get_schema_error", err)
		return nil, err
	}

	temp_table_names := []string{} // this is to keep track of the tables that we've already added

	plugin.Logger(ctx).Debug("postgres.PluginTables", "tables", schema.Tables, "patterns", config.TablesToExpose)
	for _, pattern := range config.TablesToExpose {
		for _, tableAtlas := range schema.Tables {

			if helpers.StringSliceContains(temp_table_names, tableAtlas.Name) {
				continue // we've already handled it before
			} else if ok, _ := path.Match(pattern, tableAtlas.Name); !ok {
				plugin.Logger(ctx).Debug("postgres.PluginTables.noMatch", "pattern", pattern, "table", tableAtlas.Name)
				continue // pattern didn't match, don't do what follows
			}

			// here we're sure that pattern matched and it's the first time, so process this table
			temp_table_names = append(temp_table_names, tableAtlas.Name)

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
	}

	// Manually add the raw table (that one will always exist, in addition to an unknown number of dynamic tables)
	tables["raw"] = tableRawQuery(ctx, d.Connection)

	plugin.Logger(ctx).Debug("tfbridge.PluginTables.makeTables", "tables", tables)
	return tables, nil
}
