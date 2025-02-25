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

package observer

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"time"

	gkmatch "github.com/open-policy-agent/gatekeeper/pkg/mutation/match"
	cosign "github.com/sigstore/cosign/cmd/cosign/cli"
	"github.com/sigstore/k8s-manifest-sigstore/pkg/k8smanifest"
	log "github.com/sirupsen/logrus"
	vrc "github.com/stolostron/integrity-shield/observer/pkg/apis/manifestintegritystate/v1"
	misclient "github.com/stolostron/integrity-shield/observer/pkg/client/manifestintegritystate/clientset/versioned/typed/manifestintegritystate/v1"
	midclient "github.com/stolostron/integrity-shield/reporter/pkg/client/manifestintegritydecision/clientset/versioned/typed/manifestintegritydecision/v1"
	"github.com/stolostron/integrity-shield/shield/pkg/config"
	kubeutil "github.com/stolostron/integrity-shield/shield/pkg/kubernetes"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
	kubeclient "k8s.io/client-go/kubernetes"
)

const timeFormat = "2006-01-02T15:04:05Z"

const exportDetailResult = "ENABLE_DETAIL_RESULT"
const detailResultConfigName = "OBSERVER_RESULT_CONFIG_NAME"
const detailResultConfigKey = "OBSERVER_RESULT_CONFIG_KEY"

const defaultKeyInConfigMap = "config.yaml"
const defaultPodNamespace = "integrity-shield-operator-system"
const defaultExportDetailResult = true
const defaultObserverResultDetailConfigName = "verify-result-detail"

const logLevelEnvKey = "LOG_LEVEL"
const k8sLogLevelEnvKey = "K8S_MANIFEST_SIGSTORE_LOG_LEVEL"

const VerifyResourceViolationLabel = "integrityshield.io/verifyResourceViolation"
const VerifyResourceIgnoredLabel = "integrityshield.io/verifyResourceIgnored"
const SignatureResourceLabel = "integrityshield.io/signatureResource"

var IgnoredKinds = []string{"Event", "Lease", "Endpoints", "TokenReview", "SubjectAccessReview", "SelfSubjectAccessReview", "LocalSubjectAccessReview"}

type Observer struct {
	APIResources     []groupResource
	Namespaces       []string
	DynamicClient    dynamic.Interface
	MidClient        *midclient.ApisV1Client
	MisClient        *misclient.ApisV1Client
	Clientset        *kubeclient.Clientset
	IShiledNamespace string
}

// Observer Result Detail
type VerifyResultDetail struct {
	Time                 string                            `json:"time"`
	Namespace            string                            `json:"namespace"`
	Name                 string                            `json:"name"`
	Kind                 string                            `json:"kind"`
	ApiGroup             string                            `json:"apiGroup"`
	ApiVersion           string                            `json:"apiVersion"`
	Error                bool                              `json:"error"`
	Message              string                            `json:"message"`
	Violation            bool                              `json:"violation"`
	VerifyResourceResult *k8smanifest.VerifyResourceResult `json:"verifyResourceResult"`
}
type ConstraintResult struct {
	ConstraintName  string               `json:"constraintName"`
	Violation       bool                 `json:"violation"`
	TotalViolations int                  `json:"totalViolations"`
	Results         []VerifyResultDetail `json:"results"`
	Constraint      ConstraintSpec       `json:"constraint"`
}

type ObservationDetailResults struct {
	Time              string             `json:"time"`
	ConstraintResults []ConstraintResult `json:"constraintResults"`
}

var logLevelMap = map[string]log.Level{
	"panic": log.PanicLevel,
	"fatal": log.FatalLevel,
	"error": log.ErrorLevel,
	"warn":  log.WarnLevel,
	"info":  log.InfoLevel,
	"debug": log.DebugLevel,
	"trace": log.TraceLevel,
}

func NewObserver() *Observer {
	insp := &Observer{}
	return insp
}

