import { Property } from 'csstype';
import tinycolor from 'tinycolor2';

import {
  FieldType,
  Field,
  formattedValueToString,
  reduceField,
  GrafanaTheme2,
  DisplayValue,
  LinkModel,
  DisplayValueAlignmentFactors,
  DataFrame,
} from '@grafana/data';
import {
  TableCellBackgroundDisplayMode,
  TableCellDisplayMode,
  TableCellOptions,
  TableAutoCellOptions,
  BarGaugeDisplayMode,
} from '@grafana/schema';

import { getTextColorForAlphaBackground } from '../../../utils';

import { CellColors, TableRow, TableFieldOptionsType, TableNGProps } from './types';
import { COLUMN } from './constants';

export function getCellHeight(
  text: string,
  cellWidth: number, // width of the cell without padding
  osContext: OffscreenCanvasRenderingContext2D | null,
  lineHeight: number,
  defaultRowHeight: number,
  padding = 0
) {
  const PADDING = padding * 2;

  if (osContext !== null && typeof text === 'string') {
    const words = text.split(/\s/);
    const lines = [];
    let currentLine = '';

    // Let's just wrap the lines and see how well the measurement works
    for (let i = 0; i < words.length; i++) {
      const currentWord = words[i];
      // TODO: this method is not accurate
      let lineWidth = osContext.measureText(currentLine + ' ' + currentWord).width;

      // if line width is less than the cell width, add the word to the current line and continue
      // else add the current line to the lines array and start a new line with the current word
      if (lineWidth < cellWidth) {
        currentLine += ' ' + currentWord;
      } else {
        lines.push({
          width: lineWidth,
          line: currentLine,
        });

        currentLine = currentWord;
      }

      // if we are at the last word, add the current line to the lines array
      if (i === words.length - 1) {
        lines.push({
          width: lineWidth,
          line: currentLine,
        });
      }
    }

    if (lines.length === 1) {
      return defaultRowHeight;
    }

    // TODO: double padding to adjust osContext.measureText() results
    const height = lines.length * lineHeight + PADDING * 2;

    return height;
  }

  return defaultRowHeight;
}

/**
 * getRowHeight determines cell height based on cell width + text length. Used
 * for when textWrap is enabled.
 */
export function getRowHeight(
  row: Record<string, string>,
  columnTypes: Record<string, string>,
  headerCellRefs: React.MutableRefObject<Record<string, HTMLDivElement>>,
  osContext: OffscreenCanvasRenderingContext2D | null,
  lineHeight: number,
  defaultRowHeight: number,
  padding: number
): number {
  /**
   * 0. loop through all cells in row
   * 1. find text cell in row
   * 2. find width of text cell
   * 3. calculate height based on width and text length
   * 4. return biggest height
   */

  let biggestHeight = defaultRowHeight;

  for (const key in row) {
    if (isTextCell(key, columnTypes)) {
      if (Object.keys(headerCellRefs.current).length === 0 || !headerCellRefs.current[key]) {
        return biggestHeight;
      }
      const cellWidth = headerCellRefs.current[key].offsetWidth;
      const cellText = row[key];
      const newCellHeight = getCellHeight(cellText, cellWidth, osContext, lineHeight, defaultRowHeight, padding);

      if (newCellHeight > biggestHeight) {
        biggestHeight = newCellHeight;
      }
    }
  }

  return biggestHeight;
}

function isTextCell(key: string, columnTypes: Record<string, string>): boolean {
  return columnTypes[key] === FieldType.string;
}

export function shouldTextOverflow(
  key: string,
  row: Record<string, string>,
  columnTypes: Record<string, string>,
  headerCellRefs: React.MutableRefObject<Record<string, HTMLDivElement>>,
  osContext: OffscreenCanvasRenderingContext2D | null,
  lineHeight: number,
  defaultRowHeight: number,
  padding: number,
  textWrap: boolean,
  cellInspect: boolean
): boolean {
  if (textWrap || cellInspect) {
    return false;
  }

  if (isTextCell(key, columnTypes)) {
    const cellWidth = headerCellRefs.current[key].offsetWidth;
    const cellText = row[key];
    const newCellHeight = getCellHeight(cellText, cellWidth, osContext, lineHeight, defaultRowHeight, padding);

    if (newCellHeight > defaultRowHeight) {
      return true;
    }
  }

  return false;
}

export interface TableFooterCalc {
  show: boolean;
  reducer: string[]; // actually 1 value
  fields?: string[];
  enablePagination?: boolean;
  countRows?: boolean;
}

