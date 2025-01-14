//
// Copyright 2021 IBM Corporation
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

package kubernetes

import (
	"fmt"

	_ "k8s.io/client-go/plugin/pkg/client/auth"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

func GetKubeConfig() (*rest.Config, error) {
	kubeCfg := defaultClientConfig()

	restConfig, err := kubeCfg.ClientConfig()
	if clientcmd.IsEmptyConfig(err) {
		restConfig, err := rest.InClusterConfig()
		if err != nil {
			return restConfig, fmt.Errorf("error creating REST client config in-cluster: %w", err)
		}

		return restConfig, nil
	}
	if err != nil {
		return restConfig, fmt.Errorf("error creating REST client config: %w", err)
	}

	return restConfig, nil
}

func defaultClientConfig() clientcmd.ClientConfig {
	loadingRules := clientcmd.NewDefaultClientConfigLoadingRules()
	configOverrides := &clientcmd.ConfigOverrides{}
	return clientcmd.NewNonInteractiveDeferredLoadingClientConfig(loadingRules, configOverrides)
}
