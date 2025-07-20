package resources

import (
	"path/filepath"

	corev1 "k8s.io/api/core/v1"
)

// Use the KubegresContext from resources package
func setPgDataEnvVars(r *KubegresContext, containerEnvVars []corev1.EnvVar) []corev1.EnvVar {
	pgDataPath := r.Kubegres.Spec.Database.VolumeMount
	if r.Kubegres.Spec.Database.Folder != "" {
		pgDataPath = filepath.Join(pgDataPath, r.Kubegres.Spec.Database.Folder)
	}
	containerEnvVars = append(containerEnvVars, corev1.EnvVar{
		Name:  "PGDATA",
		Value: pgDataPath,
	})
	return containerEnvVars
}
