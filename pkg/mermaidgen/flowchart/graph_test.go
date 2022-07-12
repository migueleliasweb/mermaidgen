package flowchart

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/prometheus/client_golang/api"
	v1 "github.com/prometheus/client_golang/api/prometheus/v1"
	"github.com/prometheus/prometheus/promql/parser"
)

func TestFlowchart(t *testing.T) {
	fc := New("", nil)
	fc.Lines = []Line{
		{
			InlineItems: []InlineItem{
				&Node{
					ID:   "A",
					Text: "Christmas",
				},
				&LinkArrow,
				&Node{
					ID:   "B",
					Text: "Go shopping",
				},
			},
		},
	}

	fmt.Println(fc)
}

func extractQueryMatchers(query string) []string {
	expr, _ := parser.ParseExpr(query)

	labelMatchers := parser.ExtractSelectors(expr)

	queryMatchersMap := map[string]bool{}
	result := []string{}

	for _, lms := range labelMatchers {
		for _, lm := range lms {
			if lm.Name == "__name__" {
				queryMatchersMap[lm.Value] = true
			}
		}
	}

	for k := range queryMatchersMap {
		result = append(result, k)
	}

	return result
}

func generateData() map[string][]string {
	client, err := api.NewClient(api.Config{
		// https://prometheus.demo.do.prometheus.io/rules
		Address: "https://prometheus.monitoring.prd.tyro.cloud",
	})

	if err != nil {
		fmt.Printf("Error creating client: %v\n", err)
		os.Exit(1)
	}

	v1api := v1.NewAPI(client)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	rules, _ := v1api.Rules(ctx)

	// need to add group name as metrics can repeat in different groups
	ruleNameWithRelatedQueryNames := map[string][]string{}

	for _, group := range rules.Groups {
		// groupName := group.Name
		for _, r := range group.Rules {
			switch v := r.(type) {
			case v1.RecordingRule:
				rule := r.(v1.RecordingRule)
				ruleNameWithRelatedQueryNames[rule.Name] = extractQueryMatchers(rule.Query)
			case v1.AlertingRule:
				rule := r.(v1.AlertingRule)
				ruleNameWithRelatedQueryNames[rule.Name] = extractQueryMatchers(rule.Query)
			default:
				fmt.Printf("unknown rule type %s", v)
			}
		}
	}

	return ruleNameWithRelatedQueryNames
}

// alert:CriticalLokiDown
// expr:absent(up{job="loki-headless"} == 1)

// alert:AlertmanagerClusterDown
// expr:(count by(namespace, service) (avg_over_time(up{job="alertmanager-main",namespace="monitoring"}[5m]) < 0.5) / count by(namespace, service) (up{job="alertmanager-main",namespace="monitoring"})) >= 0.5

// 1[alert:CriticalLokiDown]-->up
// 2[alert:AlertmanagerClusterDown]-->up

func outputMetricRecursive(data map[string][]string, metricName string) string {
	relatedMetrics, found := data[metricName]

	if !found {
		return ""
	}

	var b bytes.Buffer

	for _, relatedMetric := range relatedMetrics {
		b.WriteString(fmt.Sprintf("    %s-->%s\n",
			relatedMetric,
			metricName,
		))

		b.WriteString(outputMetricRecursive(data, relatedMetric))
	}

	return b.String()
}

func TestFakeMain(t *testing.T) {
	data := generateData()

	var b bytes.Buffer

	b.WriteString("flowchart TB\n")

	// topMetricName := "TAPEgressSLOTrendingHighError" // somewhat easy to understand
	// topMetricName := "TAPIngressSLIDownForExtendedPeriod" // somewhat easy to understand
	// topMetricName := "tap_ingress::job:probe_success"
	// topMetricName := "TAPEgressSLOTrendingHighError" // huge + complex
	// topMetricName := "TapNodeToApiConnectivityOkSloBudgetConsumptionHigh"
	topMetricName := "TapClusterSloBudgetConsumptionHigh"
	// topMetricName := "platform:up"
	// topMetricName := "kubernetes:job:apiserver_request_errors:ratio_rate5m"

	b.WriteString(outputMetricRecursive(data, topMetricName))

	fmt.Println(b.String())
}
