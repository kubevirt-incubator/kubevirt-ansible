package v1

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	runtime "k8s.io/apimachinery/pkg/runtime"
)

// MyOperatorResource is an example operator configuration type
type MyOperatorResource struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata"`

	Spec   MyOperatorResourceSpec   `json:"spec"`
	Status MyOperatorResourceStatus `json:"status"`
}

type MyOperatorResourceSpec struct {
	OperatorSpec `json:",inline"`
}

type MyOperatorResourceStatus struct {
	OperatorStatus `json:",inline"`
}

type ManagementState string

var (
	// Force means that the operator is actively managing its resources but will not block an upgrade
	// if unmet prereqs exist. This state puts the operator at risk for unsuccessful upgrades
	Force ManagementState = "Force"
	// Managed means that the operator is actively managing its resources and trying to keep the component active.
	// It will only upgrade the component if it is safe to do so
	Managed ManagementState = "Managed"
	// Unmanaged means that the operator will not take any action related to the component
	Unmanaged ManagementState = "Unmanaged"
	// Removed means that the operator is actively managing its resources and trying to remove all traces of the component
	Removed ManagementState = "Removed"
)

// OperatorSpec contains common fields operators need.  It is intended to be anonymous included
// inside of the Spec struct for your particular operator.
type OperatorSpec struct {
	// managementState indicates whether and how the operator should manage the component
	ManagementState ManagementState `json:"managementState"`

	// logLevel is an intent based logging for an overall component.  It does not give fine grained control, but it is a
	// simple way to manage coarse grained logging choices that operators have to interpret for their operands.
	LogLevel LogLevel `json:"logLevel"`

	// operandSpecs provide customization for functional units within the component
	OperandSpecs []OperandSpec `json:"operandSpecs"`

	// unsupportedConfigOverrides holds a sparse config that will override any previously set options.  It only needs to be the fields to override
	// it will end up overlaying in the following order:
	// 1. hardcoded defaults
	// 2. observedConfig
	// 3. unsupportedConfigOverrides
	UnsupportedConfigOverrides runtime.RawExtension `json:"unsupportedConfigOverrides"`

	// observedConfig holds a sparse config that controller has observed from the cluster state.  It exists in spec because
	// it is an input to the level for the operator
	ObservedConfig runtime.RawExtension `json:"observedConfig"`
}

type LogLevel string

var (
	// Normal is the default.  Normal, working log information, everything is fine, but helpful notices for auditing or common operations.  In kube, this is probably glog=2.
	Normal LogLevel = "Normal"

	// Debug is used when something went wrong.  Even common operations may be logged, and less helpful but more quantity of notices.  In kube, this is probably glog=4.
	Debug LogLevel = "Debug"

	// Trace is used when something went really badly and even more verbose logs are needed.  Logging every function call as part of a common operation, to tracing execution of a query.  In kube, this is probably glog=6.
	Trace LogLevel = "Trace"

	// TraceAll is used when something is broken at the level of API content/decoding.  It will dump complete body content.  If you turn this on in a production cluster
	// prepare from serious performance issues and massive amounts of logs.  In kube, this is probably glog=8.
	TraceAll LogLevel = "TraceAll"
)

// ResourcePatch is a way to represent the patch you would issue to `kubectl patch` in the API
type ResourcePatch struct {
	// type is the type of patch to apply: jsonmerge, strategicmerge
	Type string `json:"type"`
	// patch the patch itself
	Patch string `json:"patch"`
}

// OperandSpec holds information for customization of a particular functional unit - logically maps to a workload
type OperandSpec struct {
	// name is the name of this unit.  The operator must be aware of it.
	Name string `json:"name"`

	// operandContainerSpecs are per-container options
	OperandContainerSpecs []OperandContainerSpec `json:"operandContainerSpecs"`

	// unsupportedResourcePatches are applied to the workload resource for this unit. This is an unsupported
	// workaround if anything needs to be modified on the workload that is not otherwise configurable.
	// TODO Decide: alternatively, we could simply include a RawExtension which is used in place of the "normal" default manifest
	UnsupportedResourcePatches []ResourcePatch `json:"unsupportedResourcePatches"`
}

type OperandContainerSpec struct {
	// name is the name of the container to modify
	Name string `json:"name"`

	// resources are the requests and limits to place in the container.  Nil means to accept the defaults.
	Resources *corev1.ResourceRequirements `json:"resources,omitempty"`
}

