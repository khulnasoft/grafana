import { defaults, each, sortBy, uniqBy } from 'lodash';

import { DataSourceRef, PanelPluginMeta, VariableOption, VariableRefresh } from '@grafana/data';
import { getDataSourceSrv } from '@grafana/runtime';
import {
  DashboardKind,
  DashboardImportableRequirements,
  DashboardV2Spec,
  ImportableResources,
  LibraryPanelImport,
  LibraryPanelKind,
  PanelKind,
  PanelQueryKind,
  AnnotationQueryKind,
  QueryVariableKind,
} from '@grafana/schema/dist/esm/schema/dashboard/v2alpha0';
import config from 'app/core/config';
import { PanelModel } from 'app/features/dashboard/state/PanelModel';
import { getLibraryPanel } from 'app/features/library-panels/state/api';
import { variableRegex } from 'app/features/variables/utils';

import { isPanelModelLibraryPanel } from '../../../library-panels/guard';
import { LibraryElementKind } from '../../../library-panels/types';
import { DashboardJson } from '../../../manage-dashboards/types';
import { isConstant } from '../../../variables/guard';
import { DashboardModel } from '../../state/DashboardModel';
import { GridPos } from '../../state/PanelModel';

export interface InputUsage {
  libraryPanels?: LibraryPanel[];
}

export interface LibraryPanel {
  name: string;
  uid: string;
}
export interface Input {
  name: string;
  type: string;
  label: string;
  value: any;
  description: string;
  usage?: InputUsage;
}

interface Requires {
  [key: string]: {
    type: string;
    id: string;
    name: string;
    version: string;
  };
}

export interface ExternalDashboard {
  __inputs?: Input[];
  __elements?: Record<string, LibraryElementExport>;
  __requires?: Array<Requires[string]>;
  panels: Array<PanelModel | PanelWithExportableLibraryPanel>;
}

interface PanelWithExportableLibraryPanel {
  gridPos: GridPos;
  id: number;
  libraryPanel: LibraryPanel;
}

function isExportableLibraryPanel(
  p: PanelModel | PanelWithExportableLibraryPanel
): p is PanelWithExportableLibraryPanel {
  return Boolean(p.libraryPanel?.name && p.libraryPanel?.uid);
}

interface DataSources {
  [key: string]: {
    name: string;
    label: string;
    description: string;
    type: string;
    pluginId: string;
    pluginName: string;
    usage?: InputUsage;
  };
}

export interface LibraryElementExport {
  name: string;
  uid: string;
  model: any;
  kind: LibraryElementKind;
}

export type DashboardV2Json = ImportableResources;

export interface DashboardExporterLike<T, J> {
  makeExportable(dashboard: T): Promise<J | { error: unknown }>;
}

