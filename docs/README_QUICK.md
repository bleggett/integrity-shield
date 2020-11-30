# Quick Start

## Prerequisites
​
The following prerequisites must be satisfied to deploy Integrity Verifier on a cluster.
- A Kubernetes cluster and cluster admin access to the cluster to use `oc` or `kubectl` command
- Prepare a namespace to deploy Integrity Verifier. (We will use `integrity-verifier-operator-system` namespace in this document.)
- A secret resource (keyring-secret) which contains public key and certificates should be setup for enabling signature verification by Integrity Verifier.

---

## Install Integrity Verifier
​
This section describe the steps for deploying Integrity Verifier (IV) on your cluster. We will use RedHat OpenShift cluster and so use `oc` commands for installation. (You can use `kubectl` for Minikube or IBM Kubernetes Service.)

### Retrive the source from `integrity-enforcer` Git repository.

git clone this repository and moved to `integrity-enforcer` directory

```
$ git clone https://github.com/IBM/integrity-enforcer.git
$ cd integrity-verifier
$ pwd /home/repo/integrity-enforcer
```
In this document, we clone the code in `/home/repo/integrity-enforcer`.

### Prepape namespace for installing Integrity Verifier

You can deploy Integrity Verifier to any namespace. In this document, we will use `integrity-verifier-operator-system` to deploy Integrity Verifier.
```
make create-ns

```
We switch to `integrity-verifier-operator-system` namespace.
```
oc project integrity-verifier-operator-system
```
All the commands are executed on the `integrity-verifier-operator-system` namespace unless mentioned explicitly.

### Define public key secret in Integrity Verifier

Integrity Verifier requires a secret that includes a pubkey ring for verifying signatures of resources that need to be protected.  Integrity Verifier supports X509 or PGP key for signing resources. The following steps show how you can import your signature verification key to Integrity Verifier.

First, you need to export public key to a file. The following example shows a pubkey for a signer identified by an email `sample_signer@enterprise.com` is exported and stored in `/tmp/pubring.gpg`. (Use the filename `pubring.gpg`.)

```
$ gpg --export sample_signer@enterprise.com > /tmp/pubring.gpg
```

