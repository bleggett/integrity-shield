// +build !ignore_autogenerated

//
// Copyright 2020 IBM Corporation
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//

// Code generated by operator-sdk. DO NOT EDIT.

package v1alpha1

import (
	resourcesigningprofilev1alpha1 "github.com/IBM/integrity-enforcer/enforcer/pkg/apis/resourcesigningprofile/v1alpha1"
	v1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	runtime "k8s.io/apimachinery/pkg/runtime"
	intstr "k8s.io/apimachinery/pkg/util/intstr"
)

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *CertPoolConfig) DeepCopyInto(out *CertPoolConfig) {
	*out = *in
	if in.KeyValue != nil {
		in, out := &in.KeyValue, &out.KeyValue
		*out = make([]byte, len(*in))
		copy(*out, *in)
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new CertPoolConfig.
func (in *CertPoolConfig) DeepCopy() *CertPoolConfig {
	if in == nil {
		return nil
	}
	out := new(CertPoolConfig)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *EsConfig) DeepCopyInto(out *EsConfig) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new EsConfig.
func (in *EsConfig) DeepCopy() *EsConfig {
	if in == nil {
		return nil
	}
	out := new(EsConfig)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *GlobalConfig) DeepCopyInto(out *GlobalConfig) {
	*out = *in
	if in.Arch != nil {
		in, out := &in.Arch, &out.Arch
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new GlobalConfig.
func (in *GlobalConfig) DeepCopy() *GlobalConfig {
	if in == nil {
		return nil
	}
	out := new(GlobalConfig)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *HttpConfig) DeepCopyInto(out *HttpConfig) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new HttpConfig.
func (in *HttpConfig) DeepCopy() *HttpConfig {
	if in == nil {
		return nil
	}
	out := new(HttpConfig)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *IntegrityEnforcer) DeepCopyInto(out *IntegrityEnforcer) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	out.Status = in.Status
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new IntegrityEnforcer.
func (in *IntegrityEnforcer) DeepCopy() *IntegrityEnforcer {
	if in == nil {
		return nil
	}
	out := new(IntegrityEnforcer)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *IntegrityEnforcer) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *IntegrityEnforcerList) DeepCopyInto(out *IntegrityEnforcerList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]IntegrityEnforcer, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new IntegrityEnforcerList.
func (in *IntegrityEnforcerList) DeepCopy() *IntegrityEnforcerList {
	if in == nil {
		return nil
	}
	out := new(IntegrityEnforcerList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *IntegrityEnforcerList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *IntegrityEnforcerSpec) DeepCopyInto(out *IntegrityEnforcerSpec) {
	*out = *in
	if in.MaxSurge != nil {
		in, out := &in.MaxSurge, &out.MaxSurge
		*out = new(intstr.IntOrString)
		**out = **in
	}
	if in.MaxUnavailable != nil {
		in, out := &in.MaxUnavailable, &out.MaxUnavailable
		*out = new(intstr.IntOrString)
		**out = **in
	}
	if in.ReplicaCount != nil {
		in, out := &in.ReplicaCount, &out.ReplicaCount
		*out = new(int32)
		**out = **in
	}
	if in.MetaLabels != nil {
		in, out := &in.MetaLabels, &out.MetaLabels
		*out = make(map[string]string, len(*in))
		for key, val := range *in {
			(*out)[key] = val
		}
	}
	if in.SelectorLabels != nil {
		in, out := &in.SelectorLabels, &out.SelectorLabels
		*out = make(map[string]string, len(*in))
		for key, val := range *in {
			(*out)[key] = val
		}
	}
	if in.NodeSelector != nil {
		in, out := &in.NodeSelector, &out.NodeSelector
		*out = make(map[string]string, len(*in))
		for key, val := range *in {
			(*out)[key] = val
		}
	}
	if in.Affinity != nil {
		in, out := &in.Affinity, &out.Affinity
		*out = new(v1.Affinity)
		(*in).DeepCopyInto(*out)
	}
	if in.Tolerations != nil {
		in, out := &in.Tolerations, &out.Tolerations
		*out = make([]v1.Toleration, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	if in.ImagePullSecrets != nil {
		in, out := &in.ImagePullSecrets, &out.ImagePullSecrets
		*out = make([]v1.LocalObjectReference, len(*in))
		copy(*out, *in)
	}
	in.Security.DeepCopyInto(&out.Security)
	in.KeyRing.DeepCopyInto(&out.KeyRing)
	in.CertPool.DeepCopyInto(&out.CertPool)
	in.Server.DeepCopyInto(&out.Server)
	in.Logger.DeepCopyInto(&out.Logger)
	in.RegKeySecret.DeepCopyInto(&out.RegKeySecret)
	in.GlobalConfig.DeepCopyInto(&out.GlobalConfig)
	if in.EnforcerConfig != nil {
		in, out := &in.EnforcerConfig, &out.EnforcerConfig
		*out = (*in).DeepCopy()
	}
	if in.SignPolicy != nil {
		in, out := &in.SignPolicy, &out.SignPolicy
		*out = (*in).DeepCopy()
	}
	if in.PrimaryRpp != nil {
		in, out := &in.PrimaryRpp, &out.PrimaryRpp
		*out = new(resourcesigningprofilev1alpha1.ResourceSigningProfileSpec)
		(*in).DeepCopyInto(*out)
	}
	if in.DefaultRpp != nil {
		in, out := &in.DefaultRpp, &out.DefaultRpp
		*out = new(resourcesigningprofilev1alpha1.ResourceSigningProfileSpec)
		(*in).DeepCopyInto(*out)
	}
	in.WebhookNamespacedResource.DeepCopyInto(&out.WebhookNamespacedResource)
	in.WebhookClusterResource.DeepCopyInto(&out.WebhookClusterResource)
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new IntegrityEnforcerSpec.
func (in *IntegrityEnforcerSpec) DeepCopy() *IntegrityEnforcerSpec {
	if in == nil {
		return nil
	}
	out := new(IntegrityEnforcerSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *IntegrityEnforcerStatus) DeepCopyInto(out *IntegrityEnforcerStatus) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new IntegrityEnforcerStatus.
func (in *IntegrityEnforcerStatus) DeepCopy() *IntegrityEnforcerStatus {
	if in == nil {
		return nil
	}
	out := new(IntegrityEnforcerStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *KeyRingConfig) DeepCopyInto(out *KeyRingConfig) {
	*out = *in
	if in.KeyValue != nil {
		in, out := &in.KeyValue, &out.KeyValue
		*out = make([]byte, len(*in))
		copy(*out, *in)
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new KeyRingConfig.
func (in *KeyRingConfig) DeepCopy() *KeyRingConfig {
	if in == nil {
		return nil
	}
	out := new(KeyRingConfig)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *LoggerContainer) DeepCopyInto(out *LoggerContainer) {
	*out = *in
	if in.SecurityContext != nil {
		in, out := &in.SecurityContext, &out.SecurityContext
		*out = new(v1.SecurityContext)
		(*in).DeepCopyInto(*out)
	}
	if in.HttpConfig != nil {
		in, out := &in.HttpConfig, &out.HttpConfig
		*out = new(HttpConfig)
		**out = **in
	}
	in.Resources.DeepCopyInto(&out.Resources)
	if in.EsConfig != nil {
		in, out := &in.EsConfig, &out.EsConfig
		*out = new(EsConfig)
		**out = **in
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new LoggerContainer.
func (in *LoggerContainer) DeepCopy() *LoggerContainer {
	if in == nil {
		return nil
	}
	out := new(LoggerContainer)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *RegKeySecret) DeepCopyInto(out *RegKeySecret) {
	*out = *in
	if in.Value != nil {
		in, out := &in.Value, &out.Value
		*out = make([]byte, len(*in))
		copy(*out, *in)
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new RegKeySecret.
func (in *RegKeySecret) DeepCopy() *RegKeySecret {
	if in == nil {
		return nil
	}
	out := new(RegKeySecret)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *SecurityConfig) DeepCopyInto(out *SecurityConfig) {
	*out = *in
	if in.PodSecurityContext != nil {
		in, out := &in.PodSecurityContext, &out.PodSecurityContext
		*out = new(v1.PodSecurityContext)
		(*in).DeepCopyInto(*out)
	}
	if in.IEAdminSubjects != nil {
		in, out := &in.IEAdminSubjects, &out.IEAdminSubjects
		*out = make([]rbacv1.Subject, len(*in))
		copy(*out, *in)
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new SecurityConfig.
func (in *SecurityConfig) DeepCopy() *SecurityConfig {
	if in == nil {
		return nil
	}
	out := new(SecurityConfig)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ServerContainer) DeepCopyInto(out *ServerContainer) {
	*out = *in
	if in.SecurityContext != nil {
		in, out := &in.SecurityContext, &out.SecurityContext
		*out = new(v1.SecurityContext)
		(*in).DeepCopyInto(*out)
	}
	in.Resources.DeepCopyInto(&out.Resources)
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ServerContainer.
func (in *ServerContainer) DeepCopy() *ServerContainer {
	if in == nil {
		return nil
	}
	out := new(ServerContainer)
	in.DeepCopyInto(out)
	return out
}
