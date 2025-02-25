/*
Copyright The Kubernetes Authors.

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

// Code generated by client-gen. DO NOT EDIT.

package v1

import (
	"context"
	"time"

	v1 "github.com/stolostron/integrity-shield/webhook/admission-controller/pkg/apis/manifestintegrityprofile/v1"
	scheme "github.com/stolostron/integrity-shield/webhook/admission-controller/pkg/client/manifestintegrityprofile/clientset/versioned/scheme"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	rest "k8s.io/client-go/rest"
)

// ManifestIntegrityProfilesGetter has a method to return a ManifestIntegrityProfileInterface.
// A group's client should implement this interface.
type ManifestIntegrityProfilesGetter interface {
	ManifestIntegrityProfiles() ManifestIntegrityProfileInterface
}

// ManifestIntegrityProfileInterface has methods to work with ManifestIntegrityProfile resources.
type ManifestIntegrityProfileInterface interface {
	Create(ctx context.Context, manifestIntegrityProfile *v1.ManifestIntegrityProfile, opts metav1.CreateOptions) (*v1.ManifestIntegrityProfile, error)
	Update(ctx context.Context, manifestIntegrityProfile *v1.ManifestIntegrityProfile, opts metav1.UpdateOptions) (*v1.ManifestIntegrityProfile, error)
	UpdateStatus(ctx context.Context, manifestIntegrityProfile *v1.ManifestIntegrityProfile, opts metav1.UpdateOptions) (*v1.ManifestIntegrityProfile, error)
	Delete(ctx context.Context, name string, opts metav1.DeleteOptions) error
	DeleteCollection(ctx context.Context, opts metav1.DeleteOptions, listOpts metav1.ListOptions) error
	Get(ctx context.Context, name string, opts metav1.GetOptions) (*v1.ManifestIntegrityProfile, error)
	List(ctx context.Context, opts metav1.ListOptions) (*v1.ManifestIntegrityProfileList, error)
	Watch(ctx context.Context, opts metav1.ListOptions) (watch.Interface, error)
	Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts metav1.PatchOptions, subresources ...string) (result *v1.ManifestIntegrityProfile, err error)
	ManifestIntegrityProfileExpansion
}

// manifestIntegrityProfiles implements ManifestIntegrityProfileInterface
type manifestIntegrityProfiles struct {
	client rest.Interface
}

// newManifestIntegrityProfiles returns a ManifestIntegrityProfiles
func newManifestIntegrityProfiles(c *ApisV1Client) *manifestIntegrityProfiles {
	return &manifestIntegrityProfiles{
		client: c.RESTClient(),
	}
}

// Get takes name of the manifestIntegrityProfile, and returns the corresponding manifestIntegrityProfile object, and an error if there is any.
func (c *manifestIntegrityProfiles) Get(ctx context.Context, name string, options metav1.GetOptions) (result *v1.ManifestIntegrityProfile, err error) {
	result = &v1.ManifestIntegrityProfile{}
	err = c.client.Get().
		Resource("manifestintegrityprofiles").
		Name(name).
		VersionedParams(&options, scheme.ParameterCodec).
		Do(ctx).
		Into(result)
	return
}

// List takes label and field selectors, and returns the list of ManifestIntegrityProfiles that match those selectors.
func (c *manifestIntegrityProfiles) List(ctx context.Context, opts metav1.ListOptions) (result *v1.ManifestIntegrityProfileList, err error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	result = &v1.ManifestIntegrityProfileList{}
	err = c.client.Get().
		Resource("manifestintegrityprofiles").
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Do(ctx).
		Into(result)
	return
}

// Watch returns a watch.Interface that watches the requested manifestIntegrityProfiles.
func (c *manifestIntegrityProfiles) Watch(ctx context.Context, opts metav1.ListOptions) (watch.Interface, error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	opts.Watch = true
	return c.client.Get().
		Resource("manifestintegrityprofiles").
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Watch(ctx)
}

// Create takes the representation of a manifestIntegrityProfile and creates it.  Returns the server's representation of the manifestIntegrityProfile, and an error, if there is any.
func (c *manifestIntegrityProfiles) Create(ctx context.Context, manifestIntegrityProfile *v1.ManifestIntegrityProfile, opts metav1.CreateOptions) (result *v1.ManifestIntegrityProfile, err error) {
	result = &v1.ManifestIntegrityProfile{}
	err = c.client.Post().
		Resource("manifestintegrityprofiles").
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(manifestIntegrityProfile).
		Do(ctx).
		Into(result)
	return
}

// Update takes the representation of a manifestIntegrityProfile and updates it. Returns the server's representation of the manifestIntegrityProfile, and an error, if there is any.
func (c *manifestIntegrityProfiles) Update(ctx context.Context, manifestIntegrityProfile *v1.ManifestIntegrityProfile, opts metav1.UpdateOptions) (result *v1.ManifestIntegrityProfile, err error) {
	result = &v1.ManifestIntegrityProfile{}
	err = c.client.Put().
		Resource("manifestintegrityprofiles").
		Name(manifestIntegrityProfile.Name).
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(manifestIntegrityProfile).
		Do(ctx).
		Into(result)
	return
}

// UpdateStatus was generated because the type contains a Status member.
// Add a +genclient:noStatus comment above the type to avoid generating UpdateStatus().
func (c *manifestIntegrityProfiles) UpdateStatus(ctx context.Context, manifestIntegrityProfile *v1.ManifestIntegrityProfile, opts metav1.UpdateOptions) (result *v1.ManifestIntegrityProfile, err error) {
	result = &v1.ManifestIntegrityProfile{}
	err = c.client.Put().
		Resource("manifestintegrityprofiles").
		Name(manifestIntegrityProfile.Name).
		SubResource("status").
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(manifestIntegrityProfile).
		Do(ctx).
		Into(result)
	return
}

// Delete takes name of the manifestIntegrityProfile and deletes it. Returns an error if one occurs.
func (c *manifestIntegrityProfiles) Delete(ctx context.Context, name string, opts metav1.DeleteOptions) error {
	return c.client.Delete().
		Resource("manifestintegrityprofiles").
		Name(name).
		Body(&opts).
		Do(ctx).
		Error()
}

// DeleteCollection deletes a collection of objects.
func (c *manifestIntegrityProfiles) DeleteCollection(ctx context.Context, opts metav1.DeleteOptions, listOpts metav1.ListOptions) error {
	var timeout time.Duration
	if listOpts.TimeoutSeconds != nil {
		timeout = time.Duration(*listOpts.TimeoutSeconds) * time.Second
	}
	return c.client.Delete().
		Resource("manifestintegrityprofiles").
		VersionedParams(&listOpts, scheme.ParameterCodec).
		Timeout(timeout).
		Body(&opts).
		Do(ctx).
		Error()
}

// Patch applies the patch and returns the patched manifestIntegrityProfile.
func (c *manifestIntegrityProfiles) Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts metav1.PatchOptions, subresources ...string) (result *v1.ManifestIntegrityProfile, err error) {
	result = &v1.ManifestIntegrityProfile{}
	err = c.client.Patch(pt).
		Resource("manifestintegrityprofiles").
		Name(name).
		SubResource(subresources...).
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(data).
		Do(ctx).
		Into(result)
	return
}
