// Code generated - EDITING IS FUTILE. DO NOT EDIT.
//
// Generated by:
//     public/app/plugins/gen.go
// Using jennies:
//     PluginTsTypesJenny
//     PluginTsTypesJenny
//
// Run 'make gen-cue' from repository root to regenerate.

import * as common from '@grafana/schema';

export interface Options extends common.OptionsWithTimezones {
  legend: common.VizLegendOptions;
  orientation?: common.VizOrientation;
  tooltip: common.VizTooltipOptions;
}

export interface FieldConfig extends common.GraphFieldConfig {}
