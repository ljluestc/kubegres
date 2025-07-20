package resources

import (
	"fmt"
	"path"
	v1 "reactive-tech.io/kubegres/api/v1"
)

// CreateInitScript generates an initialization script for PostgreSQL data directory
// based on the Kubegres specification
func CreateInitScript(spec v1.KubegresSpec) string {
	pgDataPath := spec.Database.VolumeMount
	if spec.Database.Folder != "" {
		pgDataPath = path.Join(pgDataPath, spec.Database.Folder)
	}
	initScript := fmt.Sprintf(`
mkdir -p %s
chown -R postgres:postgres %s
`, pgDataPath, pgDataPath)
	return initScript
}
