package prom

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"gopkg.in/yaml.v3"

	"github.com/grafana/grafana/pkg/services/datasources"
	"github.com/grafana/grafana/pkg/services/ngalert/models"
	"github.com/grafana/grafana/pkg/util"
)

const (
	// ruleUIDLabel is a special label that can be used to set a custom UID for a Prometheus
	// alert rule when converting it to a Grafana alert rule. If this label is not present,
	// a stable UID will be generated automatically based on the rule's data.
	ruleUIDLabel = "__grafana_alert_rule_uid__"
)

const (
	queryRefID          = "query"
	prometheusMathRefID = "prometheus_math"
	thresholdRefID      = "threshold"
)

// Config defines the configuration options for the Prometheus to Grafana rules converter.
type Config struct {
	DatasourceUID    string
	DatasourceType   string
	FromTimeRange    *time.Duration
	EvaluationOffset *time.Duration
	ExecErrState     models.ExecutionErrorState
	NoDataState      models.NoDataState
	RecordingRules   RulesConfig
	AlertRules       RulesConfig
}

// RulesConfig contains configuration that applies to either recording or alerting rules.
type RulesConfig struct {
	IsPaused bool
}

var (
	defaultTimeRange        = 600 * time.Second
	defaultEvaluationOffset = 0 * time.Minute

	defaultConfig = Config{
		FromTimeRange:    &defaultTimeRange,
		EvaluationOffset: &defaultEvaluationOffset,
		ExecErrState:     models.ErrorErrState,
		NoDataState:      models.OK,
	}
)

type Converter struct {
	cfg Config
}

// NewConverter creates a new Converter instance with the provided configuration.
// It validates the configuration and returns an error if any required fields are missing
// or if the configuration is invalid.
func NewConverter(cfg Config) (*Converter, error) {
	if cfg.DatasourceUID == "" {
		return nil, fmt.Errorf("datasource UID is required")
	}
	if cfg.DatasourceType == "" {
		return nil, fmt.Errorf("datasource type is required")
	}
	if cfg.FromTimeRange == nil {
		cfg.FromTimeRange = defaultConfig.FromTimeRange
	}
	if cfg.EvaluationOffset == nil {
		cfg.EvaluationOffset = defaultConfig.EvaluationOffset
	}
	if cfg.ExecErrState == "" {
		cfg.ExecErrState = defaultConfig.ExecErrState
	}
	if cfg.NoDataState == "" {
		cfg.NoDataState = defaultConfig.NoDataState
	}

	if cfg.DatasourceType != datasources.DS_PROMETHEUS && cfg.DatasourceType != datasources.DS_LOKI {
		return nil, fmt.Errorf("invalid datasource type: %s", cfg.DatasourceType)
	}

	return &Converter{
		cfg: cfg,
	}, nil
}

// PrometheusRulesToGrafana converts a Prometheus rule group into Grafana Alerting rule group.
func (p *Converter) PrometheusRulesToGrafana(orgID int64, namespaceUID string, group PrometheusRuleGroup) (*models.AlertRuleGroup, error) {
	for _, rule := range group.Rules {
		err := validatePrometheusRule(rule)
		if err != nil {
			return nil, fmt.Errorf("invalid Prometheus rule '%s': %w", rule.Alert, err)
		}
	}

	grafanaGroup, err := p.convertRuleGroup(orgID, namespaceUID, group)
	if err != nil {
		return nil, fmt.Errorf("failed to convert rule group '%s': %w", group.Name, err)
	}

	return grafanaGroup, nil
}

func validatePrometheusRule(rule PrometheusRule) error {
	if rule.KeepFiringFor != nil {
		return fmt.Errorf("keep_firing_for is not supported")
	}

	return nil
}

func (p *Converter) convertRuleGroup(orgID int64, namespaceUID string, promGroup PrometheusRuleGroup) (*models.AlertRuleGroup, error) {
	uniqueNames := map[string]int{}
	rules := make([]models.AlertRule, 0, len(promGroup.Rules))
	interval := time.Duration(promGroup.Interval)
	for i, rule := range promGroup.Rules {
		gr, err := p.convertRule(orgID, namespaceUID, promGroup.Name, rule)
		if err != nil {
			return nil, fmt.Errorf("failed to convert Prometheus rule '%s' to Grafana rule: %w", rule.Alert, err)
		}
		gr.RuleGroupIndex = i + 1
		gr.IntervalSeconds = int64(interval.Seconds())

		// Check rule title uniqueness within the group.
		uniqueNames[gr.Title]++
		if val := uniqueNames[gr.Title]; val > 1 {
			gr.Title = fmt.Sprintf("%s (%d)", gr.Title, val)
		}

		uid, err := getUID(orgID, namespaceUID, promGroup.Name, i, rule)
		if err != nil {
			return nil, fmt.Errorf("failed to generate UID for rule '%s': %w", gr.Title, err)
		}
		gr.UID = uid

		rules = append(rules, gr)
	}

	result := &models.AlertRuleGroup{
		FolderUID: namespaceUID,
		Interval:  int64(interval.Seconds()),
		Rules:     rules,
		Title:     promGroup.Name,
	}

	return result, nil
}

