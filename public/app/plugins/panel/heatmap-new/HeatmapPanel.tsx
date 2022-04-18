import React, { useCallback, useMemo, useRef, useState } from 'react';
// import { css } from '@emotion/css';
// import { PanelProps, reduceField, ReducerID, TimeRange } from '@grafana/data';
import { DataFrame, PanelProps, reduceField, ReducerID, TimeRange } from '@grafana/data';
import { Portal, UPlotChart, useTheme2, VizLayout, LegendDisplayMode, usePanelContext } from '@grafana/ui';
import { PanelDataErrorView } from '@grafana/runtime';

import { HeatmapData, prepareHeatmapData } from './fields';
import { PanelOptions } from './models.gen';
import { quantizeScheme } from './palettes';
import {
  findExemplarFrameInPanelData,
  findDataFramesInPanelData,
  HeatmapHoverEvent,
  prepConfig,
  translateMatrixIndex,
  lookupDataInCell,
} from './utils';
import { HeatmapHoverView } from './HeatmapHoverView';
import { ColorScale } from './ColorScale';
import { HeatmapCalculationMode } from 'app/features/transformers/calculateHeatmap/models.gen';
import { HeatmapLookup } from './types';
import { HeatmapTab } from './hovertabs/HeatmapTab';
import { ExemplarTab } from './hovertabs/ExemplarTab';
import { ExemplarsPlugin } from './plugins/ExemplarsPlugin';

interface HeatmapPanelProps extends PanelProps<PanelOptions> {}

export const HeatmapPanel: React.FC<HeatmapPanelProps> = ({
  data,
  id,
  timeRange,
  timeZone,
  width,
  height,
  options,
  fieldConfig,
  onChangeTimeRange,
  replaceVariables,
}) => {
  const theme = useTheme2();
  // const styles = useStyles2(getStyles);

  // ugh
  let timeRangeRef = useRef<TimeRange>(timeRange);
  timeRangeRef.current = timeRange;

  const info = useMemo(
    () => prepareHeatmapData(findDataFramesInPanelData(data), options, theme),
    [data, options, theme]
  );
  const exemplars: HeatmapData | undefined = useMemo((): HeatmapData | undefined => {
    const exemplarsFrame: DataFrame | undefined = findExemplarFrameInPanelData(data);
    if (exemplarsFrame) {
      return prepareHeatmapData(
        [exemplarsFrame],
        {
          ...options,
          heatmap: {
            yAxis: {
              mode: HeatmapCalculationMode.Size,
              value: info.yBucketSize?.toString(),
            },
          },
        },
        theme
      );
    }
    return undefined;
  }, [data, info, options, theme]);
  const facets = useMemo(() => [null, info.heatmap?.fields.map((f) => f.values.toArray())], [info.heatmap]);
  const { onSplitOpen } = usePanelContext();

  const palette = useMemo(() => quantizeScheme(options.color, theme), [options.color, theme]);

  const [hover, setHover] = useState<HeatmapHoverEvent | undefined>(undefined);
  const [shouldDisplayCloseButton, setShouldDisplayCloseButton] = useState<boolean>(false);
  const isToolTipOpen = useRef<boolean>(false);

  const onCloseToolTip = () => {
    isToolTipOpen.current = false;
    setShouldDisplayCloseButton(false);
    onhover(null);
  };

  const onclick = () => {
    isToolTipOpen.current = !isToolTipOpen.current;
    setShouldDisplayCloseButton(isToolTipOpen.current);
  };

  const onhover = useCallback(
    (evt?: HeatmapHoverEvent | null) => {
      setHover(evt ?? undefined);
    },
    // eslint-disable-next-line react-hooks/exhaustive-deps
    [options, data.structureRev]
  );

  // ugh
  const dataRef = useRef<HeatmapData>(info);
  dataRef.current = info!;

  const builder = useMemo(() => {
    const b = prepConfig({
      dataRef,
      theme,
      onhover: onhover,
      onclick: options.tooltip.show ? onclick : null,
      onzoom: (evt) => {
        onChangeTimeRange({ from: evt.xMin, to: evt.xMax });
      },
      isToolTipOpen,
      timeZone,
      getTimeRange: () => timeRangeRef.current,
      palette,
      cellGap: options.cellGap,
      hideThreshold: options.hideThreshold,
    });

    b.addHook('draw', (u: uPlot) => {
      ExemplarsPlugin({
        u,
        exemplars: exemplars!,
        config: builder,
        theme: {
          ...theme,
          visualization: {
            ...theme.visualization,
            palette,
          },
        },
      });
    });
    return b;
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [options, data.structureRev]);

  const getExemplarValuesInCell = useCallback(
    (lookupRange: HeatmapLookup): DataFrame[] | undefined => {
      return lookupDataInCell(lookupRange, data.annotations?.[0]!, onSplitOpen, timeRange, timeZone);
    },
    [data.annotations, onSplitOpen, timeRange, timeZone]
  );

  const getDataValuesInCell = useCallback(
    (lookupRange: HeatmapLookup): DataFrame[] | undefined => {
      return lookupDataInCell(lookupRange, data.series?.[0]!, onSplitOpen, timeRange, timeZone);
    },
    [data.series, onSplitOpen, timeRange, timeZone]
  );

  const renderLegend = () => {
    if (options.legend.displayMode === LegendDisplayMode.Hidden || !info.heatmap) {
      return null;
    }

    const field = info.heatmap.fields[2];
    const { min, max } = reduceField({ field, reducers: [ReducerID.min, ReducerID.max] });

    let hoverValue: number | undefined = undefined;
    if (hover && info.heatmap.fields) {
      const countField = info.heatmap.fields[2];
      hoverValue = countField?.values.get(hover.index);
    }

    return (
      <VizLayout.Legend placement="bottom" maxHeight="20%">
        <ColorScale hoverValue={hoverValue} colorPalette={palette} min={min} max={max} display={info.display} />
      </VizLayout.Legend>
    );
  };

  if (info.warning || !info.heatmap) {
    return (
      <PanelDataErrorView
        panelId={id}
        fieldConfig={fieldConfig}
        data={data}
        needsNumberField={true}
        message={info.warning}
      />
    );
  }

  return (
    <>
      <VizLayout width={width} height={height} legend={renderLegend()}>
        {(vizWidth: number, vizHeight: number) => (
          <UPlotChart config={builder} data={facets as any} width={vizWidth} height={vizHeight} timeRange={timeRange} />
        )}
      </VizLayout>
      <Portal>
        {hover && (
          <HeatmapHoverView
            ttip={{
              layers: [
                HeatmapTab({
                  heatmapData: info,
                  index: hover.index,
                  getValuesInCell: getDataValuesInCell,
                  options: { showHistogram: options.tooltip.yHistogram, timeZone },
                }),
                ExemplarTab({
                  heatmapData: exemplars!,
                  getValuesInCell: getExemplarValuesInCell,
                  index: translateMatrixIndex(hover.index, info.yBucketCount!, exemplars?.yBucketCount!),
                }),
              ],
              hover,
              point: {},
            }}
            isOpen={shouldDisplayCloseButton}
            onClose={onCloseToolTip}
          />
        )}
      </Portal>
    </>
  );
};