export class DashboardExporterV1 implements DashboardExporterLike<DashboardModel, DashboardJson> {
  async makeExportable(dashboard: DashboardModel) {
    // clean up repeated rows and panels,
    // this is done on the live real dashboard instance, not on a clone
    // so we need to undo this
    // this is pretty hacky and needs to be changed
    dashboard.cleanUpRepeats();

    const saveModel = dashboard.getSaveModelCloneOld();
    saveModel.id = null;

    // undo repeat cleanup
    dashboard.processRepeats();

    const inputs: Input[] = [];
    const requires: Requires = {};
    const datasources: DataSources = {};
    const variableLookup: { [key: string]: any } = {};
    const libraryPanels: Map<string, LibraryElementExport> = new Map<string, LibraryElementExport>();

    for (const variable of saveModel.getVariables()) {
      variableLookup[variable.name] = variable;
    }

    const templateizeDatasourceUsage = (obj: any, fallback?: DataSourceRef) => {
      if (obj.datasource === undefined) {
        obj.datasource = fallback;
        return;
      }

      let datasource = obj.datasource;
      let datasourceVariable: any = null;

      const datasourceUid: string | undefined = datasource?.uid;
      const match = datasourceUid && variableRegex.exec(datasourceUid);

      // ignore data source properties that contain a variable
      if (match) {
        const varName = match[1] || match[2] || match[4];
        datasourceVariable = variableLookup[varName];
        if (datasourceVariable && datasourceVariable.current) {
          datasource = datasourceVariable.current.value;
        }
      }

      return getDataSourceSrv()
        .get(datasource)
        .then((ds) => {
          if (ds.meta?.builtIn) {
            return;
          }

          // add data source type to require list
          requires['datasource' + ds.meta?.id] = {
            type: 'datasource',
            id: ds.meta.id,
            name: ds.meta.name,
            version: ds.meta.info.version || '1.0.0',
          };

          // if used via variable we can skip templatizing usage
          if (datasourceVariable) {
            return;
          }

          const libraryPanel = obj.libraryPanel;
          const libraryPanelSuffix = !!libraryPanel ? '-for-library-panel' : '';
          let refName = 'DS_' + ds.name.replace(' ', '_').toUpperCase() + libraryPanelSuffix.toUpperCase();

          datasources[refName] = {
            name: refName,
            label: ds.name,
            description: '',
            type: 'datasource',
            pluginId: ds.meta?.id,
            pluginName: ds.meta?.name,
            usage: datasources[refName]?.usage,
          };

          if (!!libraryPanel) {
            const libPanels = datasources[refName]?.usage?.libraryPanels || [];
            libPanels.push({ name: libraryPanel.name, uid: libraryPanel.uid });

            datasources[refName].usage = {
              libraryPanels: libPanels,
            };
          }

          obj.datasource = { type: ds.meta.id, uid: '${' + refName + '}' };
        });
    };

    const processPanel = async (panel: PanelModel) => {
      if (panel.type !== 'row') {
        await templateizeDatasourceUsage(panel);

        if (panel.targets) {
          for (const target of panel.targets) {
            await templateizeDatasourceUsage(target, panel.datasource!);
          }
        }

        const panelDef: PanelPluginMeta = config.panels[panel.type];
        if (panelDef) {
          requires['panel' + panelDef.id] = {
            type: 'panel',
            id: panelDef.id,
            name: panelDef.name,
            version: panelDef.info.version,
          };
        }
      }
    };

    const processLibraryPanels = async (panel: PanelModel) => {
      if (isPanelModelLibraryPanel(panel)) {
        const { name, uid } = panel.libraryPanel;
        let model = panel.libraryPanel.model;
        if (!model) {
          const libPanel = await getLibraryPanel(uid, true);
          model = libPanel.model;
        }

        await templateizeDatasourceUsage(model);

        const { gridPos, id, ...rest } = model as any;
        if (!libraryPanels.has(uid)) {
          libraryPanels.set(uid, { name, uid, kind: LibraryElementKind.Panel, model: rest });
        }
      }
    };

    try {
      // check up panel data sources
      for (const panel of saveModel.panels) {
        await processPanel(panel);

        // handle collapsed rows
        if (panel.collapsed !== undefined && panel.collapsed === true && panel.panels) {
          for (const rowPanel of panel.panels) {
            await processPanel(rowPanel);
          }
        }
      }

      // templatize template vars
      for (const variable of saveModel.getVariables()) {
        if (variable.type === 'query') {
          await templateizeDatasourceUsage(variable);
          variable.options = [];
          variable.current = {} as unknown as VariableOption;
          variable.refresh =
            variable.refresh !== VariableRefresh.never ? variable.refresh : VariableRefresh.onDashboardLoad;
        } else if (variable.type === 'datasource') {
          variable.current = {};
        }
      }

      // templatize annotations vars
      for (const annotationDef of saveModel.annotations.list) {
        await templateizeDatasourceUsage(annotationDef);
      }

      // add grafana version
      requires['grafana'] = {
        type: 'grafana',
        id: 'grafana',
        name: 'Grafana',
        version: config.buildInfo.version,
      };

      // we need to process all panels again after all the promises are resolved
      // so all data sources, variables and targets have been templateized when we process library panels
      for (const panel of saveModel.panels) {
        await processLibraryPanels(panel);
        if (panel.collapsed !== undefined && panel.collapsed === true && panel.panels) {
          for (const rowPanel of panel.panels) {
            await processLibraryPanels(rowPanel);
          }
        }
      }

      each(datasources, (value: any) => {
        inputs.push(value);
      });

      // templatize constants
      for (const variable of saveModel.getVariables()) {
        if (isConstant(variable)) {
          const refName = 'VAR_' + variable.name.replace(' ', '_').toUpperCase();
          inputs.push({
            name: refName,
            type: 'constant',
            label: variable.label || variable.name,
            value: variable.query,
            description: '',
          });
          // update current and option
          variable.query = '${' + refName + '}';
          variable.current = {
            value: variable.query,
            text: variable.query,
            selected: false,
          };
          variable.options = [variable.current];
        }
      }

      const __elements = [...libraryPanels.entries()].reduce<Record<string, LibraryElementExport>>(
        (prev, [curKey, curLibPanel]) => {
          prev[curKey] = curLibPanel;
          return prev;
        },
        {}
      );

      // make inputs and requires a top thing
      const newObj: DashboardJson = defaults(
        {
          __inputs: inputs,
          __elements,
          __requires: sortBy(requires, ['id']),
        },
        saveModel
      );

      // Remove extraneous props from library panels
      for (let i = 0; i < newObj.panels.length; i++) {
        const libPanel = newObj.panels[i];
        if (isExportableLibraryPanel(libPanel)) {
          newObj.panels[i] = {
            gridPos: libPanel.gridPos,
            id: libPanel.id,
            libraryPanel: { uid: libPanel.libraryPanel.uid, name: libPanel.libraryPanel.name },
          };
        }
      }

      return newObj;
    } catch (err) {
      console.error('Export failed:', err);
      return {
        error: err,
      };
    }
  }
}

