# Settings

Trivy Operator reads configuration settings from ConfigMaps, as well as Secrets that hold
confidential settings (such as a GitHub token). Trivy-Operator plugins read configuration and secret data from ConfigMaps
and Secrets named after the plugin. For example, Trivy configuration is stored in the ConfigMap and Secret named
`trivy-operator-trivy-config`.

You can change the default settings with `kubectl patch` or `kubectl edit` commands. For example, by default Trivy
displays vulnerabilities with all severity levels (`UNKNOWN`, `LOW`, `MEDIUM`, `HIGH`, `CRITICAL`). However, you can
display only `HIGH` and `CRITICAL` vulnerabilities by patching the `trivy.severity` value in the `trivy-operator-trivy-config`
ConfigMap:

```
TRIVY_OPERATOR_NAMESPACE=<your trivy operator namespace>
```

```
kubectl patch cm trivy-operator-trivy-config -n $TRIVY_OPERATOR_NAMESPACE \
  --type merge \
  -p "$(cat <<EOF
{
  "data": {
    "trivy.severity": "HIGH,CRITICAL"
  }
}
EOF
)"
```

To set the GitHub token used by Trivy add the `trivy.githubToken` value to the `trivy-operator-trivy-config` Secret:

```
TRIVY_OPERATOR_NAMESPACE=<your trivy opersator namespace>
GITHUB_TOKEN=<your token>
```

```
kubectl patch secret trivy-operator-trivy-config -n $TRIVY_OPERATOR_NAMESPACE \
  --type merge \
  -p "$(cat <<EOF
{
  "data": {
    "trivy.githubToken": "$(echo -n $GITHUB_TOKEN | base64)"
  }
}
EOF
)"
```

The following table lists available settings with their default values. Check plugins' documentation to see
configuration settings for common use cases. For example, switch Trivy from [Standalone] to [ClientServer] mode.