func (self *Observer) Init() error {
	log.Info("initialize observer.")
	kubeconf, _ := kubeutil.GetKubeConfig()

	var err error
	err = self.getAPIResources(kubeconf)
	if err != nil {
		return err
	}
	err = self.getNamespaces(kubeconf)
	if err != nil {
		return err
	}

	// set kubeclients
	dynamicClient, _ := dynamic.NewForConfig(kubeconf)
	self.DynamicClient = dynamicClient
	mieClient, _ := midclient.NewForConfig(kubeconf)
	self.MidClient = mieClient
	misClient, _ := misclient.NewForConfig(kubeconf)
	self.MisClient = misClient
	clientset, _ := kubeclient.NewForConfig(kubeconf)
	self.Clientset = clientset

	namespace := os.Getenv("POD_NAMESPACE")
	if namespace == "" {
		namespace = defaultPodNamespace
	}
	self.IShiledNamespace = namespace

	// log
	if os.Getenv("LOG_FORMAT") == "json" {
		log.SetFormatter(&log.JSONFormatter{TimestampFormat: time.RFC3339Nano})
	}
	logLevelStr := os.Getenv(logLevelEnvKey)
	if logLevelStr == "" {
		logLevelStr = "info"
	}
	logLevel, ok := logLevelMap[logLevelStr]
	if !ok {
		logLevel = log.InfoLevel
	}
	os.Setenv(k8sLogLevelEnvKey, logLevelStr)
	log.SetLevel(logLevel)

	log.Info("initialize cosign.")
	_ = cosign.Initialize()
	return nil
}

func (self *Observer) Run() {
	// load requestHandlerConfig
	rhconfig, err := config.LoadRequestHandlerConfig()
	if err != nil {
		log.Error("Failed to load RequestHandlerConfig; err: ", err.Error())
	}

	// reload all namespaces
	kubeconf, _ := kubeutil.GetKubeConfig()
	err = self.getNamespaces(kubeconf)
	if err != nil {
		log.Info("failed to update namespace list")
	}

	// load constraints
	constraints, err := self.loadConstraints()
	if err != nil {
		if err.Error() == "the server could not find the requested resource" {
			log.Info("no observation results")
			return
		} else {
			log.Error("Failed to load constraints; err: ", err.Error())
		}
	}

	// ObservationDetailResults
	var constraintResults []ConstraintResult
	for _, constraint := range constraints {
		constraintName := constraint.Parameters.ConstraintName
		log.Infof("Process new constraint %s ...", constraintName)
		admissionOnly := constraint.Parameters.Action.AdmissionOnly
		if admissionOnly {
			log.Info("Reporting observation result is disabled.")
		}
		var violations []vrc.VerifyResult
		var nonViolations []vrc.VerifyResult
		narrowedGVKList := self.getPossibleProtectedGVKs(constraint.Match)
		if narrowedGVKList == nil {
			log.Infof("No resources to validate in the constraint: %s ", constraintName)
			return
		}
		log.Debug("possible Protected GVKs: ", narrowedGVKList)
		// get all resources of extracted GVKs
		resources := []unstructured.Unstructured{}
		for _, gResource := range narrowedGVKList {
			tmpResources, _ := self.getAllResoucesByGroupResource(gResource, constraint.Match.LabelSelector)
			resources = append(resources, tmpResources...)
		}

		// check all resources by verifyResource
		ignoreFields := constraint.Parameters.IgnoreFields
		secrets := constraint.Parameters.KeyConfigs
		ignoreFields = append(ignoreFields, rhconfig.RequestFilterProfile.IgnoreFields...)
		skipObjects := rhconfig.RequestFilterProfile.SkipObjects
		skipObjects = append(skipObjects, constraint.Parameters.SkipObjects...)
		results := []VerifyResultDetail{}
		for _, resource := range resources {
			log.Debugf("Observe new resource; ns:%s, kind:%s, name:%s", resource.GetNamespace(), resource.GetKind(), resource.GetName())
			// check if signature resource
			signatureResource := isSignatureResource(resource)
			if signatureResource {
				result := VerifyResultDetail{
					Time:       time.Now().Format(timeFormat),
					Kind:       resource.GroupVersionKind().Kind,
					ApiGroup:   resource.GetObjectKind().GroupVersionKind().Group,
					ApiVersion: resource.GetObjectKind().GroupVersionKind().Version,
					Name:       resource.GetName(),
					Namespace:  resource.GetNamespace(),
					Message:    "this resource is signatureResource",
					Violation:  false,
				}
				results = append(results, result)
				continue
			}
			result := ObserveResource(resource, constraint.Parameters, ignoreFields, skipObjects, secrets)
			imgAllow, imgMsg := ObserveImage(resource, constraint.Parameters.ImageProfile)
			if !imgAllow {
				if !result.Violation {
					result.Violation = true
					result.Message = imgMsg
				} else {
					result.Message = fmt.Sprintf("%s, [Image]%s", result.Message, imgMsg)
				}
			}
			result = self.checkDecisionLog(constraintName, result)
			log.Debug("Verify result: ", result)
			results = append(results, result)
		}

		// prepare for manifest integrity state
		for _, res := range results {
			// simple result
			if res.Violation {
				vres := vrc.VerifyResult{
					Namespace:  res.Namespace,
					Name:       res.Name,
					Kind:       res.Kind,
					ApiGroup:   res.ApiGroup,
					ApiVersion: res.ApiVersion,
					Result:     res.Message,
				}
				violations = append(violations, vres)
			} else {
				vres := vrc.VerifyResult{
					Namespace:  res.Namespace,
					Name:       res.Name,
					Kind:       res.Kind,
					ApiGroup:   res.ApiGroup,
					ApiVersion: res.ApiVersion,
					Result:     res.Message,
				}
				if res.VerifyResourceResult != nil {
					vres.Signer = res.VerifyResourceResult.Signer
					vres.SigRef = res.VerifyResourceResult.SigRef
					vres.SignedTime = res.VerifyResourceResult.SignedTime
				}
				nonViolations = append(nonViolations, vres)
			}
			log.WithFields(log.Fields{
				"constraintName": constraintName,
				"violation":      res.Violation,
				"kind":           res.Kind,
				"name":           res.Name,
				"namespace":      res.Namespace,
			}).Info(res.Message)
		}
		// summarize results
		var violated bool
		if len(violations) != 0 {
			violated = true
		} else {
			violated = false
		}
		count := len(violations)

		vrr := vrc.ManifestIntegrityStateSpec{
			ConstraintName:  constraintName,
			Violation:       violated,
			TotalViolations: count,
			Violations:      violations,
			NonViolations:   nonViolations,
			ObservationTime: time.Now().Format(timeFormat),
		}

		// export VerifyResult
		_ = self.exportVerifyResult(vrr, violated, admissionOnly)
		// VerifyResultDetail
		cres := ConstraintResult{
			ConstraintName:  constraintName,
			Results:         results,
			Violation:       violated,
			TotalViolations: count,
			Constraint:      constraint,
		}
		constraintResults = append(constraintResults, cres)
	}

	// export ConstraintResult
	res := ObservationDetailResults{
		ConstraintResults: constraintResults,
		Time:              time.Now().Format(timeFormat),
	}
	_ = self.exportResultDetail(res)
}