export class DashboardExporterV2 implements DashboardExporterLike<DashboardV2Spec, DashboardV2Json> {
  async makeExportable(dashboard: DashboardV2Spec) {
    const requires: DashboardImportableRequirements[] = [];
    const variableLookup: { [key: string]: any } = {};
    const libraryPanels: LibraryPanelImport[] = [];

    for (const variable of dashboard.variables) {
      variableLookup[variable.spec.name] = variable.spec;
    }

    const templateizeLibraryPanelDatasourceUsage = (obj: any) => {
      if (obj.datasource === undefined) {
        return;
      }

      let datasource = obj.datasource;
      const datasourceUid: string | undefined = obj.datasource.uid;
      const match = datasourceUid && variableRegex.exec(datasourceUid);

      const { datasourceVariable, datasourceMatch: dataSourceMatch } = getDatasourceFromMatch(match, datasource);

      if (dataSourceMatch) {
        datasource = dataSourceMatch;
      }

      return getDataSourceSrv()
        .get(datasource)
        .then((ds) => {
          if (ds.meta?.builtIn) {
            return;
          }

          requires.push({
            kind: 'datasource',
            spec: {
              id: ds.meta.id,
              name: ds.meta.name,
              version: ds.meta.info.version || '1.0.0',
            },
          });

          if (datasourceVariable) {
            return;
          }

          const libraryPanel = obj.libraryPanel;
          const libraryPanelSuffix = !!libraryPanel ? '-for-library-panel' : '';
          let refName = 'DS_' + ds.name.replace(' ', '_').toUpperCase() + libraryPanelSuffix.toUpperCase();
          obj.datasource = { type: ds.meta.id, uid: '${' + refName + '}' };
        });
    };

    const templateizeDatasourceUsage = (
      obj: AnnotationQueryKind['spec'] | QueryVariableKind['spec'] | PanelQueryKind['spec']
    ) => {
      if (obj.datasource === undefined) {
        return;
      }

      let datasource = obj.datasource;
      const datasourceUid: string | undefined = obj.datasource?.uid;
      const match = datasourceUid && variableRegex.exec(datasourceUid);

      // ignore data source properties that contain a variable
      const { datasourceMatch } = getDatasourceFromMatch(match, datasource);

      if (datasourceMatch) {
        datasource = datasourceMatch;
      }

      return getDataSourceSrv()
        .get(datasource)
        .then((ds) => {
          if (ds.meta?.builtIn) {
            return;
          }

          requires.push({
            kind: 'datasource',
            spec: {
              id: ds.meta.id,
              name: ds.meta.name,
              version: ds.meta.info.version || '1.0.0',
            },
          });
        });
    };

    const processPanel = async (panel: PanelKind) => {
      if (panel.spec.data.spec.queries) {
        for (const query of panel.spec.data.spec.queries) {
          await templateizeDatasourceUsage(query.spec);
        }
      }

      const panelDef: PanelPluginMeta = config.panels[panel.kind];
      if (panelDef) {
        requires.push({
          kind: 'panel',
          spec: {
            id: panelDef.id,
            name: panelDef.name,
            version: panelDef.info.version,
          },
        });
      }
    };

    const processLibraryPanels = async (panel: LibraryPanelKind) => {
      const { uid } = panel.spec.libraryPanel;

      const libPanel = await getLibraryPanel(uid, true);
      const exportableLibPanel: LibraryPanelImport = {
        kind: 'LibraryPanelImport',
        spec: {
          name: libPanel.name,
          uid: libPanel.uid,
          model: libPanel.model,
        },
      };

      await templateizeLibraryPanelDatasourceUsage(exportableLibPanel.spec.model);

      libraryPanels.push(exportableLibPanel);
    };

    try {
      // check panel data sources
      // support only grid layout for now
      if (dashboard.layout.kind !== 'GridLayout') {
        throw new Error('Only GridLayout is supported');
      }

      const gridLayoutItems = dashboard.layout.spec.items;
      const elements = dashboard.elements;

      for (const item of gridLayoutItems) {
        if (item.kind === 'GridLayoutItem' && !item.spec.repeat) {
          const panel = elements[item.spec.element.name];
          if (panel.kind === 'Panel') {
            await processPanel(panel);
          }
        }
      }

      // templatize template vars
      for (const variable of dashboard.variables) {
        if (variable.kind === 'QueryVariable') {
          await templateizeDatasourceUsage(variable.spec);
          variable.spec.options = [];
          variable.spec.current = {} as unknown as VariableOption;
        } else if (variable.kind === 'DatasourceVariable') {
          variable.spec.current = {
            text: '',
            value: '',
          };
        }
      }

      // templatize annotations vars
      for (const annotation of dashboard.annotations) {
        await templateizeDatasourceUsage(annotation.spec);
      }

      // add grafana version
      requires.push({
        kind: 'grafana',
        spec: {
          id: 'grafana',
          name: 'Grafana',
          version: config.buildInfo.version,
        },
      });

      for (const item of gridLayoutItems) {
        if (item.kind === 'GridLayoutItem' && !item.spec.repeat) {
          const panel = elements[item.spec.element.name];
          if (panel.kind === 'LibraryPanel') {
            await processLibraryPanels(panel);
          }
        }
      }

      const importableDashboard: DashboardKind = {
        kind: 'DashboardKind',
        spec: dashboard,
      };

      const newObj: DashboardV2Json = {
        kind: 'ImportableResources',
        spec: {
          resources: [importableDashboard, ...uniqBy(libraryPanels, 'spec.uid')],
          requirements: uniqBy(sortBy(requires, ['id']), 'spec.id'),
        },
      };

      return newObj;
    } catch (err) {
      console.error('Export failed:', err);
      return {
        error: err,
      };
    }

    function getDatasourceFromMatch(match: string | RegExpExecArray | null | undefined, datasource: any) {
      let datasourceVariable: any = null;
      if (match) {
        const varName = match[1] || match[2] || match[4];
        const datasourceVariable = variableLookup[varName];
        if (datasourceVariable && datasourceVariable.current) {
          datasource = datasourceVariable.current.value;
        }
      }
      return { datasourceVariable, datasourceMatch: datasource };
    }
  }
}

export function getDashboardExporter(): DashboardExporterLike<
  DashboardModel | DashboardV2Spec,
  DashboardJson | DashboardV2Json
> {
  return config.featureToggles.useV2DashboardsAPI ? new DashboardExporterV2() : new DashboardExporterV1();
}
