import { DataFrame, getFieldDisplayName, getFieldSeriesColor } from '@grafana/data';
import { VizLegendOptions, AxisPlacement } from '@grafana/schema';

import { useTheme2 } from '../../themes';
import { VizLayoutLegendProps } from '../VizLayout/VizLayout';
import { getDisplayValuesForCalcs } from '../uPlot/utils';

import { VizLegend } from './VizLegend';

interface LegendProps extends VizLegendOptions, Omit<VizLayoutLegendProps, 'children'> {
  data: DataFrame[];
}

export const VizLegend2 = ({ data, placement, calcs, displayMode, ...vizLayoutLegendProps }: LegendProps) => {
  const theme = useTheme2();

  const items = data.map((frame) => {
    const field = frame.fields[1];
    const fieldIndex = field.state?.origin;

    const label = getFieldDisplayName(field, frame, data);
    const color = getFieldSeriesColor(field, theme).color;
    const yAxis = field.config.custom?.axisPlacement === AxisPlacement.Right ? 2 : 1;

    return {
      fieldIndex, // { frameIndex, fieldIndex } - didn't find a usage of fieldIndex; it's only used in getItemKey
      label,
      color,
      yAxis,
      getDisplayValues: () => getDisplayValuesForCalcs(calcs, field, theme),
      getItemKey: () => `${label}-${fieldIndex?.frameIndex}-${fieldIndex?.fieldIndex}`,
    };
  });

  return (
    <VizLegend
      placement={placement}
      items={items}
      displayMode={displayMode}
      sortBy={vizLayoutLegendProps.sortBy}
      sortDesc={vizLayoutLegendProps.sortDesc}
      isSortable={true}
    />
  );
};