func (self *Observer) checkDecisionLog(constraintName string, res VerifyResultDetail) VerifyResultDetail {
	// load manifest integrity decision
	mie, err := self.MidClient.ManifestIntegrityDecisions(self.IShiledNamespace).Get(context.Background(), constraintName, metav1.GetOptions{})
	if err != nil {
		return res
	}
	for _, ex := range mie.Spec.AdmissionResults {
		if ex.Namespace == res.Namespace && ex.Name == res.Name &&
			ex.Kind == res.Kind && ex.ApiGroup == res.ApiGroup && ex.ApiVersion == res.ApiVersion {
			if ex.Allow {
				res.Violation = false
				res.Message = fmt.Sprintf("Created by skipUser: %s", ex.UserName)
				log.Debug("Decision log found. Created by skipUser: ", res)
				return res
			}
		}
	}
	return res
}

func (self *Observer) exportVerifyResult(vrr vrc.ManifestIntegrityStateSpec, violated, admissionOnly bool) error {
	// label
	vrv := "false"
	vri := "false"
	if violated {
		vrv = "true"
	}
	if admissionOnly {
		vri = "true"
	}
	labels := map[string]string{
		VerifyResourceViolationLabel: vrv,
		VerifyResourceIgnoredLabel:   vri,
	}

	obj, err := self.MisClient.ManifestIntegrityStates(self.IShiledNamespace).Get(context.Background(), vrr.ConstraintName, metav1.GetOptions{})
	if err != nil || obj == nil {
		log.Infof("creating new ManifestIntegrityState resource %s ...", vrr.ConstraintName)
		newVRC := &vrc.ManifestIntegrityState{
			ObjectMeta: metav1.ObjectMeta{
				Name: vrr.ConstraintName,
			},
			Spec: vrr,
		}

		newVRC.Labels = labels
		_, err = self.MisClient.ManifestIntegrityStates(self.IShiledNamespace).Create(context.Background(), newVRC, metav1.CreateOptions{})
		if err != nil {
			log.Error("failed to create ManifestIntegrityStates:", err.Error())
			return err
		}
	} else {
		log.Infof("updating ManifestIntegrityStates resource %s ...", vrr.ConstraintName)
		obj.Spec = vrr
		obj.Labels = labels
		_, err = self.MisClient.ManifestIntegrityStates(self.IShiledNamespace).Update(context.Background(), obj, metav1.UpdateOptions{})
		if err != nil {
			log.Error("failed to update ManifestIntegrityStates:", err.Error())
			return err
		}
	}
	return nil
}

