package v1alpha1

import "github.com/redhat-appstudio/operator-toolkit/conditions"

const (
	// healthConditionType is the type used to track the health of a Foo resource
	healthConditionType conditions.ConditionType = "Health"
)

const (
	// HealthyReason is the reason set when the resource is healthy
	HealthyReason conditions.ConditionReason = "Healthy"

	// NotEnoughReplicasReason is the reason set when the resource needs to scale up
	NotEnoughReplicasReason conditions.ConditionReason = "NotEnoughReplicas"

	// TooManyReplicasReason is the reason set when the resource needs to scale down
	TooManyReplicasReason conditions.ConditionReason = "TooManyReplicas"
)
