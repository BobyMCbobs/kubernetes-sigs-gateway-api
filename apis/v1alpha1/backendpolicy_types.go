/*
Copyright 2020 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +genclient
// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

// BackendPolicy defines policies associated with backends. For the purpose of
// this API, a backend is defined as any resource that a route can forward
// traffic to. A common example of a backend is a Service. Configuration that is
// implementation specific may be represented with similar implementation
// specific custom resources.
type BackendPolicy struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   BackendPolicySpec   `json:"spec,omitempty"`
	Status BackendPolicyStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// BackendPolicyList contains a list of BackendPolicy
type BackendPolicyList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []BackendPolicy `json:"items"`
}

// BackendPolicySpec defines desired policy for a backend.
type BackendPolicySpec struct {
	// BackendRefs define which backends this policy should be applied to. This
	// policy can only apply to backends within the same namespace. If more than
	// one BackendPolicy targets the same backend, precedence must be given to
	// the oldest BackendPolicy.
	//
	// Support: Core
	// +kubebuilder:validation:MaxItems=16
	BackendRefs []BackendRef `json:"backendRefs"`

	// TLS is the TLS configuration for these backends.
	//
	// Support: Extended
	// +optional
	TLS *BackendTLSConfig `json:"tls,omitempty"`
}

// BackendRef identifies an API object within a known namespace that defaults
// group to core and resource to services if unspecified.
type BackendRef struct {
	// Group is the group of the referent.  Omitting the value or specifying
	// the empty string indicates the core API group.  For example, use the
	// following to specify a service:
	//
	// fooRef:
	//   resource: services
	//   name: myservice
	//
	// Otherwise, if the core API group is not desired, specify the desired
	// group:
	//
	// fooRef:
	//   group: acme.io
	//   resource: foos
	//   name: myfoo
	//
	// +kubebuilder:default=core
	// +kubebuilder:validation:MaxLength=253
	Group string `json:"group,omitempty"`

	// Resource is the API resource name of the referent. Omitting the value
	// or specifying the empty string indicates the services resource. For example,
	// use the following to specify a services resource:
	//
	// fooRef:
	//   name: myservice
	//
	// Otherwise, if the services resource is not desired, specify the desired
	// group:
	//
	// fooRef:
	//   group: acme.io
	//   resource: foos
	//   name: myfoo
	//
	// +kubebuilder:default=services
	// +kubebuilder:validation:MaxLength=253
	Resource string `json:"resource,omitempty"`

	// Name is the name of the referent.
	//
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:MaxLength=253
	Name string `json:"name"`

	// Port is the port of the referent. If unspecified, this policy applies to
	// all ports on the backend.
	// +optional
	// +kubebuilder:validation:Minimum=1
	// +kubebuilder:validation:Maximum=65535
	Port *int32 `json:"port,omitempty"`
}

// BackendTLSConfig describes TLS configuration for a backend.
type BackendTLSConfig struct {
	// ClientCertificateRef is a reference to a TLS client certificate-key pair
	// that may be used to connect to these backends. If an entry in this list
	// omits or specifies the empty string for both the group and the resource,
	// the resource defaults to "secrets". An implementation may support other
	// resources (for example, resource "mycertificates" in group
	// "networking.acme.io").
	//
	// If a Secret is referenced, it must be of type "kubernetes.io/tls" and
	// contain tls.crt and tls.key data fields that contain the certificate and
	// private key to use for TLS.
	//
	// Support: Extended
	//
	// +optional
	ClientCertificateRef *CertificateObjectReference `json:"clientCertificateRef,omitempty"`

	// CertificateAuthorityRef is a reference to a resource that includes
	// trusted CA certificates for the associated backends. If an entry in this
	// list omits or specifies the empty string for both the group and the
	// resource, the resource defaults to "secrets". An implementation may
	// support other resources (for example, resource "mycertificates" in group
	// "networking.acme.io").
	//
	// When stored in a Secret, certificates must be PEM encoded and specified
	// within the "ca.crt" data field of the Secret. Multiple certificates can
	// be specified, concatenated by new lines.
	//
	// Support: Extended
	//
	// +optional
	CertificateAuthorityRef *CertificateObjectReference `json:"certificateAuthorityRef,omitempty"`

	// Options are a list of key/value pairs to give extended options to the
	// provider.
	//
	// Support: Implementation-specific.
	// +optional
	Options map[string]string `json:"options,omitempty"`
}

// BackendPolicyStatus defines the observed state of BackendPolicy. Conditions
// that are related to a specific Route or Gateway should be placed on the
// Route(s) using backends configured by this BackendPolicy.
type BackendPolicyStatus struct {
	// Conditions describe the current conditions of the BackendPolicy.
	// +optional
	// +kubebuilder:validation:MaxItems=8
	Conditions []metav1.Condition `json:"conditions,omitempty"`
}

// BackendPolicyConditionType is a type of condition associated with a
// BackendPolicy.
type BackendPolicyConditionType string

const (
	// ConditionNoSuchBackend indicates that one or more of the the specified
	// Backends does not exist.
	ConditionNoSuchBackend BackendPolicyConditionType = "NoSuchBackend"
)