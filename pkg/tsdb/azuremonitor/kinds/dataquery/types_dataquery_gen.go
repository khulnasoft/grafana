// Code generated - EDITING IS FUTILE. DO NOT EDIT.
//
// Generated by:
//     public/app/plugins/gen.go
// Using jennies:
//     PluginGoTypesJenny
//
// Run 'make gen-cue' from repository root to regenerate.

// Code generated - EDITING IS FUTILE. DO NOT EDIT.

package dataquery

import (
	json "encoding/json"
	errors "errors"
	fmt "fmt"
)

type AzureMonitorQuery struct {
	// A unique identifier for the query within the list of targets.
	// In server side expressions, the refId is used as a variable name to identify results.
	// By default, the UI will assign A->Z; however setting meaningful names may be useful.
	RefId string `json:"refId"`
	// If hide is set to true, Grafana will filter out the response(s) associated with this query before returning it to the panel.
	Hide *bool `json:"hide,omitempty"`
	// Specify the query flavor
	// TODO make this required and give it a default
	QueryType *string `json:"queryType,omitempty"`
	// Azure subscription containing the resource(s) to be queried.
	// Also used for template variable queries
	Subscription *string `json:"subscription,omitempty"`
	// Subscriptions to be queried via Azure Resource Graph.
	Subscriptions []string `json:"subscriptions,omitempty"`
	// Azure Monitor Metrics sub-query properties.
	AzureMonitor *AzureMetricQuery `json:"azureMonitor,omitempty"`
	// Azure Monitor Logs sub-query properties.
	AzureLogAnalytics *AzureLogsQuery `json:"azureLogAnalytics,omitempty"`
	// Azure Resource Graph sub-query properties.
	AzureResourceGraph *AzureResourceGraphQuery `json:"azureResourceGraph,omitempty"`
	// Application Insights Traces sub-query properties.
	AzureTraces *AzureTracesQuery `json:"azureTraces,omitempty"`
	// @deprecated Legacy template variable support.
	GrafanaTemplateVariableFn *GrafanaTemplateVariableQuery `json:"grafanaTemplateVariableFn,omitempty"`
	// Resource group used in template variable queries
	ResourceGroup *string `json:"resourceGroup,omitempty"`
	// Namespace used in template variable queries
	Namespace *string `json:"namespace,omitempty"`
	// Resource used in template variable queries
	Resource *string `json:"resource,omitempty"`
	// Region used in template variable queries
	Region *string `json:"region,omitempty"`
	// Custom namespace used in template variable queries
	CustomNamespace *string `json:"customNamespace,omitempty"`
	// For mixed data sources the selected datasource is on the query level.
	// For non mixed scenarios this is undefined.
	// TODO find a better way to do this ^ that's friendly to schema
	// TODO this shouldn't be unknown but DataSourceRef | null
	Datasource any `json:"datasource,omitempty"`
	// Used only for exemplar queries from Prometheus
	Query *string `json:"query,omitempty"`
}

// NewAzureMonitorQuery creates a new AzureMonitorQuery object.
func NewAzureMonitorQuery() *AzureMonitorQuery {
	return &AzureMonitorQuery{}
}

// Defines the supported queryTypes. GrafanaTemplateVariableFn is deprecated
type AzureQueryType string

const (
	AzureQueryTypeAzureMonitor              AzureQueryType = "Azure Monitor"
	AzureQueryTypeLogAnalytics              AzureQueryType = "Azure Log Analytics"
	AzureQueryTypeAzureResourceGraph        AzureQueryType = "Azure Resource Graph"
	AzureQueryTypeAzureTraces               AzureQueryType = "Azure Traces"
	AzureQueryTypeSubscriptionsQuery        AzureQueryType = "Azure Subscriptions"
	AzureQueryTypeResourceGroupsQuery       AzureQueryType = "Azure Resource Groups"
	AzureQueryTypeNamespacesQuery           AzureQueryType = "Azure Namespaces"
	AzureQueryTypeResourceNamesQuery        AzureQueryType = "Azure Resource Names"
	AzureQueryTypeMetricNamesQuery          AzureQueryType = "Azure Metric Names"
	AzureQueryTypeWorkspacesQuery           AzureQueryType = "Azure Workspaces"
	AzureQueryTypeLocationsQuery            AzureQueryType = "Azure Regions"
	AzureQueryTypeGrafanaTemplateVariableFn AzureQueryType = "Grafana Template Variable Function"
	AzureQueryTypeTraceExemplar             AzureQueryType = "traceql"
	AzureQueryTypeCustomNamespacesQuery     AzureQueryType = "Azure Custom Namespaces"
	AzureQueryTypeCustomMetricNamesQuery    AzureQueryType = "Azure Custom Metric Names"
)

