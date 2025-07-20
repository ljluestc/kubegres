package resources

import (
	"fmt"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"path/filepath"
	kubegresv1 "reactive-tech.io/kubegres/api/v1"
)

// KubegresContext defines the context for Kubegres operations
type KubegresContext struct {
	Name      string
	Namespace string
	Kubegres  *kubegresv1.Kubegres
	Password  string
}

// GetNamespace returns the namespace
func (k *KubegresContext) GetNamespace() string {
	return k.Namespace
}

// GetName returns the name
func (k *KubegresContext) GetName() string {
	return k.Name
}

// GetDatabasePassword returns the database password
func (k *KubegresContext) GetDatabasePassword() string {
	return k.Password
}

// ReconcileResources handles reconciliation of resources
func ReconcileResources(ctx *KubegresContext) error {
	// Implementation placeholder
	return nil
}

func createPrimaryInitConfigMap(r *KubegresContext) corev1.ConfigMap {
	resource := r.Kubegres
	db := resource.Spec.Database

	pgDataPath := db.VolumeMount
	if db.Folder != "" {
		pgDataPath = filepath.Join(pgDataPath, db.Folder)
	}

	// Get the database password
	dbPassword := r.Password

	commands := fmt.Sprintf(`
			set -ex

			export POSTGRES_PASSWORD=%s
			export POSTGRES_USER=postgres
			export PGDATA=%s

			mkdir -p $PGDATA
			chown -R postgres:postgres $PGDATA

			[...]
`, dbPassword, pgDataPath)

	return corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: r.GetNamespace(),
			Name:      r.GetName() + "-init",
		},
		Data: map[string]string{
			"init.sh": commands,
		},
	}
}