type OperatorStatus struct {
	// observedGeneration is the last generation change you've dealt with
	// +optional
	ObservedGeneration int64 `json:"observedGeneration,omitempty"`

	// conditions is a list of conditions and their status
	Conditions []OperatorCondition `json:"conditions,omitempty"`

	// version is the level this availability applies to
	Version string `json:"version"`

	// readyReplicas indicates how many replicas are ready and at the desired state
	ReadyReplicas int32 `json:"readyReplicas"`

	// generations are used to determine when an item needs to be reconciled or has changed in a way that needs a reaction.
	Generations []GenerationStatus `json:"generations"`
}

// GenerationStatus keeps track of the generation for a given resource so that decisions about forced updates can be made.
type GenerationStatus struct {
	// group is the group of the thing you're tracking
	Group string `json:"group"`
	// resource is the resource type of the thing you're tracking
	Resource string `json:"resource"`
	// namespace is where the thing you're tracking is
	Namespace string `json:"namespace"`
	// name is the name of the thing you're tracking
	Name string `json:"name"`
	// lastGeneration is the last generation of the workload controller involved
	LastGeneration int64 `json:"lastGeneration"`
	// hash is an optional field set for resources without generation that are content sensitive like secrets and configmaps
	Hash string `json:"hash"`
}

var (
	// Available indicates that the operand is present and accessible in the cluster
	OperatorStatusTypeAvailable = "Available"
	// Progressing indicates that the operator is trying to transition the operand to a different state
	OperatorStatusTypeProgressing = "Progressing"
	// Failing indicates that the operator (not the operand) is unable to fulfill the user intent
	OperatorStatusTypeFailing = "Failing"
	// PrereqsSatisfied indicates that the things this operator depends on are present and at levels compatible with the
	// current and desired states.
	OperatorStatusTypePrereqsSatisfied = "PrereqsSatisfied"
	// Upgradeable indicates that the operator configuration itself (not prereqs) can be auto-upgraded by the CVO
	OperatorStatusTypeUpgradeable = "Upgradeable"
)

// OperatorCondition is just the standard condition fields.
type OperatorCondition struct {
	Type               string          `json:"type"`
	Status             ConditionStatus `json:"status"`
	LastTransitionTime metav1.Time     `json:"lastTransitionTime,omitempty"`
	Reason             string          `json:"reason,omitempty"`
	Message            string          `json:"message,omitempty"`
}

type ConditionStatus string

const (
	ConditionTrue    ConditionStatus = "True"
	ConditionFalse   ConditionStatus = "False"
	ConditionUnknown ConditionStatus = "Unknown"
)

// StaticPodOperatorSpec is spec for controllers that manage static pods.
type StaticPodOperatorSpec struct {
	OperatorSpec           `json:",inline"`

	// failedRevisionLimit is the number of failed static pod installer revisions to keep on disk and in the api
	// -1 = unlimited, 0 or unset = 5 (default)
	FailedRevisionLimit    int32 `json:"failedRevisionLimit,omitempty"`
	// succeededRevisionLimit is the number of successful static pod installer revisions to keep on disk and in the api
	// -1 = unlimited, 0 or unset = 5 (default)
	SucceededRevisionLimit int32 `json:"succeededRevisionLimit,omitempty"`
}

// StaticPodOperatorStatus is status for controllers that manage static pods.  There are different needs because individual
// node status must be tracked.
type StaticPodOperatorStatus struct {
	OperatorStatus `json:",inline"`

	// latestAvailableRevision is the deploymentID of the most recent deployment
	LatestAvailableRevision int32 `json:"latestAvailableRevision"`

	// nodeStatuses track the deployment values and errors across individual nodes
	NodeStatuses []NodeStatus `json:"nodeStatuses"`
}

// NodeStatus provides information about the current state of a particular node managed by this operator.
type NodeStatus struct {
	// nodeName is the name of the node
	NodeName string `json:"nodeName"`

	// currentRevision is the generation of the most recently successful deployment
	CurrentRevision int32 `json:"currentRevision"`
	// targetRevision is the generation of the deployment we're trying to apply
	TargetRevision int32 `json:"targetRevision"`
	// lastFailedRevision is the generation of the deployment we tried and failed to deploy.
	LastFailedRevision int32 `json:"lastFailedRevision"`

	// lastFailedRevisionErrors is a list of the errors during the failed deployment referenced in lastFailedRevision
	LastFailedRevisionErrors []string `json:"lastFailedRevisionErrors"`
}