type AzureMetricQuery struct {
	// Array of resource URIs to be queried.
	Resources []AzureMonitorResource `json:"resources,omitempty"`
	// metricNamespace is used as the resource type (or resource namespace).
	// It's usually equal to the target metric namespace. e.g. microsoft.storage/storageaccounts
	// Kept the name of the variable as metricNamespace to avoid backward incompatibility issues.
	MetricNamespace *string `json:"metricNamespace,omitempty"`
	// Used as the value for the metricNamespace property when it's different from the resource namespace.
	CustomNamespace *string `json:"customNamespace,omitempty"`
	// The metric to query data for within the specified metricNamespace. e.g. UsedCapacity
	MetricName *string `json:"metricName,omitempty"`
	// The Azure region containing the resource(s).
	Region *string `json:"region,omitempty"`
	// The granularity of data points to be queried. Defaults to auto.
	TimeGrain *string `json:"timeGrain,omitempty"`
	// The aggregation to be used within the query. Defaults to the primaryAggregationType defined by the metric.
	Aggregation *string `json:"aggregation,omitempty"`
	// Filters to reduce the set of data returned. Dimensions that can be filtered on are defined by the metric.
	DimensionFilters []AzureMetricDimension `json:"dimensionFilters,omitempty"`
	// Maximum number of records to return. Defaults to 10.
	Top *string `json:"top,omitempty"`
	// Time grains that are supported by the metric.
	AllowedTimeGrainsMs []int64 `json:"allowedTimeGrainsMs,omitempty"`
	// Aliases can be set to modify the legend labels. e.g. {{ resourceGroup }}. See docs for more detail.
	Alias *string `json:"alias,omitempty"`
	// @deprecated
	TimeGrainUnit *string `json:"timeGrainUnit,omitempty"`
	// @deprecated This property was migrated to dimensionFilters and should only be accessed in the migration
	Dimension *string `json:"dimension,omitempty"`
	// @deprecated This property was migrated to dimensionFilters and should only be accessed in the migration
	DimensionFilter *string `json:"dimensionFilter,omitempty"`
	// @deprecated Use metricNamespace instead
	MetricDefinition *string `json:"metricDefinition,omitempty"`
	// @deprecated Use resourceGroup, resourceName and metricNamespace instead
	ResourceUri *string `json:"resourceUri,omitempty"`
	// @deprecated Use resources instead
	ResourceGroup *string `json:"resourceGroup,omitempty"`
	// @deprecated Use resources instead
	ResourceName *string `json:"resourceName,omitempty"`
}

// NewAzureMetricQuery creates a new AzureMetricQuery object.
func NewAzureMetricQuery() *AzureMetricQuery {
	return &AzureMetricQuery{}
}

// Azure Monitor Logs sub-query properties
type AzureLogsQuery struct {
	// KQL query to be executed.
	Query *string `json:"query,omitempty"`
	// Specifies the format results should be returned as.
	ResultFormat *ResultFormat `json:"resultFormat,omitempty"`
	// Array of resource URIs to be queried.
	Resources []string `json:"resources,omitempty"`
	// If set to true the dashboard time range will be used as a filter for the query. Otherwise the query time ranges will be used. Defaults to false.
	DashboardTime *bool `json:"dashboardTime,omitempty"`
	// If dashboardTime is set to true this value dictates which column the time filter will be applied to. Defaults to the first tables timeSpan column, the first datetime column found, or TimeGenerated
	TimeColumn *string `json:"timeColumn,omitempty"`
	// If set to true the query will be run as a basic logs query
	BasicLogsQuery *bool `json:"basicLogsQuery,omitempty"`
	// Workspace ID. This was removed in Grafana 8, but remains for backwards compat.
	Workspace *string `json:"workspace,omitempty"`
	// Denotes if logs query editor is in builder mode
	Mode *LogsEditorMode `json:"mode,omitempty"`
	// @deprecated Use resources instead
	Resource *string `json:"resource,omitempty"`
	// @deprecated Use dashboardTime instead
	IntersectTime *bool `json:"intersectTime,omitempty"`
}

