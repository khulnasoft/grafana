package advisor

manifest: {
	appName:       "advisor"
	groupOverride: "advisor.grafana.app"
	kinds: [
		datasourcecheck,
		plugincheck,
	]
}
