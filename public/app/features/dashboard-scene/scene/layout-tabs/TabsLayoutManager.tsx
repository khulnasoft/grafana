import { css } from '@emotion/css';

import { GrafanaTheme2 } from '@grafana/data';
import { SceneComponentProps, sceneGraph, SceneObjectBase, SceneObjectState, VizPanel } from '@grafana/scenes';
import { TabContent, TabsBar, useStyles2 } from '@grafana/ui';

import { DashboardScene } from '../DashboardScene';
import { ResponsiveGridLayoutManager } from '../layout-responsive-grid/ResponsiveGridLayoutManager';
import { DashboardLayoutManager, LayoutRegistryItem } from '../types';

import { TabItem } from './TabItem';
import { TabItemRepeaterBehavior } from './TabItemRepeaterBehavior';

interface TabsLayoutManagerState extends SceneObjectState {
  tabs: TabItem[];
  currentTab: TabItem;
}

export class TabsLayoutManager extends SceneObjectBase<TabsLayoutManagerState> implements DashboardLayoutManager {
  public isDashboardLayoutManager: true = true;

  public editModeChanged(isEditing: boolean): void {}

  public addPanel(vizPanel: VizPanel): void {
    this.state.currentTab.onAddPanel(vizPanel);
  }

  public addNewTab(): void {
    const currentTab = new TabItem({
      title: 'New tab',
      layout: ResponsiveGridLayoutManager.createEmpty(),
    });

    this.setState({
      tabs: [...this.state.tabs, currentTab],
      currentTab,
    });
  }

  public removeTab(tab: TabItem): void {
    if (this.state.tabs.length === 1) {
      throw new Error('TabsLayoutManager: Cannot remove last tab');
    }

    if (this.state.currentTab === tab) {
      const currentTabIndex = this.state.tabs.indexOf(tab);
      const nextTabIndex = currentTabIndex === 0 ? 1 : currentTabIndex - 1;
      const nextTab = this.state.tabs[nextTabIndex];
      this.setState({ tabs: this.state.tabs.filter((t) => t !== tab), currentTab: nextTab });
      return;
    }

    this.setState({
      tabs: this.state.tabs.filter((tab) => tab !== this.state.currentTab),
      currentTab: this.state.tabs[this.state.tabs.length - 1],
    });
  }

  public changeTab(tab: TabItem): void {
    this.setState({ currentTab: tab });
  }

  public addNewRow(): void {
    this.state.currentTab.getLayout().addNewRow();
  }

  public getMaxPanelId(): number {
    return Math.max(...this.state.tabs.map((tab) => tab.getLayout().getMaxPanelId()));
  }

  public getNextPanelId(): number {
    return 0;
  }

  public removePanel(panel: VizPanel) {
    this.state.currentTab.getLayout().removePanel(panel);
  }

  public duplicatePanel(panel: VizPanel): void {
    this.state.currentTab.getLayout();
    throw new Error('Method not implemented.');
  }

  public getVizPanels(): VizPanel[] {
    const panels: VizPanel[] = [];

    for (const tab of this.state.tabs) {
      const innerPanels = tab.getLayout().getVizPanels();
      panels.push(...innerPanels);
    }

    return panels;
  }

  public getOptions() {
    return [];
  }

  public activateRepeaters() {
    this.state.tabs.forEach((tab) => {
      if (tab.state.$behaviors) {
        for (const behavior of tab.state.$behaviors) {
          if (behavior instanceof TabItemRepeaterBehavior && !tab.isActive) {
            tab.activate();
            break;
          }
        }

        if (!tab.getLayout().isActive) {
          tab.getLayout().activate();
        }
      }
    });
  }

  public getDescriptor(): LayoutRegistryItem {
    return TabsLayoutManager.getDescriptor();
  }

  public getSelectedObject() {
    return sceneGraph.getAncestor(this, DashboardScene).state.editPane.state.selectedObject?.resolve();
  }

  public static getDescriptor(): LayoutRegistryItem {
    return {
      name: 'Tabs',
      description: 'Tabs layout',
      id: 'tabs-layout',
      createFromLayout: TabsLayoutManager.createFromLayout,
    };
  }

  public static createEmpty() {
    const tab = new TabItem({
      layout: ResponsiveGridLayoutManager.createEmpty(),
      title: 'Tab title',
    });
    return new TabsLayoutManager({ tabs: [tab], currentTab: tab });
  }

  public static createFromLayout(layout: DashboardLayoutManager): TabsLayoutManager {
    const tab = new TabItem({ layout: layout.clone(), title: 'Tab title' });
    return new TabsLayoutManager({ tabs: [tab], currentTab: tab });
  }

  public static Component = ({ model }: SceneComponentProps<TabsLayoutManager>) => {
    const styles = useStyles2(getStyles);
    const { tabs, currentTab } = model.useState();
    const { layout } = currentTab.useState();

    return (
      <>
        <TabsBar className={styles.tabsContainer}>
          {tabs.map((tab) => (
            <tab.Component model={tab} key={tab.state.key!} />
          ))}
        </TabsBar>
        <TabContent className={styles.tabContentContainer}>{layout && <layout.Component model={layout} />}</TabContent>
      </>
    );
  };
}

const getStyles = (theme: GrafanaTheme2) => ({
  tabsContainer: css({
    flexShrink: 1,
    padding: '2px 2px 0 2px',
    marginBottom: theme.spacing(1),
  }),
  tabContentContainer: css({
    backgroundColor: 'transparent',
    display: 'flex',
    flex: 1,
    height: '100%',
    overflow: 'auto',
    scrollbarWidth: 'thin',
    padding: '2px 2px 0 2px',
  }),
});