// NewAzureLogsQuery creates a new AzureLogsQuery object.
func NewAzureLogsQuery() *AzureLogsQuery {
	return &AzureLogsQuery{}
}

// Application Insights Traces sub-query properties
type AzureTracesQuery struct {
	// Specifies the format results should be returned as.
	ResultFormat *ResultFormat `json:"resultFormat,omitempty"`
	// Array of resource URIs to be queried.
	Resources []string `json:"resources,omitempty"`
	// Operation ID. Used only for Traces queries.
	OperationId *string `json:"operationId,omitempty"`
	// Types of events to filter by.
	TraceTypes []string `json:"traceTypes,omitempty"`
	// Filters for property values.
	Filters []AzureTracesFilter `json:"filters,omitempty"`
	// KQL query to be executed.
	Query *string `json:"query,omitempty"`
}

// NewAzureTracesQuery creates a new AzureTracesQuery object.
func NewAzureTracesQuery() *AzureTracesQuery {
	return &AzureTracesQuery{}
}

type AzureTracesFilter struct {
	// Property name, auto-populated based on available traces.
	Property string `json:"property"`
	// Comparison operator to use. Either equals or not equals.
	Operation string `json:"operation"`
	// Values to filter by.
	Filters []string `json:"filters"`
}

// NewAzureTracesFilter creates a new AzureTracesFilter object.
func NewAzureTracesFilter() *AzureTracesFilter {
	return &AzureTracesFilter{}
}

type ResultFormat string

const (
	ResultFormatTable      ResultFormat = "table"
	ResultFormatTimeSeries ResultFormat = "time_series"
	ResultFormatTrace      ResultFormat = "trace"
	ResultFormatLogs       ResultFormat = "logs"
)

type LogsEditorMode string

const (
	LogsEditorModeBuilder LogsEditorMode = "builder"
	LogsEditorModeRaw     LogsEditorMode = "raw"
)

type AzureResourceGraphQuery struct {
	// Azure Resource Graph KQL query to be executed.
	Query *string `json:"query,omitempty"`
	// Specifies the format results should be returned as. Defaults to table.
	ResultFormat *string `json:"resultFormat,omitempty"`
}

// NewAzureResourceGraphQuery creates a new AzureResourceGraphQuery object.
func NewAzureResourceGraphQuery() *AzureResourceGraphQuery {
	return &AzureResourceGraphQuery{}
}

type AzureMonitorResource struct {
	Subscription    *string `json:"subscription,omitempty"`
	ResourceGroup   *string `json:"resourceGroup,omitempty"`
	ResourceName    *string `json:"resourceName,omitempty"`
	MetricNamespace *string `json:"metricNamespace,omitempty"`
	Region          *string `json:"region,omitempty"`
}

// NewAzureMonitorResource creates a new AzureMonitorResource object.
func NewAzureMonitorResource() *AzureMonitorResource {
	return &AzureMonitorResource{}
}

type AzureMetricDimension struct {
	// Name of Dimension to be filtered on.
	Dimension *string `json:"dimension,omitempty"`
	// String denoting the filter operation. Supports 'eq' - equals,'ne' - not equals, 'sw' - starts with. Note that some dimensions may not support all operators.
	Operator *string `json:"operator,omitempty"`
	// Values to match with the filter.
	Filters []string `json:"filters,omitempty"`
	// @deprecated filter is deprecated in favour of filters to support multiselect.
	Filter *string `json:"filter,omitempty"`
}

// NewAzureMetricDimension creates a new AzureMetricDimension object.
func NewAzureMetricDimension() *AzureMetricDimension {
	return &AzureMetricDimension{}
}

type GrafanaTemplateVariableQueryType string