If you do not have any PGP key or you want to use new key, generate new one and export it to a file. See [this GitHub document](https://docs.github.com/en/free-pro-team@latest/github/authenticating-to-github/generating-a-new-gpg-key).

Then, create a secret that includes a pubkey ring for verifying signatures of resources

```
oc create secret generic --save-config keyring-secret  -n integrity-verifier-operator-system --from-file=/tmp/pubring.gpg
```

### Define signers for each namespace


You can define signer who can provide signature for resources on each namespace. It can be configured when deploying the Integrity Verifier. For that, configure signPolicy in the following Integrity Verifier Custom Resource [file](../operator/config/samples/apis.integrityverifier.io_v1alpha1_integrityverifier_cr.yaml). Example below shows a signer `SampleSigner` identified by email `sample_signer@enterprise.com` is configured to sign rosources to be protected in any namespace.

```yaml
# Edit integrity-verifier-operator/config/samples/apis_v1alpha1_integrityverifier.yaml
apiVersion: apis.integrityverifier.io/v1alpha1
kind: IntegrityVerifier
metadata:
  name: integrity-verifier-server
spec:
  namespace: integrity-verifier-operator-system
  verifierConfig:
    verifyType: pgp # x509
    inScopeNamespaceSelector:
      include:
      - "*"
      exclude:
      - "kube-*"
      - "openshift-*"
  signPolicy:
    policies:
    - namespaces:
      - "*"
      signers:
      - "SampleSigner"
    - scope: "Cluster"
      signers:
      - "SampleSigner"
    signers:
    - name: "SampleSigner"
      secret: keyring-secret
      subjects:
      - email: "signer@enterprise.com"
  keyRingConfigs:
  - name: keyring-secret


```



### Install Integrity Verifier to a cluster

Integrity Verifier can be installed to a cluster using a series of steps which are bundled in a script called [`install_verifier.sh`](../scripts/install_verifier.sh). Before executing the script `install_verifier.sh`, setup local environment as follows:
- `IV_REPO_ROOT=<set absolute path of the root directory of cloned integrity-verifier source repository`
- `KUBECONFIG=~/kube/config/minikube`  (for deploying IV on minikube cluster)

`~/kube/config/minikube` is the Kuebernetes config file with credentials for accessing a cluster via `kubectl`.

Example:
```
$ export KUBECONFIG=~/kube/config/minikube
$ export IV_REPO_ROOT=/home/repo/integrity-enforcer
```

Execute the following make commands to build Integrity Verifier images.
```
$ cd integrity-verifier
$ make build-images
$ make tag-images-to-local
```

Then, execute the following script to deploy Integrity Verifier in a cluster.

```
$ make install-crds
$ make install-operator
```

After successful installation, you should see two pods are running in the namespace `integrity-verifier-operator-system`.

```
$ oc get pod -n integrity-verifier-operator-system
integrity-verifier-operator-c4699c95c-4p8wp   1/1     Running   0          5m
integrity-verifier-server-85c787bf8c-h5bnj    2/2     Running   0          82m
```

---

## Protect Resources with Integrity Verifier
​
Once Integrity Verifier is deployed to a cluster, you are ready to put resources on the cluster into signature-based protection. To start actual protection, you need to define which resources should be protected specifically. This section describes the execution flow for protecting a specific resource (e.g. ConfigMap) in a specific namespace (e.g. `secure-ns`) on your cluster.

The steps for protecting resources include:
- Define which reource(s) should be protected.
- Create a resource with signature.

### Define which reource(s) should be protected

You can define which resources should be protected with signature in a cluster by Integrity Verifier. A custom resource `ResourceSigningProfile` (RSP) includes the definition and it is created in the same namespace as resources. Example below illustrates how to define RSP to protect three resources ConfigMap, Deployment, and Service in a namespace `secure-ns`. After this, any resources specified here cannot be created/updated without valid signature.

```
$ cat <<EOF | oc apply -n secure-ns -f -
apiVersion: apis.integrityverifier.io/v1alpha1
kind: ResourceSigningProfile
metadata:
  name: sample-rsp
spec:
  rules:
  - match:
    - kind: ConfigMap
    - kind: Deployment
    - kind: Service
EOF

resourcesigningprofile.apis.integrityverifier.io/sample-rsp created
```

See [Define Protected Resources](README_FOR_RESOURCE_PROTECTION_PROFILE.md) for detail specs.

### Create a resource with signature

Any configmap cannot be created without signature in `secure-ns` namespace. Run the following command to create a sample configmap.

```
cat << EOF > /tmp/test-cm.yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: test-cm
data:
  key1: val1
  key2: val2
  key4: val4
EOF
```


run the command below for trying to create the configmap in `secure-ns` namespace without signature. You will see it is blocked because no signature for this resource is stored in the cluster.

```
$ oc apply -f /tmp/test-cm.yaml -n secure-ns
Error from server: error when creating "test-cm.yaml": admission webhook "ac-server.integrity-verifier-operator-system.svc" denied the request: No signature found
```


To generate a signature for a resource, you can use a [utility script](../scripts/gpg-rs-sign.sh) (Use [yq](https://github.com/mikefarah/yq) in the script)

Run the following script to generate a signature

```
$ ./scripts/gpg-rs-sign.sh signer@enterprise.com /tmp/test-cm.yaml /tmp/test-cm-rs.yaml
```


Then, output file `/tmp/test-cm-rs.yaml` is A custom resource `ResourceSignature` which includes signature of the input yaml.


```yaml
apiVersion: apis.integrityverifier.io/v1alpha1
kind: ResourceSignature
metadata:
  annotations:
    messageScope: spec
    signature: LS0tLS1CRUdJTiBQR1AgU0lHTkFUVVJFLS0t
  name: rsig-test-cm
spec:
  data:
    - message: YXBpVmVyc2lvbjogdjEKa2luZDogQ29u
      signature: LS0tLS1CRUdJTiBQR1AgU0lHTkFUVVJFLS0t
      type: resource
```


Create this resource.
```
$ oc create -f /tmp/test-cm-rs.yaml -n secure-ns
resourcesignature.apis.integrityverifier.io/rsig-test-cm created
```


Then, run the same command again to create ConfigMap. It should be successful this time because a corresponding ResourceSignature is available in the cluster.

```
$ oc create -f /tmp/test-cm.yaml -n secure-ns
configmap/test-cm created
```


Integrity Verifier generates logs while processing admission requests in a cluster. Two types of logs are available. You can see Integrity Verifier server processing logs by a script called [`log_server.sh `](../script/log_server.sh). This includes when requests come and go, as well as errors which occured during processing. 

If you want to see the result of admission check, you can see the detail by using a script called [`log_logging.sh  `](../script/log_logging.sh).
```json
{
  "abortReason": "",
  "aborted": false,
  "allowed": false,
  "apiGroup": "",
  "apiVersion": "v1",
  "breakglass": false,
  "claim.ownerApiVersion": "",
  "claim.ownerKind": "",
  "claim.ownerName": "",
  "claim.ownerNamespace": "secure-ns",
  "creator": "",
  "detectOnly": false,
  "ivresource": false,
  "ignoreSA": false,
  "kind": "ConfigMap",
  "ma.checked": "false",
  "ma.diff": "",
  "ma.errOccured": false,
  "ma.filtered": "",
  "ma.mutated": "false",
  "maIntegrity.serviceAccount": "",
  "maIntegrity.signature": "",
  "msg": "Failed to verify signature; Signature is invalid",
  "name": "test-cm",
  "namespace": "secure-ns",
  "objLabels": "",
  "objMetaName": "test-cm",
  "operation": "CREATE",
  "org.ownerApiVersion": "",
  "org.ownerKind": "",
  "org.ownerName": "",
  "org.ownerNamespace": "secure-ns",
  "own.errOccured": false,
  "own.owners": "null",
  "own.verified": false,
  "protected": true,
  "reasonCode": "invalid-signature",
  "request.dump": "",
  "request.objectHash": "",
  "request.objectHashType": "",
  "request.uid": "bdb62f22-22f8-4a4d-9ead-cc034e4ce07b",
  "requestScope": "Namespaced",
  "sessionTrace": "time=2020-09-23T02:45:19Z level=trace msg=New Admission Request Sent aborted=false allowed=true apiVersion=apis.integrityverifier.io/v1alpha1 kind=ResourceSigningProfile name=sample-rsp namespace=secure-ns operation=UPDATE\n",
  "sig.allow": false,
  "sig.errMsg": "",
  "sig.errOccured": true,
  "sig.errReason": "Failed to verify signature; Signature is invalid",
  "timestamp": "2020-09-23T02:45:19.728Z",
  "type": "",
  "userInfo": "{\"username\":\"IAM#gajan@jp.ibm.com\",\"groups\":[\"admin\",\"iv-group\",\"system:authenticated\"]}",
  "userName": "IAM#gajan@jp.ibm.com",
  "verified": false
}
```

### Clean up Integrity Verifier from the cluster

When you want to remove Integrity Verifier from a cluster, run the uninstaller script [`delete_verifier.sh`](../scripts/delete_verifier.sh).
```
$ cd integrity-verifier
$ make delete-tmp-cr
$ make delete-operator
```

