import { SceneGridItemLike, SceneGridRow, SceneObjectBase, SceneObjectState, VizPanel } from '@grafana/scenes';
import { t } from 'app/core/internationalization';

import { ConditionalRendering } from '../../conditional-rendering/ConditionalRendering';
import { DashboardOutlineItemType, DashboardOutlineRowItem } from '../../outline/types';
import { isClonedKey } from '../../utils/clone';
import { dashboardSceneGraph } from '../../utils/dashboardSceneGraph';
import { getDashboardSceneFor } from '../../utils/utils';
import { DashboardGridItem } from '../layout-default/DashboardGridItem';
import { DefaultGridLayoutManager } from '../layout-default/DefaultGridLayoutManager';
import { RowRepeaterBehavior } from '../layout-default/RowRepeaterBehavior';
import { TabsLayoutManager } from '../layout-tabs/TabsLayoutManager';
import { DashboardLayoutManager } from '../types/DashboardLayoutManager';
import { LayoutRegistryItem } from '../types/LayoutRegistryItem';

import { RowItem } from './RowItem';
import { RowItemRepeaterBehavior } from './RowItemRepeaterBehavior';
import { RowLayoutManagerRenderer } from './RowsLayoutManagerRenderer';

interface RowsLayoutManagerState extends SceneObjectState {
  rows: RowItem[];
}

