// Code generated - EDITING IS FUTILE. DO NOT EDIT.
//
// Generated by:
//     public/app/plugins/gen.go
// Using jennies:
//     TSTypesJenny
//     PluginTsTypesJenny
//
// Run 'make gen-cue' from repository root to regenerate.

import * as common from '@grafana/schema';

export const pluginVersion = "11.6.0-pre";

export interface Options {
  dedupStrategy: common.LogsDedupStrategy;
  enableInfiniteScrolling?: boolean;
  enableLogDetails: boolean;
  onNewLogsReceived?: unknown;
  showTime: boolean;
  sortOrder: common.LogsSortOrder;
  wrapLogMessage: boolean;
}
