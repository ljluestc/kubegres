package resources

import (
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"path"
	kubegresv1 "reactive-tech.io/kubegres/api/v1"
)

// CreatePostgresPod creates a PostgreSQL pod with proper configuration
// based on the Kubegres specification
func CreatePostgresPod(spec kubegresv1.KubegresSpec) *v1.Pod {
	// Create a basic pod structure
	pod := &v1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "postgres-pod",
			Namespace: "default",
			Labels: map[string]string{
				"app": "postgres",
			},
		},
		Spec: v1.PodSpec{
			Containers: []v1.Container{
				{
					Name:  "postgres",
					Image: spec.Image,
				},
			},
		},
	}

	// Get the container for modifications
	container := &pod.Spec.Containers[0]

	// Set PGDATA environment variable with proper path
	pgDataPath := spec.Database.VolumeMount
	if spec.Database.Folder != "" {
		pgDataPath = path.Join(pgDataPath, spec.Database.Folder)
	}
	container.Env = append(container.Env, v1.EnvVar{Name: "PGDATA", Value: pgDataPath})

	return pod
}
