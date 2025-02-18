package schemaversion

import (
	"strconv"
)

const LATEST_VERSION = 41

type SchemaVersionMigrationFunc func(map[string]interface{}) error

type DataSourceInfo struct {
	Default bool
	UID     string
	Name    string
	ID      int64
}

type DataSourceInfoProvider interface {
	GetDataSourceInfo() []DataSourceInfo
}

func GetMigrations(dsInfoProvider DataSourceInfoProvider) map[int]SchemaVersionMigrationFunc {
	return map[int]SchemaVersionMigrationFunc{
		36: V36(dsInfoProvider),
		37: V37,
		38: V38,
		39: V39,
		40: V40,
		41: V41,
	}
}

func GetSchemaVersion(dash map[string]interface{}) int {
	if v, ok := dash["schemaVersion"]; ok {
		switch v := v.(type) {
		case int:
			return v
		case float64:
			return int(v)
		case string:
			if version, err := strconv.Atoi(v); err == nil {
				return version
			}
			return 0
		}
	}
	return 0
}
