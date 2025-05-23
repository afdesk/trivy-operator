package controller

import (
	"context"
	"fmt"
	"strings"

	"github.com/go-logr/logr"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/aquasecurity/trivy-operator/pkg/apis/aquasecurity/v1alpha1"
	"github.com/aquasecurity/trivy-operator/pkg/configauditreport"
	"github.com/aquasecurity/trivy-operator/pkg/ext"
	"github.com/aquasecurity/trivy-operator/pkg/infraassessment"
	"github.com/aquasecurity/trivy-operator/pkg/kube"
	"github.com/aquasecurity/trivy-operator/pkg/operator/etc"
	"github.com/aquasecurity/trivy-operator/pkg/policy"
	"github.com/aquasecurity/trivy-operator/pkg/trivyoperator"
	"github.com/aquasecurity/trivy/pkg/iac/scan"
)

func Policies(ctx context.Context, config etc.Config, c client.Client, cac configauditreport.ConfigAuditConfig, log logr.Logger, pl policy.Loader, clusterVersion ...string) (*policy.Policies, error) {
	cm := &corev1.ConfigMap{}

	err := c.Get(ctx, client.ObjectKey{
		Namespace: config.Namespace,
		Name:      trivyoperator.PoliciesConfigMapName,
	}, cm)
	if err != nil {
		if !apierrors.IsNotFound(err) {
			return nil, fmt.Errorf("failed getting policies from configmap: %s/%s: %w", config.Namespace, trivyoperator.PoliciesConfigMapName, err)
		}
	}
	var version string
	if len(clusterVersion) > 0 {
		version = clusterVersion[0]
	}
	return policy.NewPolicies(cm.Data, cac, log, pl, version, config.CacheReportTTL), nil
}

func ConfigurePolicies(ctx context.Context, config etc.Config, c client.Client, cac configauditreport.ConfigAuditConfig, log logr.Logger, pl policy.Loader, clusterVersion ...string) (*policy.Policies, error) {
	policies, err := Policies(ctx, config, c, cac, log, pl, clusterVersion...)
	if err != nil {
		return nil, err
	}
	if err := policies.Load(); err != nil {
		return nil, fmt.Errorf("load policies: %w", err)
	}

	if err := policies.InitScanner(); err != nil {
		return nil, fmt.Errorf("init scanner: %w", err)
	}
	return policies, nil
}

func filter(results scan.Results, resource client.Object, bi trivyoperator.BuildInfo, cd trivyoperator.ConfigData, c etc.Config, defaultSeverity string) Misconfiguration {
	misconfiguration := Misconfiguration{}
	infraChecks := make([]v1alpha1.Check, 0)
	checks := make([]v1alpha1.Check, 0)

	for _, result := range results {
		if !policy.HasSeverity(result.Severity(), defaultSeverity) {
			continue
		}

		id := policy.GetResultID(result)

		// record only misconfig failed checks
		if cd.ReportRecordFailedChecksOnly() && result.Status() == scan.StatusPassed {
			continue
		}
		currentCheck := getCheck(result, id)
		if currentCheck == nil {
			continue
		}
		if infraCheck(id) {
			if strings.HasPrefix(id, "N/A") {
				continue
			}
			if k8sCoreComponent(resource) {
				infraChecks = append(infraChecks, *currentCheck)
			}
			continue
		}
		checks = append(checks, *currentCheck)
	}
	kind := resource.GetObjectKind().GroupVersionKind().Kind
	if kube.IsRoleTypes(kube.Kind(kind)) && !c.MergeRbacFindingWithConfigAudit {
		misconfiguration.rbacAssessmentReportData = v1alpha1.RbacAssessmentReportData{
			Scanner: scanner(bi),
			Summary: v1alpha1.RbacAssessmentSummaryFromChecks(checks),
			Checks:  checks,
		}
		return misconfiguration
	}
	misconfiguration.configAuditReportData = v1alpha1.ConfigAuditReportData{
		UpdateTimestamp: metav1.NewTime(ext.NewSystemClock().Now()),
		Scanner:         scanner(bi),
		Summary:         v1alpha1.ConfigAuditSummaryFromChecks(checks),
		Checks:          checks,
	}
	if c.InfraAssessmentScannerEnabled {
		misconfiguration.infraAssessmentReportData = v1alpha1.InfraAssessmentReportData{
			Scanner: scanner(bi),
			Summary: v1alpha1.InfraAssessmentSummaryFromChecks(infraChecks),
			Checks:  infraChecks,
		}
	}

	return misconfiguration
}

func evaluate(ctx context.Context, policies *policy.Policies, resource client.Object, bi trivyoperator.BuildInfo, cd trivyoperator.ConfigData, c etc.Config, inputs ...[]byte) (Misconfiguration, error) {
	results, err := policies.Eval(ctx, resource, inputs...)
	if err != nil {
		return Misconfiguration{}, err
	}

	return filter(results, resource, bi, cd, c, policies.GetDefaultSeverity()), nil
}

func scanner(bi trivyoperator.BuildInfo) v1alpha1.Scanner {
	return v1alpha1.Scanner{
		Name:    v1alpha1.ScannerNameTrivy,
		Vendor:  "Aqua Security",
		Version: bi.Version,
	}
}

func hasInfraReport(ctx context.Context, node *corev1.Node, infraReadWriter infraassessment.ReadWriter) (bool, error) {
	report, err := infraReadWriter.FindClusterReportByOwner(ctx, kube.ObjectRef{Kind: kube.KindNode, Name: node.Name})
	if err != nil {
		return false, err
	}
	return report != nil, nil
}
