import { PanelBuilders, SceneFlexItem, SceneQueryRunner } from '@grafana/scenes';
import { DataSourceRef, GraphDrawStyle, TooltipDisplayMode } from '@grafana/schema';

import { INSTANCE_ID, PANEL_STYLES, overrideToFixedColor } from '../../home/Insights';
import { InsightsMenuButton } from '../InsightsMenuButton';

export function getGrafanaEvalSuccessVsFailuresScene(datasource: DataSourceRef, panelTitle: string) {
  const exprA = INSTANCE_ID
    ? `sum(grafanacloud_grafana_instance_alerting_rule_evaluations_total:rate5m{id="${INSTANCE_ID}"}) - sum(grafanacloud_grafana_instance_alerting_rule_evaluation_failures_total:rate5m{id="${INSTANCE_ID}"})`
    : `sum(grafanacloud_grafana_instance_alerting_rule_evaluations_total:rate5m) - sum(grafanacloud_grafana_instance_alerting_rule_evaluation_failures_total:rate5m)`;

  const exprB = INSTANCE_ID
    ? `sum(grafanacloud_grafana_instance_alerting_rule_evaluation_failures_total:rate5m{id="${INSTANCE_ID}"})`
    : `sum(grafanacloud_grafana_instance_alerting_rule_evaluation_failures_total:rate5m)`;

  const query = new SceneQueryRunner({
    datasource,
    queries: [
      {
        refId: 'A',
        expr: exprA,
        range: true,
        legendFormat: 'success',
      },
      {
        refId: 'B',
        expr: exprB,
        range: true,
        legendFormat: 'failed',
      },
    ],
  });

  return new SceneFlexItem({
    ...PANEL_STYLES,
    body: PanelBuilders.timeseries()
      .setTitle(panelTitle)
      .setDescription('The number of successful and failed alert rule evaluations')
      .setData(query)
      .setOption('tooltip', { mode: TooltipDisplayMode.Multi })
      .setCustomFieldConfig('drawStyle', GraphDrawStyle.Line)
      .setOverrides((b) =>
        b
          .matchFieldsWithName('success')
          .overrideColor(overrideToFixedColor('success'))
          .matchFieldsWithName('failed')
          .overrideColor(overrideToFixedColor('failed'))
      )
      .setHeaderActions([new InsightsMenuButton({ panel: panelTitle })])
      .build(),
  });
}