const (
	GrafanaTemplateVariableQueryTypeAppInsightsMetricNameQuery GrafanaTemplateVariableQueryType = "AppInsightsMetricNameQuery"
	GrafanaTemplateVariableQueryTypeAppInsightsGroupByQuery    GrafanaTemplateVariableQueryType = "AppInsightsGroupByQuery"
	GrafanaTemplateVariableQueryTypeSubscriptionsQuery         GrafanaTemplateVariableQueryType = "SubscriptionsQuery"
	GrafanaTemplateVariableQueryTypeResourceGroupsQuery        GrafanaTemplateVariableQueryType = "ResourceGroupsQuery"
	GrafanaTemplateVariableQueryTypeResourceNamesQuery         GrafanaTemplateVariableQueryType = "ResourceNamesQuery"
	GrafanaTemplateVariableQueryTypeMetricNamespaceQuery       GrafanaTemplateVariableQueryType = "MetricNamespaceQuery"
	GrafanaTemplateVariableQueryTypeMetricNamesQuery           GrafanaTemplateVariableQueryType = "MetricNamesQuery"
	GrafanaTemplateVariableQueryTypeWorkspacesQuery            GrafanaTemplateVariableQueryType = "WorkspacesQuery"
	GrafanaTemplateVariableQueryTypeUnknownQuery               GrafanaTemplateVariableQueryType = "UnknownQuery"
)

type BaseGrafanaTemplateVariableQuery struct {
	RawQuery *string `json:"rawQuery,omitempty"`
}

// NewBaseGrafanaTemplateVariableQuery creates a new BaseGrafanaTemplateVariableQuery object.
func NewBaseGrafanaTemplateVariableQuery() *BaseGrafanaTemplateVariableQuery {
	return &BaseGrafanaTemplateVariableQuery{}
}

type UnknownQuery struct {
	RawQuery *string `json:"rawQuery,omitempty"`
	Kind     string  `json:"kind"`
}

// NewUnknownQuery creates a new UnknownQuery object.
func NewUnknownQuery() *UnknownQuery {
	return &UnknownQuery{
		Kind: "UnknownQuery",
	}
}

type AppInsightsMetricNameQuery struct {
	RawQuery *string `json:"rawQuery,omitempty"`
	Kind     string  `json:"kind"`
}

// NewAppInsightsMetricNameQuery creates a new AppInsightsMetricNameQuery object.
func NewAppInsightsMetricNameQuery() *AppInsightsMetricNameQuery {
	return &AppInsightsMetricNameQuery{
		Kind: "AppInsightsMetricNameQuery",
	}
}

type AppInsightsGroupByQuery struct {
	RawQuery   *string `json:"rawQuery,omitempty"`
	Kind       string  `json:"kind"`
	MetricName string  `json:"metricName"`
}

// NewAppInsightsGroupByQuery creates a new AppInsightsGroupByQuery object.
func NewAppInsightsGroupByQuery() *AppInsightsGroupByQuery {
	return &AppInsightsGroupByQuery{
		Kind: "AppInsightsGroupByQuery",
	}
}

type SubscriptionsQuery struct {
	RawQuery *string `json:"rawQuery,omitempty"`
	Kind     string  `json:"kind"`
}

// NewSubscriptionsQuery creates a new SubscriptionsQuery object.
func NewSubscriptionsQuery() *SubscriptionsQuery {
	return &SubscriptionsQuery{
		Kind: "SubscriptionsQuery",
	}
}

type ResourceGroupsQuery struct {
	RawQuery     *string `json:"rawQuery,omitempty"`
	Kind         string  `json:"kind"`
	Subscription string  `json:"subscription"`
}

// NewResourceGroupsQuery creates a new ResourceGroupsQuery object.
func NewResourceGroupsQuery() *ResourceGroupsQuery {
	return &ResourceGroupsQuery{
		Kind: "ResourceGroupsQuery",
	}
}

type ResourceNamesQuery struct {
	RawQuery        *string `json:"rawQuery,omitempty"`
	Kind            string  `json:"kind"`
	Subscription    string  `json:"subscription"`
	ResourceGroup   string  `json:"resourceGroup"`
	MetricNamespace string  `json:"metricNamespace"`
}

