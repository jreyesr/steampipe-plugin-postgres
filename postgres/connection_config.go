package postgres

import (
	"fmt"
	"os"

	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/schema"
)

type PostgresConfig struct {
	ConnectionString *string  `cty:"connection_string"`
	Schema           *string  `cty:"schema"`
	TablesToExpose   []string `cty:"tables_to_expose"`
}

var ConfigSchema = map[string]*schema.Attribute{
	"connection_string": {Type: schema.TypeString},
	"schema":            {Type: schema.TypeString},
	"tables_to_expose":  {Type: schema.TypeList, Elem: &schema.Attribute{Type: schema.TypeString}},
}

func ConfigInstance() interface{} {
	return &PostgresConfig{}
}

// GetConfig :: retrieve and cast connection config from query data
func GetConfig(connection *plugin.Connection) PostgresConfig {
	if connection == nil || connection.Config == nil {
		return PostgresConfig{}
	}
	config, _ := connection.Config.(PostgresConfig)
	return config
}

func (c PostgresConfig) String() string {
	return fmt.Sprintf(
		"PostgresConfig{schema=%s}",
		c.GetSchema()) // can't print connection_string, since it has credentials embedded
}

/*
GetSchema returns the schema that was configured in the .spc file, if available, and "public" otherwise
*/
func (c PostgresConfig) GetSchema() string {
	if c.Schema != nil && *c.Schema != "" {
		return *c.Schema
	}
	return "public"
}

func (c PostgresConfig) GetConnectionString() (string, error) {
	if c.ConnectionString != nil && *c.ConnectionString != "" {
		return *c.ConnectionString, nil
	}

	if v := os.Getenv("DATABASE_URL"); v != "" {
		return v, nil
	}

	return "", fmt.Errorf("please provide either the connection_string param or the DATABASE_URL envvar")
}

/*
GetTablesToExpose returns the slice of table blobs that was configured in the .spc file, if set, and ["*"] otherwise (which will expose every table)
*/
func (c PostgresConfig) GetTablesToExpose() []string {
	if len(c.TablesToExpose) > 0 {
		return c.TablesToExpose
	}
	return []string{"*"}
}