export class RowsLayoutManager
  extends SceneObjectBase<RowsLayoutManagerState>
  implements DashboardLayoutManager<{}, DashboardOutlineRowItem>
{
  public static Component = RowLayoutManagerRenderer;

  public readonly isDashboardLayoutManager = true;

  public static readonly descriptor: LayoutRegistryItem = {
    get name() {
      return t('dashboard.rows-layout.name', 'Rows');
    },
    get description() {
      return t('dashboard.rows-layout.description', 'Rows layout');
    },
    id: 'rows-layout',
    createFromLayout: RowsLayoutManager.createFromLayout,

    kind: 'RowsLayout',
  };

  public readonly descriptor = RowsLayoutManager.descriptor;

  public addPanel(vizPanel: VizPanel) {
    // Try to add new panels to the selected row
    const selectedRows = dashboardSceneGraph.getAllSelectedObjects(this).filter((obj) => obj instanceof RowItem);
    if (selectedRows.length > 0) {
      return selectedRows.forEach((row) => row.onAddPanel(vizPanel));
    }

    // If we don't have selected row add it to the first row
    if (this.state.rows.length > 0) {
      return this.state.rows[0].onAddPanel(vizPanel);
    }

    // Otherwise fallback to adding a new row and a panel
    this.addNewRow();
    this.state.rows[this.state.rows.length - 1].onAddPanel(vizPanel);
  }

  public getVizPanels(): VizPanel[] {
    const panels: VizPanel[] = [];

    for (const row of this.state.rows) {
      const innerPanels = row.getLayout().getVizPanels();
      panels.push(...innerPanels);
    }

    return panels;
  }

  public hasVizPanels(): boolean {
    for (const row of this.state.rows) {
      if (row.getLayout().hasVizPanels()) {
        return true;
      }
    }

    return false;
  }

  public addNewRow() {
    this.setState({ rows: [...this.state.rows, new RowItem({ $behaviors: [ConditionalRendering.createEmpty()] })] });
  }

  public addNewTab() {
    const shouldAddTab = this.hasVizPanels();
    const tabsLayout = TabsLayoutManager.createFromLayout(this);

    if (shouldAddTab) {
      tabsLayout.addNewTab();
    }

    getDashboardSceneFor(this).switchLayout(tabsLayout);
  }

  public getOutline(): DashboardOutlineRowItem[] {
    return this.state.rows.map((row) => ({
      type: DashboardOutlineItemType.ROW,
      item: row,
      children: row.getLayout().getOutline(),
    }));
  }

  public editModeChanged(isEditing: boolean) {
    this.state.rows.forEach((row) => row.getLayout().editModeChanged?.(isEditing));
  }

  public activateRepeaters() {
    this.state.rows.forEach((row) => {
      if (!row.isActive) {
        row.activate();
      }

      const behavior = (row.state.$behaviors ?? []).find((b) => b instanceof RowItemRepeaterBehavior);

      if (!behavior?.isActive) {
        behavior?.activate();
      }

      row.getLayout().activateRepeaters?.();
    });
  }

  public addRowAbove(row: RowItem) {
    const rows = this.state.rows;
    const index = rows.indexOf(row);
    rows.splice(index, 0, new RowItem({ $behaviors: [ConditionalRendering.createEmpty()] }));
    this.setState({ rows });
  }

  public addRowBelow(row: RowItem) {
    const rows = this.state.rows;
    let index = rows.indexOf(row);

    // Be sure we don't add a row between an original row and one of its clones
    while (rows[index + 1] && isClonedKey(rows[index + 1].state.key!)) {
      index = index + 1;
    }

    rows.splice(index + 1, 0, new RowItem({ $behaviors: [ConditionalRendering.createEmpty()] }));
    this.setState({ rows });
  }

  public removeRow(row: RowItem) {
    const rows = this.state.rows.filter((r) => r !== row);
    this.setState({
      rows: rows.length === 0 ? [new RowItem({ $behaviors: [ConditionalRendering.createEmpty()] })] : rows,
    });
  }

  public moveRowUp(row: RowItem) {
    const rows = this.state.rows;
    const index = rows.indexOf(row);

    if (index === 0) {
      return;
    }

    rows.splice(index, 1);
    rows.splice(index - 1, 0, row);
    this.setState({ rows });
  }

  public moveRowDown(row: RowItem) {
    const rows = this.state.rows;
    let index = rows.indexOf(row);
    rows.splice(index, 1);

    // Be sure we don't add a row between an original row and one of its clones
    while (rows[index + 1] && isClonedKey(rows[index + 1].state.key!)) {
      index = index + 1;
    }

    rows.splice(index, 0, row);
    this.setState({ rows });
  }

  public isFirstRow(row: RowItem): boolean {
    return this.state.rows[0] === row;
  }

  public isLastRow(row: RowItem): boolean {
    const filteredRow = this.state.rows.filter((r) => !isClonedKey(r.state.key!));
    return filteredRow[filteredRow.length - 1] === row;
  }

  public static createEmpty(): RowsLayoutManager {
    return new RowsLayoutManager({
      rows: [new RowItem({ $behaviors: [ConditionalRendering.createEmpty()] })],
    });
  }

  public static createFromLayout(layout: DashboardLayoutManager): RowsLayoutManager {
    let rows: RowItem[];

    if (layout instanceof DefaultGridLayoutManager) {
      const config: Array<{
        title?: string;
        isCollapsed?: boolean;
        isDraggable?: boolean;
        isResizable?: boolean;
        children: SceneGridItemLike[];
        repeat?: string;
      }> = [];
      let children: SceneGridItemLike[] | undefined;

      layout.state.grid.forEachChild((child) => {
        if (!(child instanceof DashboardGridItem) && !(child instanceof SceneGridRow)) {
          throw new Error('Child is not a DashboardGridItem or SceneGridRow, invalid scene');
        }

        if (child instanceof SceneGridRow) {
          if (!isClonedKey(child.state.key!)) {
            const behaviour = child.state.$behaviors?.find((b) => b instanceof RowRepeaterBehavior);

            config.push({
              title: child.state.title,
              isCollapsed: !!child.state.isCollapsed,
              isDraggable: child.state.isDraggable ?? layout.state.grid.state.isDraggable,
              isResizable: child.state.isResizable ?? layout.state.grid.state.isResizable,
              children: child.state.children,
              repeat: behaviour?.state.variableName,
            });

            // Since we encountered a row item, any subsequent panels should be added to a new row
            children = undefined;
          }
        } else {
          if (!children) {
            children = [];
            config.push({ children });
          }

          children.push(child);
        }
      });

      rows = config.map(
        (rowConfig) =>
          new RowItem({
            title: rowConfig.title,
            isCollapsed: !!rowConfig.isCollapsed,
            layout: DefaultGridLayoutManager.fromGridItems(
              rowConfig.children,
              rowConfig.isDraggable,
              rowConfig.isResizable
            ),
            $behaviors: rowConfig.repeat
              ? [ConditionalRendering.createEmpty(), new RowItemRepeaterBehavior({ variableName: rowConfig.repeat })]
              : [ConditionalRendering.createEmpty()],
          })
      );
    } else {
      rows = [new RowItem({ layout: layout.clone(), $behaviors: [ConditionalRendering.createEmpty()] })];
    }

    // Ensure we always get at least one row
    if (rows.length === 0) {
      rows = [new RowItem({ $behaviors: [ConditionalRendering.createEmpty()] })];
    }

    return new RowsLayoutManager({ rows });
  }
}