| CONFIG   KEY                                                        | DEFAULT                               | DESCRIPTION                                                                                                                                                                                                                         |
|------------------------------------------------|---------------------------------------|-------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| `vulnerabilityReports.scanner`                                       | `Trivy`                               | The name of the plugin that generates vulnerability reports. Either `Trivy` or `Aqua`.                                                                                                                                              |
| `vulnerabilityReports.scanJobsInSameNamespace`                       | `"false"`                             | Whether to run vulnerability scan jobs in same namespace of workload. Set `"true"` to enable.                                                                                                                                       |
| `configAuditReports.scanner`                                         | `Trivy`                               | The name of the plugin that generates config audit reports.                                                                                                                                                                         |
| `scanJob.affinity`                                                   | N/A                                   | JSON representation of the [affinity] to be applied to the scanner pods and node-collector. Example: `'{"nodeAffinity":{"requiredDuringSchedulingIgnoredDuringExecution":{"nodeSelectorTerms":[{"matchExpressions":[{"key":"kubernetes.io/os","operator":"In","values":["linux"]}]},{"matchExpressions":[{"key":"virtual-kubelet.io/provider","operator":"DoesNotExist"}]}]}}}'`           |
| `scanJob.tolerations`                                                | N/A                                   | JSON representation of the [tolerations] to be applied to the scanner pods and node-collector so that they can run on nodes with matching taints. Example: `'[{"key":"key1", "operator":"Equal", "value":"value1", "effect":"NoSchedule"}]'`           |
| `nodeCollector.volumeMounts`| see helm/values.yaml | node-collector pod volumeMounts definition for collecting config files information
| `nodeCollector.volumes`| see helm/values.yaml | node-collector pod volumes definition for collecting config files information
| `scanJob.nodeSelector`                                                | N/A                                   | JSON representation of the [nodeSelector] to be applied to the scanner pods so that they can run on nodes with matching labels. Example: `'{"example.com/node-type":"worker", "cpu-type": "sandylake"}'`           |
| `scanJob.annotations`                                                 | N/A                                   | One-line comma-separated representation of the annotations which the user wants the scanner pods to be annotated with. Example: `foo=bar,env=stage` will annotate the scanner pods with the annotations `foo: bar` and `env: stage` |
| `scanJob.templateLabel`                                               | N/A                                   | One-line comma-separated representation of the template labels which the user wants the scanner pods to be labeled with. Example: `foo=bar,env=stage` will labeled the scanner pods with the labels `foo: bar` and `env: stage`     |
| `scanJob.podTemplatePodSecurityContext`                               | N/A                                   | One-line JSON representation of the template securityContext which the user wants the scanner and node collector pods to be secured with. Example: `{"RunAsUser": 1000, "RunAsGroup": 1000, "RunAsNonRoot": true}`                |
| `scanJob.podTemplateContainerSecurityContext`                         | N/A| One-line JSON representation of the template securityContext which the user wants the scanner and node collector containers (and their initContainers) to be amended with. Example: `{"allowPrivilegeEscalation": false, "capabilities": { "drop": ["ALL"]},"privileged": false, "readOnlyRootFilesystem": true }`|
| `scanJob.podPriorityClassName`                                       | `""`                                  | The value of the priorityClassName for job                                                                                                                                                                   |
| `compliance.failEntriesLimit`                                         | `"10"`                                | Limit the number of fail entries per control check in the cluster compliance detail report.                                                                                                                                         |
| `compliance.reportType`               | `summary`                   | this flag control the type of report generated summary or all                |
| `compliance.cron`                     | `0 */6 * * *`                   | this flag control the cron interval for compliance report generation                |
| `scanJob.compressLogs`                                         | `"true"`                              | Control whether scanjob output should be compressed                                                                                                                                     |
| `nodeCollector.excludeNodes`                        | `""`                      | excludeNodes comma-separated node labels that the node-collector job should exclude from scanning (example kubernetes.io/arch=arm64,team=dev)                                                                                                                                                                                                                                      |
| `alternateReportStorage.enabled`| `"false"` | Control where reports are written. By default this is false, so reports will be written normally as CRDs in ETCD memory. However, if you would rather reports be written to a persistent volume, flip this to true. If done a persistent volume claim will be inluded in your installation and all reports will be written there.|
| `alternateReportStorage.mountPath`|`"/mnt/data/trivy-operator"`| The mount path for your persistent volume.|
| `alternateReportStorage.volumeName`|`"trivy-operator-pvc"`| Name of your persistant volume.|
|`alternateReportStorage.storage`|`"10Gi"`| Amount of storage for your persistent volume.|
|`alternateReportStorage.podSecurityContext.runAsUser`| `10000` | Specifies the UNIX user ID that all processes in the container should run as (for the persistent volume), ensuring they don’t execute as the root user and limiting their privileges.|
|`alternateReportStorage.podSecurityContext.fsGroup`| `10000` | Defines a UNIX group ID that Kubernetes will use to change the ownership of any mounted volumes so that files created by the container (persistent volume) are accessible to processes running under that group.|

!!! note
    For parameters that use time values, such as `ScanJobTTL`, valid time units are "ns", "us" (or "µs"), "ms", "s", "m", and "h".

!!! tip
    You can delete a configuration key.For example, the following `kubectl patch` command deletes the `trivy.httpProxy` key:
    ```
    TRIVY_OPERATOR_NAMESPACE=<your trivy operator namespace>
    ```
    ```
    kubectl patch cm trivy-operator-trivy-config -n $TRIVY_OPERATOR_NAMESPACE \
      --type json \
      -p '[{"op": "remove", "path": "/data/trivy.httpProxy"}]'
    ```

[Standalone]: ./docs/vulnerability-scanning/trivy.md#standalone
[ClientServer]: ./docs/vulnerability-scanning/trivy.md#clientserver
[tolerations]: https://kubernetes.io/docs/concepts/scheduling-eviction/taint-and-toleration