export function getFooterItemNG(rows: TableRow[], field: Field, options: TableFooterCalc | undefined): string {
  if (!options) {
    return '';
  }

  if (field.type !== FieldType.number) {
    return '';
  }

  const calc = options.reducer[0];
  if (calc === undefined) {
    return '';
  }

  const value = reduceField({
    field: {
      ...field,
      values: rows.map((row) => row[field.name]),
    },
    reducers: options.reducer,
  })[calc];

  const formattedValue = formattedValueToString(field.display!(value));

  return formattedValue;
}

const CELL_COLOR_DARKENING_MULTIPLIER = 10;
const CELL_GRADIENT_DARKENING_MULTIPLIER = 15;
const CELL_GRADIENT_HUE_ROTATION_DEGREES = 5;

export function getCellColors(
  theme: GrafanaTheme2,
  cellOptions: TableCellOptions,
  displayValue: DisplayValue
): CellColors {
  // Convert RGBA hover color to hex to prevent transparency issues on cell hover
  const autoCellBackgroundHoverColor = convertRGBAToHex(theme.colors.background.primary, theme.colors.action.hover);

  // How much to darken elements depends upon if we're in dark mode
  const darkeningFactor = theme.isDark ? 1 : -0.7;

  // Setup color variables
  let textColor: string | undefined = undefined;
  let bgColor: string | undefined = undefined;
  let bgHoverColor: string = autoCellBackgroundHoverColor;

  if (cellOptions.type === TableCellDisplayMode.ColorText) {
    textColor = displayValue.color;
  } else if (cellOptions.type === TableCellDisplayMode.ColorBackground) {
    const mode = cellOptions.mode ?? TableCellBackgroundDisplayMode.Gradient;

    if (mode === TableCellBackgroundDisplayMode.Basic) {
      textColor = getTextColorForAlphaBackground(displayValue.color!, theme.isDark);
      bgColor = tinycolor(displayValue.color).toRgbString();
      bgHoverColor = tinycolor(displayValue.color)
        .darken(CELL_COLOR_DARKENING_MULTIPLIER * darkeningFactor)
        .toRgbString();
    } else if (mode === TableCellBackgroundDisplayMode.Gradient) {
      const hoverColor = tinycolor(displayValue.color)
        .darken(CELL_GRADIENT_DARKENING_MULTIPLIER * darkeningFactor)
        .toRgbString();
      const bgColor2 = tinycolor(displayValue.color)
        .darken(CELL_COLOR_DARKENING_MULTIPLIER * darkeningFactor)
        .spin(CELL_GRADIENT_HUE_ROTATION_DEGREES);
      textColor = getTextColorForAlphaBackground(displayValue.color!, theme.isDark);
      bgColor = `linear-gradient(120deg, ${bgColor2.toRgbString()}, ${displayValue.color})`;
      bgHoverColor = `linear-gradient(120deg, ${bgColor2.toRgbString()}, ${hoverColor})`;
    }
  }

  return { textColor, bgColor, bgHoverColor };
}

/**
 * @internal
 */
export const getCellLinks = (field: Field, rowIdx: number) => {
  let links: Array<LinkModel<unknown>> | undefined;
  if (field.getLinks) {
    links = field.getLinks({
      valueRowIndex: rowIdx,
    });
  }

  if (!links) {
    return;
  }

  for (let i = 0; i < links?.length; i++) {
    if (links[i].onClick) {
      const origOnClick = links[i].onClick;

      links[i].onClick = (event) => {
        // Allow opening in new tab
        if (!(event.ctrlKey || event.metaKey || event.shiftKey)) {
          event.preventDefault();
          origOnClick!(event, {
            field,
            rowIndex: rowIdx,
          });
        }
      };
    }
  }

  return links;
};

/** Extracts numeric pixel value from theme spacing */
export const extractPixelValue = (spacing: string | number): number => {
  return typeof spacing === 'number' ? spacing : parseFloat(spacing) || 0;
};

export function getTextAlign(field?: Field): Property.JustifyContent {
  if (!field) {
    return 'flex-start';
  }

  if (field.config.custom) {
    const custom: TableFieldOptionsType = field.config.custom;

    switch (custom.align) {
      case 'right':
        return 'flex-end';
      case 'left':
        return 'flex-start';
      case 'center':
        return 'center';
    }
  }

  if (field.type === FieldType.number) {
    return 'flex-end';
  }

  return 'flex-start';
}

