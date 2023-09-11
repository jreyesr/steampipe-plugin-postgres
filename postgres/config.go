package postgres

import (
	"fmt"

	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/schema"
)

type PostgresConfig struct {
	ConnectionString *string `cty:"connection_string"`
	Schema           *string `cty:"schema"`
}

var ConfigSchema = map[string]*schema.Attribute{
	"connection_string": {Type: schema.TypeString},
	"schema":            {Type: schema.TypeString},
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