// NewResourceNamesQuery creates a new ResourceNamesQuery object.
func NewResourceNamesQuery() *ResourceNamesQuery {
	return &ResourceNamesQuery{
		Kind: "ResourceNamesQuery",
	}
}

type MetricNamespaceQuery struct {
	RawQuery        *string `json:"rawQuery,omitempty"`
	Kind            string  `json:"kind"`
	Subscription    string  `json:"subscription"`
	ResourceGroup   string  `json:"resourceGroup"`
	MetricNamespace *string `json:"metricNamespace,omitempty"`
	ResourceName    *string `json:"resourceName,omitempty"`
}

// NewMetricNamespaceQuery creates a new MetricNamespaceQuery object.
func NewMetricNamespaceQuery() *MetricNamespaceQuery {
	return &MetricNamespaceQuery{
		Kind: "MetricNamespaceQuery",
	}
}

// @deprecated Use MetricNamespaceQuery instead
type MetricDefinitionsQuery struct {
	RawQuery        *string `json:"rawQuery,omitempty"`
	Kind            string  `json:"kind"`
	Subscription    string  `json:"subscription"`
	ResourceGroup   string  `json:"resourceGroup"`
	MetricNamespace *string `json:"metricNamespace,omitempty"`
	ResourceName    *string `json:"resourceName,omitempty"`
}

// NewMetricDefinitionsQuery creates a new MetricDefinitionsQuery object.
func NewMetricDefinitionsQuery() *MetricDefinitionsQuery {
	return &MetricDefinitionsQuery{
		Kind: "MetricDefinitionsQuery",
	}
}

type MetricNamesQuery struct {
	RawQuery        *string `json:"rawQuery,omitempty"`
	Kind            string  `json:"kind"`
	Subscription    string  `json:"subscription"`
	ResourceGroup   string  `json:"resourceGroup"`
	ResourceName    string  `json:"resourceName"`
	MetricNamespace string  `json:"metricNamespace"`
}

// NewMetricNamesQuery creates a new MetricNamesQuery object.
func NewMetricNamesQuery() *MetricNamesQuery {
	return &MetricNamesQuery{
		Kind: "MetricNamesQuery",
	}
}

type WorkspacesQuery struct {
	RawQuery     *string `json:"rawQuery,omitempty"`
	Kind         string  `json:"kind"`
	Subscription string  `json:"subscription"`
}

// NewWorkspacesQuery creates a new WorkspacesQuery object.
func NewWorkspacesQuery() *WorkspacesQuery {
	return &WorkspacesQuery{
		Kind: "WorkspacesQuery",
	}
}

type GrafanaTemplateVariableQuery = AppInsightsMetricNameQueryOrAppInsightsGroupByQueryOrSubscriptionsQueryOrResourceGroupsQueryOrResourceNamesQueryOrMetricNamespaceQueryOrMetricDefinitionsQueryOrMetricNamesQueryOrWorkspacesQueryOrUnknownQuery

// NewGrafanaTemplateVariableQuery creates a new GrafanaTemplateVariableQuery object.
func NewGrafanaTemplateVariableQuery() *GrafanaTemplateVariableQuery {
	return NewAppInsightsMetricNameQueryOrAppInsightsGroupByQueryOrSubscriptionsQueryOrResourceGroupsQueryOrResourceNamesQueryOrMetricNamespaceQueryOrMetricDefinitionsQueryOrMetricNamesQueryOrWorkspacesQueryOrUnknownQuery()
}

type AppInsightsMetricNameQueryOrAppInsightsGroupByQueryOrSubscriptionsQueryOrResourceGroupsQueryOrResourceNamesQueryOrMetricNamespaceQueryOrMetricDefinitionsQueryOrMetricNamesQueryOrWorkspacesQueryOrUnknownQuery struct {
	AppInsightsMetricNameQuery *AppInsightsMetricNameQuery `json:"AppInsightsMetricNameQuery,omitempty"`
	AppInsightsGroupByQuery    *AppInsightsGroupByQuery    `json:"AppInsightsGroupByQuery,omitempty"`
	SubscriptionsQuery         *SubscriptionsQuery         `json:"SubscriptionsQuery,omitempty"`
	ResourceGroupsQuery        *ResourceGroupsQuery        `json:"ResourceGroupsQuery,omitempty"`
	ResourceNamesQuery         *ResourceNamesQuery         `json:"ResourceNamesQuery,omitempty"`
	MetricNamespaceQuery       *MetricNamespaceQuery       `json:"MetricNamespaceQuery,omitempty"`
	MetricDefinitionsQuery     *MetricDefinitionsQuery     `json:"MetricDefinitionsQuery,omitempty"`
	MetricNamesQuery           *MetricNamesQuery           `json:"MetricNamesQuery,omitempty"`
	WorkspacesQuery            *WorkspacesQuery            `json:"WorkspacesQuery,omitempty"`
	UnknownQuery               *UnknownQuery               `json:"UnknownQuery,omitempty"`
}