const defaultCellOptions: TableAutoCellOptions = { type: TableCellDisplayMode.Auto };

export function getCellOptions(field: Field): TableCellOptions {
  if (field.config.custom?.displayMode) {
    return migrateTableDisplayModeToCellOptions(field.config.custom?.displayMode);
  }

  if (!field.config.custom?.cellOptions) {
    return defaultCellOptions;
  }

  return field.config.custom.cellOptions;
}

/**
 * Migrates table cell display mode to new object format.
 *
 * @param displayMode The display mode of the cell
 * @returns TableCellOptions object in the correct format
 * relative to the old display mode.
 */
export function migrateTableDisplayModeToCellOptions(displayMode: TableCellDisplayMode): TableCellOptions {
  switch (displayMode) {
    // In the case of the gauge we move to a different option
    case 'basic':
    case 'gradient-gauge':
    case 'lcd-gauge':
      let gaugeMode = BarGaugeDisplayMode.Basic;

      if (displayMode === 'gradient-gauge') {
        gaugeMode = BarGaugeDisplayMode.Gradient;
      } else if (displayMode === 'lcd-gauge') {
        gaugeMode = BarGaugeDisplayMode.Lcd;
      }

      return {
        type: TableCellDisplayMode.Gauge,
        mode: gaugeMode,
      };
    // Also true in the case of the color background
    case 'color-background':
    case 'color-background-solid':
      let mode = TableCellBackgroundDisplayMode.Basic;

      // Set the new mode field, somewhat confusingly the
      // color-background mode is for gradient display
      if (displayMode === 'color-background') {
        mode = TableCellBackgroundDisplayMode.Gradient;
      }

      return {
        type: TableCellDisplayMode.ColorBackground,
        mode: mode,
      };
    default:
      return {
        // @ts-ignore
        type: displayMode,
      };
  }
}

/**
 * Getting gauge or sparkline values to align is very tricky without looking at all values and passing them through display processor.
 * For very large tables that could pretty expensive. So this is kind of a compromise. We look at the first 1000 rows and cache the longest value.
 * If we have a cached value we just check if the current value is longer and update the alignmentFactor. This can obviously still lead to
 * unaligned gauges but it should a lot less common.
 **/
export function getAlignmentFactor(
  field: Field,
  displayValue: DisplayValue,
  rowIndex: number
): DisplayValueAlignmentFactors {
  let alignmentFactor = field.state?.alignmentFactors;

  if (alignmentFactor) {
    // check if current alignmentFactor is still the longest
    if (formattedValueToString(alignmentFactor).length < formattedValueToString(displayValue).length) {
      alignmentFactor = { ...displayValue };
      field.state!.alignmentFactors = alignmentFactor;
    }
    return alignmentFactor;
  } else {
    // look at the next 1000 rows
    alignmentFactor = { ...displayValue };
    const maxIndex = Math.min(field.values.length, rowIndex + 1000);

    for (let i = rowIndex + 1; i < maxIndex; i++) {
      const nextDisplayValue = field.display!(field.values[i]);
      if (formattedValueToString(alignmentFactor).length > formattedValueToString(nextDisplayValue).length) {
        alignmentFactor.text = displayValue.text;
      }
    }

    if (field.state) {
      field.state.alignmentFactors = alignmentFactor;
    } else {
      field.state = { alignmentFactors: alignmentFactor };
    }

    return alignmentFactor;
  }
}

/** Converts an RGBA color to hex by blending it with a background color */
function convertRGBAToHex(backgroundColor: string, rgbaColor: string): string {
  const bg = tinycolor(backgroundColor);
  const rgba = tinycolor(rgbaColor);
  return tinycolor.mix(bg, rgba, rgba.getAlpha() * 100).toHexString();
}

/** Returns true if the DataFrame contains nested frames */
export const getIsNestedTable = (dataFrame: DataFrame): boolean =>
  dataFrame.fields.some(({ type }) => type === FieldType.nestedFrames);

/** Returns the width of a column based on overrides, field config, and defaults */
export const getColumnWidth = (field: Field, fieldConfig: TableNGProps['fieldConfig'], key: string): number => {
  const overrideWidth = fieldConfig?.overrides
    ?.find(({ matcher: { id, options } }) => id === 'byName' && options === key)
    ?.properties?.find(({ id }) => id === 'width')?.value;

  return overrideWidth ?? field.config?.custom?.width ?? fieldConfig?.defaults?.custom?.width ?? COLUMN.DEFAULT_WIDTH;
};
