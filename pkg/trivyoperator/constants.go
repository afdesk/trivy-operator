package trivyoperator

const (
	// NamespaceName the name of the namespace in which Trivy-operator stores its
	// configuration and where it runs scan jobs.
	NamespaceName = "trivy-operator"

	// ConfigMapName the name of the ConfigMap where Trivy-operator stores its
	// configuration.
	ConfigMapName = "trivy-operator"

	// SecretName the name of the secret where Trivy-operator stores is sensitive
	// configuration.
	SecretName = "trivy-operator"

	// PoliciesConfigMapName the name of the ConfigMap used to store OPA Rego
	// policies.
	PoliciesConfigMapName = "trivy-operator-policies-config"

	// PoliciesConfigMapName the name of the ConfigMap used to store Trivy configuration.
	TrivyConfigMapName = "trivy-operator-trivy-config"
)

const (
	LabelResourceKind      = "trivy-operator.resource.kind"
	LabelResourceName      = "trivy-operator.resource.name"
	LabelResourceNameHash  = "trivy-operator.resource.name-hash"
	LabelResourceNamespace = "trivy-operator.resource.namespace"
	LabelContainerName     = "trivy-operator.container.name"
	LabelResourceSpecHash  = "resource-spec-hash"
	LabelPluginConfigHash  = "plugin-config-hash"
	LabelResourceImageID   = "resource-image-id"
	LabelReusedReport      = "reused-report"
	LabelCoreComponent     = "component"
	LabelAddon             = "k8s-app"

	LabelVulnerabilityReportScanner = "vulnerabilityReport.scanner"
	LabelNodeInfoCollector          = "node-info.collector"

	LabelK8SAppManagedBy = "app.kubernetes.io/managed-by"
	AppTrivyOperator     = "trivy-operator"

	// openshift core component
	LabelOpenShiftAPIServer         = "apiserver"
	LabelOpenShiftControllerManager = "kube-controller-manager"
	LabelOpenShiftScheduler         = "scheduler"
	LabelOpenShiftEtcd              = "etcd"
	LabelKbom                       = "trivy-operator.aquasecurity.github.io/sbom-type"
)

const (
	AnnotationContainerImages = "trivy-operator.container-images"
)