// NewAppInsightsMetricNameQueryOrAppInsightsGroupByQueryOrSubscriptionsQueryOrResourceGroupsQueryOrResourceNamesQueryOrMetricNamespaceQueryOrMetricDefinitionsQueryOrMetricNamesQueryOrWorkspacesQueryOrUnknownQuery creates a new AppInsightsMetricNameQueryOrAppInsightsGroupByQueryOrSubscriptionsQueryOrResourceGroupsQueryOrResourceNamesQueryOrMetricNamespaceQueryOrMetricDefinitionsQueryOrMetricNamesQueryOrWorkspacesQueryOrUnknownQuery object.
func NewAppInsightsMetricNameQueryOrAppInsightsGroupByQueryOrSubscriptionsQueryOrResourceGroupsQueryOrResourceNamesQueryOrMetricNamespaceQueryOrMetricDefinitionsQueryOrMetricNamesQueryOrWorkspacesQueryOrUnknownQuery() *AppInsightsMetricNameQueryOrAppInsightsGroupByQueryOrSubscriptionsQueryOrResourceGroupsQueryOrResourceNamesQueryOrMetricNamespaceQueryOrMetricDefinitionsQueryOrMetricNamesQueryOrWorkspacesQueryOrUnknownQuery {
	return &AppInsightsMetricNameQueryOrAppInsightsGroupByQueryOrSubscriptionsQueryOrResourceGroupsQueryOrResourceNamesQueryOrMetricNamespaceQueryOrMetricDefinitionsQueryOrMetricNamesQueryOrWorkspacesQueryOrUnknownQuery{}
}

// MarshalJSON implements a custom JSON marshalling logic to encode `AppInsightsMetricNameQueryOrAppInsightsGroupByQueryOrSubscriptionsQueryOrResourceGroupsQueryOrResourceNamesQueryOrMetricNamespaceQueryOrMetricDefinitionsQueryOrMetricNamesQueryOrWorkspacesQueryOrUnknownQuery` as JSON.
func (resource AppInsightsMetricNameQueryOrAppInsightsGroupByQueryOrSubscriptionsQueryOrResourceGroupsQueryOrResourceNamesQueryOrMetricNamespaceQueryOrMetricDefinitionsQueryOrMetricNamesQueryOrWorkspacesQueryOrUnknownQuery) MarshalJSON() ([]byte, error) {
	if resource.AppInsightsMetricNameQuery != nil {
		return json.Marshal(resource.AppInsightsMetricNameQuery)
	}
	if resource.AppInsightsGroupByQuery != nil {
		return json.Marshal(resource.AppInsightsGroupByQuery)
	}
	if resource.SubscriptionsQuery != nil {
		return json.Marshal(resource.SubscriptionsQuery)
	}
	if resource.ResourceGroupsQuery != nil {
		return json.Marshal(resource.ResourceGroupsQuery)
	}
	if resource.ResourceNamesQuery != nil {
		return json.Marshal(resource.ResourceNamesQuery)
	}
	if resource.MetricNamespaceQuery != nil {
		return json.Marshal(resource.MetricNamespaceQuery)
	}
	if resource.MetricDefinitionsQuery != nil {
		return json.Marshal(resource.MetricDefinitionsQuery)
	}
	if resource.MetricNamesQuery != nil {
		return json.Marshal(resource.MetricNamesQuery)
	}
	if resource.WorkspacesQuery != nil {
		return json.Marshal(resource.WorkspacesQuery)
	}
	if resource.UnknownQuery != nil {
		return json.Marshal(resource.UnknownQuery)
	}

	return nil, fmt.Errorf("no value for disjunction of refs")
}

