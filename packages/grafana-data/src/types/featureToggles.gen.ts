// NOTE: This file was auto generated.  DO NOT EDIT DIRECTLY!
// To change feature flags, edit:
//  pkg/services/featuremgmt/registry.go
// Then run tests in:
//  pkg/services/featuremgmt/toggles_gen_test.go

/**
 * Describes available feature toggles in Grafana. These can be configured via
 * conf/custom.ini to enable features under development or not yet available in
 * stable version.
 *
 * Only enabled values will be returned in this interface.
 *
 * NOTE: the possible values may change between versions without notice, although
 * this may cause compilation issues when depending on removed feature keys, the
 * runtime state will continue to work.
 *
 * @public
 */
export interface FeatureToggles {
  disableEnvelopeEncryption?: boolean;
  ['live-service-web-worker']?: boolean;
  queryOverLive?: boolean;
  panelTitleSearch?: boolean;
  publicDashboardsEmailSharing?: boolean;
  publicDashboardsScene?: boolean;
  lokiExperimentalStreaming?: boolean;
  featureHighlights?: boolean;
  storage?: boolean;
  correlations?: boolean;
  canvasPanelNesting?: boolean;
  vizActions?: boolean;
  disableSecretsCompatibility?: boolean;
  logRequestsInstrumentedAsUnknown?: boolean;
  grpcServer?: boolean;
  cloudWatchCrossAccountQuerying?: boolean;
  showDashboardValidationWarnings?: boolean;
  mysqlAnsiQuotes?: boolean;
  accessControlOnCall?: boolean;
  nestedFolders?: boolean;
  alertingBacktesting?: boolean;
  editPanelCSVDragAndDrop?: boolean;
  logsContextDatasourceUi?: boolean;
  lokiShardSplitting?: boolean;
  lokiQuerySplitting?: boolean;
  lokiQuerySplittingConfig?: boolean;
  individualCookiePreferences?: boolean;
  influxdbBackendMigration?: boolean;
  influxqlStreamingParser?: boolean;
  influxdbRunQueriesInParallel?: boolean;
  prometheusRunQueriesInParallel?: boolean;
  lokiLogsDataplane?: boolean;
  dataplaneFrontendFallback?: boolean;
  disableSSEDataplane?: boolean;
  alertStateHistoryLokiSecondary?: boolean;
  alertStateHistoryLokiPrimary?: boolean;
  alertStateHistoryLokiOnly?: boolean;
  unifiedRequestLog?: boolean;
  renderAuthJWT?: boolean;
  refactorVariablesTimeRange?: boolean;
  faroDatasourceSelector?: boolean;
  enableDatagridEditing?: boolean;
  extraThemes?: boolean;
  lokiPredefinedOperations?: boolean;
  pluginsFrontendSandbox?: boolean;
  frontendSandboxMonitorOnly?: boolean;
  pluginsDetailsRightPanel?: boolean;
  sqlDatasourceDatabaseSelection?: boolean;
  recordedQueriesMulti?: boolean;
  logsExploreTableVisualisation?: boolean;
  awsDatasourcesTempCredentials?: boolean;
  transformationsRedesign?: boolean;
  mlExpressions?: boolean;
  traceQLStreaming?: boolean;
  metricsSummary?: boolean;
  datasourceAPIServers?: boolean;
  grafanaAPIServerWithExperimentalAPIs?: boolean;
  provisioning?: boolean;
  grafanaAPIServerEnsureKubectlAccess?: boolean;
  featureToggleAdminPage?: boolean;
  awsAsyncQueryCaching?: boolean;
  permissionsFilterRemoveSubquery?: boolean;
  configurableSchedulerTick?: boolean;
  alertingNoDataErrorExecution?: boolean;
  angularDeprecationUI?: boolean;
  dashgpt?: boolean;
  aiGeneratedDashboardChanges?: boolean;
  reportingRetries?: boolean;
  sseGroupByDatasource?: boolean;
  libraryPanelRBAC?: boolean;
  lokiRunQueriesInParallel?: boolean;
  wargamesTesting?: boolean;
  alertingInsights?: boolean;
  externalCorePlugins?: boolean;
  pluginsAPIMetrics?: boolean;
  externalServiceAccounts?: boolean;
  panelMonitoring?: boolean;
  enableNativeHTTPHistogram?: boolean;
  disableClassicHTTPHistogram?: boolean;
  formatString?: boolean;
  transformationsVariableSupport?: boolean;
  kubernetesPlaylists?: boolean;
  kubernetesSnapshots?: boolean;
  kubernetesDashboards?: boolean;
  kubernetesCliDashboards?: boolean;
  kubernetesRestore?: boolean;
  kubernetesFoldersServiceV2?: boolean;
  datasourceQueryTypes?: boolean;
  queryService?: boolean;
  queryServiceRewrite?: boolean;
  queryServiceFromUI?: boolean;
  cloudWatchBatchQueries?: boolean;
  recoveryThreshold?: boolean;
  lokiStructuredMetadata?: boolean;
  teamHttpHeaders?: boolean;
  cachingOptimizeSerializationMemoryUsage?: boolean;
  managedPluginsInstall?: boolean;
  prometheusPromQAIL?: boolean;
  prometheusCodeModeMetricNamesSearch?: boolean;
  addFieldFromCalculationStatFunctions?: boolean;
  alertmanagerRemoteSecondary?: boolean;
  alertmanagerRemotePrimary?: boolean;
  alertmanagerRemoteOnly?: boolean;
  annotationPermissionUpdate?: boolean;
  extractFieldsNameDeduplication?: boolean;
  dashboardSceneForViewers?: boolean;
  dashboardSceneSolo?: boolean;
  dashboardScene?: boolean;
  dashboardNewLayouts?: boolean;
  panelFilterVariable?: boolean;
  pdfTables?: boolean;
  ssoSettingsApi?: boolean;
  canvasPanelPanZoom?: boolean;
  logsInfiniteScrolling?: boolean;
  exploreMetrics?: boolean;
  alertingSimplifiedRouting?: boolean;
  logRowsPopoverMenu?: boolean;
  pluginsSkipHostEnvVars?: boolean;
  tableSharedCrosshair?: boolean;
  regressionTransformation?: boolean;
  lokiQueryHints?: boolean;
  kubernetesFeatureToggles?: boolean;
  cloudRBACRoles?: boolean;
  alertingQueryOptimization?: boolean;
  newFolderPicker?: boolean;
  jitterAlertRulesWithinGroups?: boolean;
  onPremToCloudMigrations?: boolean;
  alertingSaveStatePeriodic?: boolean;
  alertingSaveStateCompressed?: boolean;
  scopeApi?: boolean;
  promQLScope?: boolean;
  logQLScope?: boolean;
  sqlExpressions?: boolean;
  nodeGraphDotLayout?: boolean;
  groupToNestedTableTransformation?: boolean;
  newPDFRendering?: boolean;
  tlsMemcached?: boolean;
  kubernetesAggregator?: boolean;
  expressionParser?: boolean;
  groupByVariable?: boolean;
  scopeFilters?: boolean;
  ssoSettingsSAML?: boolean;
  oauthRequireSubClaim?: boolean;
  newDashboardWithFiltersAndGroupBy?: boolean;
  cloudWatchNewLabelParsing?: boolean;
  accessActionSets?: boolean;
  disableNumericMetricsSortingInExpressions?: boolean;
  grafanaManagedRecordingRules?: boolean;
  queryLibrary?: boolean;
  logsExploreTableDefaultVisualization?: boolean;
  newDashboardSharingComponent?: boolean;
  alertingListViewV2?: boolean;
  dashboardRestore?: boolean;
  datasourceProxyDisableRBAC?: boolean;
  alertingDisableSendAlertsExternal?: boolean;
  preserveDashboardStateWhenNavigating?: boolean;
  alertingCentralAlertHistory?: boolean;
  pluginProxyPreserveTrailingSlash?: boolean;
  sqlQuerybuilderFunctionParameters?: boolean;
  azureMonitorPrometheusExemplars?: boolean;
  pinNavItems?: boolean;
  authZGRPCServer?: boolean;
  ssoSettingsLDAP?: boolean;
  failWrongDSUID?: boolean;
  zanzana?: boolean;
  reloadDashboardsOnParamsChange?: boolean;
  enableScopesInMetricsExplore?: boolean;
  alertingApiServer?: boolean;
  cloudWatchRoundUpEndTime?: boolean;
  prometheusAzureOverrideAudience?: boolean;
  alertingFilterV2?: boolean;
  dataplaneAggregator?: boolean;
  newFiltersUI?: boolean;
  lokiSendDashboardPanelNames?: boolean;
  alertingPrometheusRulesPrimary?: boolean;
  exploreLogsShardSplitting?: boolean;
  exploreLogsAggregatedMetrics?: boolean;
  exploreLogsLimitedTimeRange?: boolean;
  homeSetupGuide?: boolean;
  appPlatformGrpcClientAuth?: boolean;
  appSidecar?: boolean;
  groupAttributeSync?: boolean;
  alertingQueryAndExpressionsStepMode?: boolean;
  improvedExternalSessionHandling?: boolean;
  useSessionStorageForRedirection?: boolean;
  rolePickerDrawer?: boolean;
  unifiedStorageSearch?: boolean;
  unifiedStorageSearchSprinkles?: boolean;
  unifiedStorageSearchPermissionFiltering?: boolean;
  pluginsSriChecks?: boolean;
  unifiedStorageBigObjectsSupport?: boolean;
  timeRangeProvider?: boolean;
  prometheusUsesCombobox?: boolean;
  userStorageAPI?: boolean;
  azureMonitorDisableLogLimit?: boolean;
  preinstallAutoUpdate?: boolean;
  playlistsReconciler?: boolean;
  passwordlessMagicLinkAuthentication?: boolean;
  exploreMetricsRelatedLogs?: boolean;
  prometheusSpecialCharsInLabelValues?: boolean;
  enableExtensionsAdminPage?: boolean;
  enableSCIM?: boolean;
  crashDetection?: boolean;
  jaegerBackendMigration?: boolean;
  reportingUseRawTimeRange?: boolean;
  alertingUIOptimizeReducer?: boolean;
  azureMonitorEnableUserAuth?: boolean;
  alertingNotificationsStepMode?: boolean;
  useV2DashboardsAPI?: boolean;
  feedbackButton?: boolean;
  unifiedStorageSearchUI?: boolean;
  elasticsearchCrossClusterSearch?: boolean;
  unifiedHistory?: boolean;
  lokiLabelNamesQueryApi?: boolean;
  investigationsBackend?: boolean;
  k8SFolderCounts?: boolean;
  k8SFolderMove?: boolean;
  improvedExternalSessionHandlingSAML?: boolean;
  teamHttpHeadersMimir?: boolean;
  ABTestFeatureToggleA?: boolean;
  templateVariablesUsesCombobox?: boolean;
  ABTestFeatureToggleB?: boolean;
  queryLibraryDashboards?: boolean;
  grafanaAdvisor?: boolean;
  elasticsearchImprovedParsing?: boolean;
  exploreMetricsUseExternalAppPlugin?: boolean;
  datasourceConnectionsTab?: boolean;
  fetchRulesUsingPost?: boolean;
  alertingConversionAPI?: boolean;
  alertingAlertmanagerExtraDedupStage?: boolean;
  alertingAlertmanagerExtraDedupStageStopPipeline?: boolean;
  newLogsPanel?: boolean;
  grafanaconThemes?: boolean;
  pluginsCDNSyncLoader?: boolean;
}
