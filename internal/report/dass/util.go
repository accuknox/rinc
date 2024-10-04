package dass

import (
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
)

// isDeploymentAvailable checks whether a deployment is available. A
// deployment is considered available when at least the minimum
// available replicas required are up and running for at least
// `minReadySeconds`.
func isDeploymentAvailable(conds []appsv1.DeploymentCondition) bool {
	for _, c := range conds {
		if c.Type != appsv1.DeploymentAvailable {
			continue
		}
		return c.Status == corev1.ConditionTrue
	}
	return false
}

// deploymentHasReplicaFailure checks whether a deployment has a replica
// failure. Replica failure occurs when one of the pods fails to be
// created or deleted.
func deploymentHasReplicaFailure(conds []appsv1.DeploymentCondition) bool {
	for _, c := range conds {
		if c.Type != appsv1.DeploymentReplicaFailure {
			continue
		}
		return c.Status == corev1.ConditionTrue
	}
	return false
}
