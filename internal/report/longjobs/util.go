package longjobs

import (
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
)

// isFinished checks whether a job is finished. A job is considered
// finished when it reaches a terminal condition, either "Complete" or
// "Failed".
func isFinished(conds []batchv1.JobCondition) bool {
	for _, c := range conds {
		if c.Type == batchv1.JobComplete && c.Status == corev1.ConditionTrue {
			return true
		}
		if c.Type == batchv1.JobFailed && c.Status == corev1.ConditionTrue {
			return true
		}
	}
	return false
}

// isSuspended checks whether a job is suspended.
func isSuspended(conds []batchv1.JobCondition) bool {
	for _, c := range conds {
		if c.Type == batchv1.JobSuspended && c.Status == corev1.ConditionTrue {
			return true
		}
	}
	return false
}