func (self *Observer) exportResultDetail(results ObservationDetailResults) error {
	exportStr := os.Getenv(exportDetailResult)
	export := defaultExportDetailResult
	if exportStr != "" {
		export, _ = strconv.ParseBool(exportStr)
	}
	if !export {
		return nil
	}

	if len(results.ConstraintResults) == 0 {
		log.Info("no observation results")
		return nil
	}
	namespace := os.Getenv("POD_NAMESPACE")
	if namespace == "" {
		namespace = defaultPodNamespace
	}
	configName := os.Getenv(detailResultConfigName)
	if configName == "" {
		configName = defaultObserverResultDetailConfigName
	}
	configKey := os.Getenv(detailResultConfigKey)
	if configKey == "" {
		configKey = defaultKeyInConfigMap
	}

	// load
	cm, err := self.Clientset.CoreV1().ConfigMaps(namespace).Get(context.Background(), configName, metav1.GetOptions{})
	if err != nil {
		// create
		log.Info("creating new configmap to store verify result...", configName)
		newcm := &v1.ConfigMap{
			ObjectMeta: metav1.ObjectMeta{
				Name: configName,
			},
		}
		resByte, _ := json.Marshal(results)
		newcm.Data = map[string]string{
			configKey: string(resByte),
		}
		_, err := self.Clientset.CoreV1().ConfigMaps(namespace).Create(context.Background(), newcm, metav1.CreateOptions{})
		if err != nil {
			log.Error("failed to create configmap", err.Error())
			return err
		}
		return nil
	} else {
		// update
		log.Info("updating configmap ...", configName)
		resByte, _ := json.Marshal(results)
		cm.Data = map[string]string{
			configKey: string(resByte),
		}
		_, err := self.Clientset.CoreV1().ConfigMaps(namespace).Update(context.Background(), cm, metav1.UpdateOptions{})
		if err != nil {
			log.Error("failed to update configmap", err.Error())
			return err
		}
	}
	return nil
}

//
// Constraint
//

type ConstraintSpec struct {
	Match      gkmatch.Match          `json:"match,omitempty"`
	Parameters config.ParameterObject `json:"parameters,omitempty"`
}

type Kinds struct {
	Kinds     []string `json:"kinds,omitempty"`
	ApiGroups []string `json:"apiGroups,omitempty"`
}

func (self *Observer) loadConstraints() ([]ConstraintSpec, error) {
	gvr := schema.GroupVersionResource{
		Group:    "constraints.gatekeeper.sh",
		Version:  "v1beta1",
		Resource: "manifestintegrityconstraint",
	}
	constraintList, err := self.DynamicClient.Resource(gvr).List(context.Background(), metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	micList := []ConstraintSpec{}
	for _, unstructed := range constraintList.Items {
		var mic ConstraintSpec
		spec, ok := unstructed.Object["spec"]
		if !ok {
			log.Error("failed to get spec in constraint", unstructed.GetName())
		}
		jsonStr, err := json.Marshal(spec)
		if err != nil {
			log.Error(err)
		}
		if err := json.Unmarshal(jsonStr, &mic); err != nil {
			log.Error(err)
		}
		micList = append(micList, mic)
	}
	return micList, nil
}

func isSignatureResource(resource unstructured.Unstructured) bool {
	var label bool
	if !(resource.GetKind() == "ConfigMap") {
		return label
	}
	labelsMap := resource.GetLabels()
	_, found := labelsMap[SignatureResourceLabel]
	if found {
		label = true
	}
	return label
}
