

# Custom Resource: IntegrityShield

Integrity Shield can be deployed with operator. You can configure IntegrityShield custom resource to define the configuration of IShield.

## Type of Signature Verification

Integrity Shield supports two modes of signature verification.
- `pgp`: use [gpg key](https://www.gnupg.org/index.html) for signing. certificate is not used.
- `x509`: use signing key with X509 public key certificate.

`spec.verifyType` should be set either `pgp` (default) or `x509`.

```yaml
apiVersion: apis.integrityshield.io/v1alpha1
kind: IntegrityShield
metadata:
  name: integrity-shield-server
spec:
  shieldConfig:
    verifyType: pgp
```

<!-- ## Enable Helm plugin

You can enable Helm plugin to support verification of Helm provenance and integrity (https://helm.sh/docs/topics/provenance/). By enabling this, Helm package installation is verified with its provenance file.

```yaml
spec:
  shieldConfig:
    policy:
      plugin:
      - name: helm
        enabled: false
``` -->

## Verification Key and Sign Policy Configuration

The list of verification key names should be set as `keyRingConfigs` in this CR.
The operator will start installing Integrity Shield when all key secrets listed here are ready.

Also, you can set SignPolicy here.
This policy defines signers that are allowed to create/update resources with their signature in some namespaces.
(see [How to configure SignPolicy](README_CONFIG_SIGNER_POLICY.md) for detail.)

```yaml
spec:
  keyRingConfigs:
  - name: keyring-secret
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
      - email: "sample_signer@signer.com"
```

## Resource Signing Profile Configuration
You can define one or more ResourceSigningProfiles that are installed by this operator.
This configuration is not set by default.
(see [How to configure ResourceSigningProfile](README_FOR_RESOURCE_PROTECTION_PROFILE.md) for detail.)

```yaml
spec:
  resourceSigningProfiles:
  - name: sample-rsp
    targetNamespaceSelector:
      include:
      - "secure-ns"
    protectRules:
    - match:
      - kind: "ConfigMap"
        name: "*"
```

## Define In-scope Namespaces
You can define which namespace is not checked by Integrity Shield even if ResourceSigningProfile is there.
Wildcard "*" can be used for this config. By default, Integrity Shield checks RSPs in all namespaces except ones in `kube-*` and `openshift-*` namespaces.

```yaml
spec:
  inScopeNamespaceSelector:
    include:
    - "*"
    exclude:
    - "kube-*"
    - "openshift-*"
```

## Unprocessed Requests
Some resources are not relevant to the signature-based protection by Integrity Shield.
The resources defined here are not processed in IShield admission controller (always returns `allowed`).

```yaml
spec:
  shieldConfig:
    ignore:
    - kind: Event
    - kind: Lease
    - kind: Endpoints
    - kind: TokenReview
    - kind: SubjectAccessReview
    - kind: SelfSubjectAccessReview
```

## IShield Run mode
You can set run mode. Two modes are available. `enforce` mode is default. `detect` mode always allows any admission request, but signature verification is conducted and logged for all protected resources. `enforce` is set unless specified.

```yaml
spec:
  shieldConfig:
    mode: "detect"
```

<!-- ## Install on OpenShift

When deploying OpenShift cluster, this should be set `true` (default). Then, SecurityContextConstratint (SCC) will be deployed automatically during installation. For IKS or Minikube, this should be set to `false`.

```yaml
spec:
  globalConfig:
    openShift: true
``` -->

## IShield admin

Specify user group for IShield admin with comma separated strings like the following. This value is empty by default.

```yaml
spec:
  shieldConfig:
    iShieldAdminUserGroup: "system:masters,system:cluster-admins"
```

Also, you can define IShield admin role. This role will be created automatically during installation when `autoIShieldAdminRoleCreationDisabled` is `false` (default).

```yaml
spec
  security
    iShieldAdminSubjects:
      - apiGroup: rbac.authorization.k8s.io
        kind: Group
        name: system:masters
    autoIShieldAdminRoleCreationDisabled: false
```

<!-- 
## Webhook configuration

You can specify webhook filter configuration for processing requests in IShield. As default, all requests for namespaced resources and selected cluster-scope resources are forwarded to IShield. If you want to protect a resource by IShield, it must be covered with this filter condition.

```yaml
spec:
  webhookNamespacedResource:
    apiGroups: ["*"]
    apiVersions: ["*"]
    resources: ["*"]
  webhookClusterResource:
    apiGroups: ["*"]
    apiVersions: ["*"]
    resources:
    - podsecuritypolicies
    - clusterrolebindings
    - clusterroles
``` -->

## Logging

Console log includes stdout logging from IShield server. Context log includes admission control results. Both are enabled as default. You can specify namespaces in scope. `'*'` is wildcard. `'-'` is empty stiring, which implies cluster-scope resource.
```yaml
spec:
  shieldConfig:
    log:
      consoleLog:
        enabled: true
        inScope:
        - namespace: '*'
        - namespace: '-'
      contextLog:
        enabled: true
        inScope:
        - namespace: '*'
        - namespace: '-'
      logLevel: info
```