// getUID returns a UID for a Prometheus rule.
// If the rule has a special label its value is used.
// Otherwise, a stable UUID is generated by using a hash of the rule's data.
func getUID(orgID int64, namespaceUID string, group string, position int, promRule PrometheusRule) (string, error) {
	if uid, ok := promRule.Labels[ruleUIDLabel]; ok {
		if err := util.ValidateUID(uid); err != nil {
			return "", fmt.Errorf("invalid UID label value: %s; %w", uid, err)
		}
		return uid, nil
	}

	// Generate stable UUID based on the orgID, namespace, group and position.
	uidData := fmt.Sprintf("%d|%s|%s|%d", orgID, namespaceUID, group, position)
	u := uuid.NewSHA1(uuid.NameSpaceOID, []byte(uidData))

	return u.String(), nil
}

func (p *Converter) convertRule(orgID int64, namespaceUID, group string, rule PrometheusRule) (models.AlertRule, error) {
	var forInterval time.Duration
	if rule.For != nil {
		forInterval = time.Duration(*rule.For)
	}

	var query []models.AlertQuery
	var title string
	var isPaused bool
	var record *models.Record
	var err error

	isRecordingRule := rule.Record != ""
	query, err = p.createQuery(rule.Expr, isRecordingRule)
	if err != nil {
		return models.AlertRule{}, err
	}

	if isRecordingRule {
		record = &models.Record{
			From:   queryRefID,
			Metric: rule.Record,
		}

		isPaused = p.cfg.RecordingRules.IsPaused
		title = rule.Record
	} else {
		isPaused = p.cfg.AlertRules.IsPaused
		title = rule.Alert
	}

	labels := make(map[string]string, len(rule.Labels)+1)
	for k, v := range rule.Labels {
		labels[k] = v
	}

	originalRuleDefinition, err := yaml.Marshal(rule)
	if err != nil {
		return models.AlertRule{}, fmt.Errorf("failed to marshal original rule definition: %w", err)
	}

	result := models.AlertRule{
		OrgID:        orgID,
		NamespaceUID: namespaceUID,
		Title:        title,
		Data:         query,
		Condition:    query[len(query)-1].RefID,
		NoDataState:  p.cfg.NoDataState,
		ExecErrState: p.cfg.ExecErrState,
		Annotations:  rule.Annotations,
		Labels:       labels,
		For:          forInterval,
		RuleGroup:    group,
		IsPaused:     isPaused,
		Record:       record,
		Metadata: models.AlertRuleMetadata{
			PrometheusStyleRule: &models.PrometheusStyleRule{
				OriginalRuleDefinition: string(originalRuleDefinition),
			},
		},
	}

	return result, nil
}

// createQuery constructs the alert query nodes for a given Prometheus rule expression.
// It returns a slice of AlertQuery that represent the evaluation steps for the rule.
//
// For recording rules it generates a single query node that
// executes the PromQL query in the configured datasource.
//
// For alerting rules, it generates three query nodes:
//  1. Query Node (query): Executes the PromQL query using the configured datasource.
//  2. Math Node (prometheus_math): Applies a math expression "is_number($query) || is_nan($query) || is_inf($query)".
//  3. Threshold Node (threshold): Gets the result from the math node and checks that it's greater than 0.
//
// This is needed to ensure that we keep the Prometheus behaviour, where any returned result
// is considered alerting, and only when the query returns no data is the alert treated as normal.
func (p *Converter) createQuery(expr string, isRecordingRule bool) ([]models.AlertQuery, error) {
	queryNode, err := createQueryNode(p.cfg.DatasourceUID, p.cfg.DatasourceType, expr, *p.cfg.FromTimeRange, *p.cfg.EvaluationOffset)
	if err != nil {
		return nil, err
	}

	if isRecordingRule {
		return []models.AlertQuery{queryNode}, nil
	}

	mathNode, err := createMathNode()
	if err != nil {
		return nil, err
	}

	thresholdNode, err := createThresholdNode()
	if err != nil {
		return nil, err
	}

	return []models.AlertQuery{queryNode, mathNode, thresholdNode}, nil
}
