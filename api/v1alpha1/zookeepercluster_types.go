/*
 * Copyright 2020 Skulup Ltd, Open Collaborators
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *      http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package v1alpha1

import (
	"fmt"
	"github.com/skulup/operator-helper/basetype"
	"github.com/skulup/operator-helper/k8s"
	"github.com/skulup/operator-helper/operator/promethues"
	"github.com/skulup/operator-helper/reconciler"
	"github.com/skulup/zookeeper-operator/internal"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"strings"
)

var (
	_ reconciler.Defaulting = &ZookeeperCluster{}
)

const defaultRepository = "skulup/zookeeper"
const defaultTag = "latest"

const (
	defaultClusterSize = 3
	defaultDataDir     = "/data"
)
const (
	defaultAdminPort          = 8080
	defaultClientPort         = 2181
	defaultMetricsPort        = 7000
	defaultSecureClientPort   = -1
	defaultQuorumPort         = 2888
	defaultLeaderElectionPort = 3888
)

const (
	// VolumeReclaimPolicyDelete deletes the volume after the cluster is deleted
	VolumeReclaimPolicyDelete = "Delete"
	// VolumeReclaimPolicyRetain retains the volume after the cluster is deleted
	VolumeReclaimPolicyRetain = "Retain"
)

const (
	defaultStorageVolumeSize = "20Gi"
	defaultClusterDomain     = "cluster.local"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// ZookeeperClusterSpec defines the desired state of ZookeeperCluster
type ZookeeperClusterSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// +kubebuilder:validation:Required
	// +kubebuilder:validation:Minimum=1
	Size int32 `json:"size,omitempty"`

	Dirs *Dirs `json:"dirs,omitempty"`

	Ports *Ports `json:"ports,omitempty"`

	// Image defines the container image to use.
	Image basetype.Image `json:"image,omitempty"`

	// ZkCfg defines the zoo.cfg data
	ZkCfg string `json:"zkCfg,omitempty"`

	// Log4jProps defines the log4j.properties data
	Log4jProps string `json:"log4jProps,omitempty"`

	// Log4jQuietProps defines the log4j-quiet.properties data
	Log4jQuietProps string `json:"log4jQuietProps,omitempty"`

	PersistenceVolume *PersistenceVolume `json:"persistence,omitempty"`

	// PodConfig defines common configuration for the zookeeper pods
	PodConfig basetype.PodConfig `json:"pod,omitempty"`

	Metrics *promethues.MetricSpec `json:"metrics,omitempty"`

	Env []v1.EnvVar `json:"env,omitempty"`

	// Labels defines the labels to attach to the broker deployment
	Labels map[string]string `json:"labels,omitempty"`

	// Annotations defines the annotations to attach to the broker deployment
	Annotations map[string]string `json:"annotations,omitempty"`

	// ClusterDomain defines the cluster domain for the cluster
	// It defaults to cluster.local
	ClusterDomain string `json:"clusterDomain,omitempty"`
}

type Ports struct {
	Client       int32 `json:"client,omitempty"`
	SecureClient int32 `json:"secureClient,omitempty"`
	Metrics      int32 `json:"metrics,omitempty"`
	Quorum       int32 `json:"quorum,omitempty"`
	Leader       int32 `json:"leader,omitempty"`
	Admin        int32 `json:"admin,omitempty"`
}

func (in *Ports) setDefaults() (changed bool) {
	if in.Admin == 0 {
		changed = true
		in.Admin = defaultAdminPort
	}
	if in.Client == 0 {
		changed = true
		in.Client = defaultClientPort
	}
	if in.Metrics == 0 {
		changed = true
		in.Metrics = defaultMetricsPort
	}
	if in.SecureClient == 0 {
		changed = true
		in.SecureClient = defaultSecureClientPort
	}
	if in.Quorum == 0 {
		changed = true
		in.Quorum = defaultQuorumPort
	}
	if in.Leader == 0 {
		changed = true
		in.Leader = defaultLeaderElectionPort
	}
	return
}

type Dirs struct {
	Data string `json:"data,omitempty"`
	Log  string `json:"log,omitempty"`
}

// VolumeReclaimPolicy defines the possible volume reclaim policy: Delete or Retain
type VolumeReclaimPolicy string

// PersistenceVolume defines cluster node persistence volume is configured
type PersistenceVolume struct {
	// ReclaimPolicy decides the fate of the PVCs after the cluster is deleted.
	// If it's set to Delete and the zookeeper cluster is deleted, the corresponding PVCs will be deleted.
	// The default value is Retain.
	// +kubebuilder:validation:Enum="Delete";"Retain"
	ReclaimPolicy VolumeReclaimPolicy `json:"reclaimPolicy,omitempty"`
	// ClaimSpec describes the common attributes of storage devices
	// and allows a Source for provider-specific attributes
	ClaimSpec v1.PersistentVolumeClaimSpec `json:"claimSpec,omitempty"`
}

func (in *PersistenceVolume) setDefault() (changed bool) {
	if in.ReclaimPolicy != VolumeReclaimPolicyDelete && in.ReclaimPolicy != VolumeReclaimPolicyRetain {
		in.ReclaimPolicy = VolumeReclaimPolicyRetain
		changed = true
	}
	storage, ok := in.ClaimSpec.Resources.Requests[v1.ResourceStorage]
	if !ok || storage.IsZero() {
		changed = true
		if in.ClaimSpec.Resources.Requests == nil {
			in.ClaimSpec.Resources.Requests = map[v1.ResourceName]resource.Quantity{}
		}
		in.ClaimSpec.Resources.Requests[v1.ResourceStorage] = resource.MustParse(defaultStorageVolumeSize)
	}
	in.ClaimSpec.AccessModes = []v1.PersistentVolumeAccessMode{v1.ReadWriteOnce}
	return
}

func (in *ZookeeperCluster) setMetricsDefault() (changed bool) {
	return false
}

// ZookeeperClusterStatus defines the observed state of ZookeeperCluster
type ZookeeperClusterStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	Metadata Metadata `json:"metadata,omitempty"`
}

// Metadata defines the metadata status of the ZookeeperCluster
type Metadata struct {
	SizeZnodeCreated      bool    `json:"sizeZnodeCreated,omitempty"`
	ServiceMonitorVersion *string `json:"ServiceMonitorVersion,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

// ZookeeperCluster is the Schema for the zookeeperclusters API
type ZookeeperCluster struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ZookeeperClusterSpec   `json:"spec,omitempty"`
	Status ZookeeperClusterStatus `json:"status,omitempty"`
}

func (in *ZookeeperCluster) nameHasZkIndicator() bool {
	return strings.Contains(in.Name, "zk") || strings.Contains(in.Name, "zookeeper")
}

func (in *ZookeeperCluster) CreateLabels(addPodLabels bool, more map[string]string) map[string]string {
	labels := in.Spec.Labels
	if labels == nil {
		labels = map[string]string{}
	}
	if addPodLabels {
		for k, v := range in.Spec.PodConfig.Labels {
			labels[k] = v
		}
	}
	for k, v := range more {
		labels[k] = v
	}
	labels[k8s.LabelAppManagedBy] = internal.OperatorName
	labels[k8s.LabelAppName] = in.Name
	return labels
}

// ConfigMapName defines the name of the configmap object
func (in *ZookeeperCluster) ConfigMapName() string {
	if in.nameHasZkIndicator() {
		return in.Name
	}
	return fmt.Sprintf("%s-zk", in.GetName())
}

// StatefulSetName defines the name of the statefulset object
func (in *ZookeeperCluster) StatefulSetName() string {
	if in.nameHasZkIndicator() {
		return in.Name
	}
	return fmt.Sprintf("%s-zk", in.GetName())
}

// ClientServiceName defines the name of the client service object
func (in *ZookeeperCluster) ClientServiceName() string {
	if in.nameHasZkIndicator() {
		return fmt.Sprintf("%s", in.GetName())
	}
	return fmt.Sprintf("%s-zk", in.GetName())
}

// HeadlessServiceName defines the name of the headless service object
func (in *ZookeeperCluster) HeadlessServiceName() string {
	return fmt.Sprintf("%s-headless", in.ClientServiceName())
}

// ClientServiceFQDN defines the FQDN of the client service object
func (in *ZookeeperCluster) ClientServiceFQDN() string {
	return fmt.Sprintf("%s.%s.svc.%s", in.ClientServiceName(), in.Namespace, in.Spec.ClusterDomain)
}

// HeadlessServiceFQDN defines the FQDN of the headless service object
func (in *ZookeeperCluster) HeadlessServiceFQDN() string {
	return fmt.Sprintf("%s.%s.svc.%s", in.HeadlessServiceName(), in.Namespace, in.Spec.ClusterDomain)
}

// IsSslClientSupported returns whether SSL client is supported
func (in *ZookeeperCluster) IsSslClientSupported() bool {
	return in.Spec.Ports.SecureClient > 0
}

// SetSpecDefaults set the defaults for the cluster spec and returns true otherwise false
func (in *ZookeeperCluster) SetSpecDefaults() (changed bool) {
	if in.Spec.Image.SetDefaults(defaultRepository, defaultTag, v1.PullIfNotPresent) {
		changed = true
	}
	if in.Spec.Size == 0 {
		changed = true
		in.Spec.Size = defaultClusterSize
	}
	if in.Spec.ClusterDomain == "" {
		changed = true
		in.Spec.ClusterDomain = defaultClusterDomain
	}
	if in.Spec.Dirs == nil {
		changed = true
		in.Spec.Dirs = &Dirs{
			Data: defaultDataDir,
		}
	}
	if in.Spec.Ports == nil {
		in.Spec.Ports = &Ports{}
		in.Spec.Ports.setDefaults()
		changed = true
	} else if in.Spec.Ports.setDefaults() {
		changed = true
	}
	if in.Spec.PersistenceVolume == nil {
		in.Spec.PersistenceVolume = &PersistenceVolume{}
		in.Spec.PersistenceVolume.setDefault()
		changed = true
	} else if in.Spec.PersistenceVolume.setDefault() {
		changed = true
	}
	if in.setMetricsDefault() {
		changed = true
	}
	return
}

// SetStatusDefaults set the defaults for the cluster status and returns true otherwise false
func (in *ZookeeperCluster) SetStatusDefaults() (changed bool) {
	return
}

// +kubebuilder:object:root=true

// ZookeeperClusterList contains a list of ZookeeperCluster
type ZookeeperClusterList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []ZookeeperCluster `json:"items"`
}

func init() {
	SchemeBuilder.Register(&ZookeeperCluster{}, &ZookeeperClusterList{})
}
