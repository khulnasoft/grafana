package migration

import "github.com/grafana/grafana/pkg/apis/dashboard/migration/schemaversion"

// Initialize provides the migrator singleton with required dependencies and builds the map of migrations.
func Initialize(dsInfoProvider schemaversion.DataSourceInfoProvider) {
	migratorInstance.init(dsInfoProvider)
}

// Migrate migrates the given dashboard to the target version.
// This will block until the migrator is initialized.
func Migrate(dash map[string]interface{}, targetVersion int) error {
	return migratorInstance.migrate(dash, targetVersion)
}

var migratorInstance = &migrator{
	migrations: map[int]schemaversion.SchemaVersionMigrationFunc{},
	ready:      make(chan struct{}),
}

type migrator struct {
	ready      chan struct{}
	migrations map[int]schemaversion.SchemaVersionMigrationFunc
}

func (m *migrator) init(dsInfoProvider schemaversion.DataSourceInfoProvider) {
	m.migrations = schemaversion.GetMigrations(dsInfoProvider)
	close(m.ready)
}

func (m *migrator) migrate(dash map[string]interface{}, targetVersion int) error {
	if dash == nil {
		return schemaversion.NewMigrationError("dashboard is nil", 0, targetVersion)
	}

	// wait for the migrator to be initialized
	<-m.ready

	inputVersion := schemaversion.GetSchemaVersion(dash)
	dash["schemaVersion"] = inputVersion

	for nextVersion := inputVersion + 1; nextVersion <= targetVersion; nextVersion++ {
		if migration, ok := m.migrations[nextVersion]; ok {
			if err := migration(dash); err != nil {
				return schemaversion.NewMigrationError("migration failed", inputVersion, nextVersion)
			}
			dash["schemaVersion"] = nextVersion
		}
	}

	if schemaversion.GetSchemaVersion(dash) != targetVersion {
		return schemaversion.NewMigrationError("schema version not migrated to target version", inputVersion, targetVersion)
	}

	return nil
}