// UnmarshalJSON implements a custom JSON unmarshalling logic to decode `AppInsightsMetricNameQueryOrAppInsightsGroupByQueryOrSubscriptionsQueryOrResourceGroupsQueryOrResourceNamesQueryOrMetricNamespaceQueryOrMetricDefinitionsQueryOrMetricNamesQueryOrWorkspacesQueryOrUnknownQuery` from JSON.
func (resource *AppInsightsMetricNameQueryOrAppInsightsGroupByQueryOrSubscriptionsQueryOrResourceGroupsQueryOrResourceNamesQueryOrMetricNamespaceQueryOrMetricDefinitionsQueryOrMetricNamesQueryOrWorkspacesQueryOrUnknownQuery) UnmarshalJSON(raw []byte) error {
	if raw == nil {
		return nil
	}

	// FIXME: this is wasteful, we need to find a more efficient way to unmarshal this.
	parsedAsMap := make(map[string]any)
	if err := json.Unmarshal(raw, &parsedAsMap); err != nil {
		return err
	}

	discriminator, found := parsedAsMap["kind"]
	if !found {
		return errors.New("discriminator field 'kind' not found in payload")
	}

	switch discriminator {
	case "AppInsightsGroupByQuery":
		var appInsightsGroupByQuery AppInsightsGroupByQuery
		if err := json.Unmarshal(raw, &appInsightsGroupByQuery); err != nil {
			return err
		}

		resource.AppInsightsGroupByQuery = &appInsightsGroupByQuery
		return nil
	case "AppInsightsMetricNameQuery":
		var appInsightsMetricNameQuery AppInsightsMetricNameQuery
		if err := json.Unmarshal(raw, &appInsightsMetricNameQuery); err != nil {
			return err
		}

		resource.AppInsightsMetricNameQuery = &appInsightsMetricNameQuery
		return nil
	case "MetricDefinitionsQuery":
		var metricDefinitionsQuery MetricDefinitionsQuery
		if err := json.Unmarshal(raw, &metricDefinitionsQuery); err != nil {
			return err
		}

		resource.MetricDefinitionsQuery = &metricDefinitionsQuery
		return nil
	case "MetricNamesQuery":
		var metricNamesQuery MetricNamesQuery
		if err := json.Unmarshal(raw, &metricNamesQuery); err != nil {
			return err
		}

		resource.MetricNamesQuery = &metricNamesQuery
		return nil
	case "MetricNamespaceQuery":
		var metricNamespaceQuery MetricNamespaceQuery
		if err := json.Unmarshal(raw, &metricNamespaceQuery); err != nil {
			return err
		}

		resource.MetricNamespaceQuery = &metricNamespaceQuery
		return nil
	case "ResourceGroupsQuery":
		var resourceGroupsQuery ResourceGroupsQuery
		if err := json.Unmarshal(raw, &resourceGroupsQuery); err != nil {
			return err
		}

		resource.ResourceGroupsQuery = &resourceGroupsQuery
		return nil
	case "ResourceNamesQuery":
		var resourceNamesQuery ResourceNamesQuery
		if err := json.Unmarshal(raw, &resourceNamesQuery); err != nil {
			return err
		}

		resource.ResourceNamesQuery = &resourceNamesQuery
		return nil
	case "SubscriptionsQuery":
		var subscriptionsQuery SubscriptionsQuery
		if err := json.Unmarshal(raw, &subscriptionsQuery); err != nil {
			return err
		}

		resource.SubscriptionsQuery = &subscriptionsQuery
		return nil
	case "UnknownQuery":
		var unknownQuery UnknownQuery
		if err := json.Unmarshal(raw, &unknownQuery); err != nil {
			return err
		}

		resource.UnknownQuery = &unknownQuery
		return nil
	case "WorkspacesQuery":
		var workspacesQuery WorkspacesQuery
		if err := json.Unmarshal(raw, &workspacesQuery); err != nil {
			return err
		}

		resource.WorkspacesQuery = &workspacesQuery
		return nil
	}

	return fmt.Errorf("could not unmarshal resource with `kind = %v`", discriminator)
}
