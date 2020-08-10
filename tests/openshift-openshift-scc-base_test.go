package tests_test

import (
	"sigs.k8s.io/kustomize/v3/k8sdeps/kunstruct"
	"sigs.k8s.io/kustomize/v3/k8sdeps/transformer"
	"sigs.k8s.io/kustomize/v3/pkg/fs"
	"sigs.k8s.io/kustomize/v3/pkg/loader"
	"sigs.k8s.io/kustomize/v3/pkg/plugins"
	"sigs.k8s.io/kustomize/v3/pkg/resmap"
	"sigs.k8s.io/kustomize/v3/pkg/resource"
	"sigs.k8s.io/kustomize/v3/pkg/target"
	"sigs.k8s.io/kustomize/v3/pkg/validators"
	"testing"
)

func writeOpenshiftSccBase(th *KustTestHarness) {
	th.writeF("/manifests/openshift/openshift-scc/base/kubeflow-anyuid-scc-istio.yaml", `
apiVersion: security.openshift.io/v1
kind: SecurityContextConstraints
metadata:
  annotations:
    kubernetes.io/description: kubeflow-anyuid provides all features of the restricted SCC
      but allows users to run with any UID and any GID.
  name: kubeflow-anyuid-istio
allowHostDirVolumePlugin: false
allowHostIPC: false
allowHostNetwork: false
allowHostPID: false
allowHostPorts: false
allowPrivilegeEscalation: true
allowPrivilegedContainer: false
allowedCapabilities: null
defaultAddCapabilities: null
fsGroup:
  type: RunAsAny
groups:
- system:cluster-admins
priority: 10
readOnlyRootFilesystem: false
requiredDropCapabilities:
- MKNOD
runAsUser:
  type: RunAsAny
seLinuxContext:
  type: MustRunAs
supplementalGroups:
  type: RunAsAny
users:
#Component: istio/istio-install
- system:serviceaccount:istio-system:istio-egressgateway-service-account
- system:serviceaccount:istio-system:istio-citadel-service-account
- system:serviceaccount:istio-system:istio-ingressgateway-service-account
- system:serviceaccount:istio-system:istio-cleanup-old-ca-service-account
- system:serviceaccount:istio-system:istio-mixer-post-install-account
- system:serviceaccount:istio-system:istio-mixer-service-account
- system:serviceaccount:istio-system:istio-pilot-service-account
- system:serviceaccount:istio-system:istio-sidecar-injector-service-account
- system:serviceaccount:istio-system:istio-sidecar-injector-service-account
- system:serviceaccount:istio-system:istio-galley-service-account
- system:serviceaccount:istio-system:prometheus
#Component: istio/cluster-local-gateway
- system:serviceaccount:cluster-local-gateway-service-account
volumes:
- configMap
- downwardAPI
- emptyDir
- persistentVolumeClaim
- projected
- secret
`)
	th.writeF("/manifests/openshift/openshift-scc/base/kubeflow-anyuid-scc.yaml", `
apiVersion: security.openshift.io/v1
kind: SecurityContextConstraints
metadata:
  annotations:
    kubernetes.io/description: kubeflow-anyuid provides all features of the restricted SCC
      but allows users to run with any UID and any GID.
  name: kubeflow-anyuid-$(NAMESPACE)
allowHostDirVolumePlugin: false
allowHostIPC: false
allowHostNetwork: false
allowHostPID: false
allowHostPorts: false
allowPrivilegeEscalation: true
allowPrivilegedContainer: false
allowedCapabilities: null
defaultAddCapabilities: null
fsGroup:
  type: RunAsAny
groups:
- system:cluster-admins
priority: 10
readOnlyRootFilesystem: false
requiredDropCapabilities:
- MKNOD
runAsUser:
  type: RunAsAny
seLinuxContext:
  type: MustRunAs
supplementalGroups:
  type: RunAsAny
users:
#Metadata DB accesses files owned by root
- system:serviceaccount:$(NAMESPACE):metadatadb
#Minio accesses files owned by root
- system:serviceaccount:$(NAMESPACE):minio
#Katib injects container into pods which does not run as non-root user, trying to find Dockerfile for that image and fix it
#- system:serviceaccount:kubeflow:default
volumes:
- configMap
- downwardAPI
- emptyDir
- persistentVolumeClaim
- projected
- secret
`)
	th.writeF("/manifests/openshift/openshift-scc/base/params.yaml", `
varReference:
- path: users
  kind: SecurityContextConstraints
- path: metadata/name
  kind: SecurityContextConstraints
`)
	th.writeF("/manifests/openshift/openshift-scc/base/params.env", `
namespace=kubeflow
`)
	th.writeK("/manifests/openshift/openshift-scc/base", `
apiVersion: kustomize.config.k8s.io/v1beta1
kind: Kustomization
namespace: kubeflow
resources:
- kubeflow-anyuid-scc-istio.yaml
- kubeflow-anyuid-scc.yaml
vars:
  - name: NAMESPACE
    objref:
      apiVersion: v1
      kind: ConfigMap
      name: scc-namespace-check
    fieldref:
      fieldpath: metadata.namespace
configMapGenerator:
  - name: scc-namespace-check
configurations:
- params.yaml
`)
}

func TestOpenshiftSccBase(t *testing.T) {
	th := NewKustTestHarness(t, "/manifests/openshift/openshift-scc/base")
	writeOpenshiftSccBase(th)
	m, err := th.makeKustTarget().MakeCustomizedResMap()
	if err != nil {
		t.Fatalf("Err: %v", err)
	}
	expected, err := m.AsYaml()
	if err != nil {
		t.Fatalf("Err: %v", err)
	}
	targetPath := "../openshift/openshift-scc/base"
	fsys := fs.MakeRealFS()
	lrc := loader.RestrictionRootOnly
	_loader, loaderErr := loader.NewLoader(lrc, validators.MakeFakeValidator(), targetPath, fsys)
	if loaderErr != nil {
		t.Fatalf("could not load kustomize loader: %v", loaderErr)
	}
	rf := resmap.NewFactory(resource.NewFactory(kunstruct.NewKunstructuredFactoryImpl()), transformer.NewFactoryImpl())
	pc := plugins.DefaultPluginConfig()
	kt, err := target.NewKustTarget(_loader, rf, transformer.NewFactoryImpl(), plugins.NewLoader(pc, rf))
	if err != nil {
		th.t.Fatalf("Unexpected construction error %v", err)
	}
	actual, err := kt.MakeCustomizedResMap()
	if err != nil {
		t.Fatalf("Err: %v", err)
	}
	th.assertActualEqualsExpected(actual, string(expected))
}
