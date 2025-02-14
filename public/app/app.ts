import 'symbol-observable';
import 'regenerator-runtime/runtime';

import 'whatwg-fetch'; // fetch polyfill needed for PhantomJs rendering
import 'file-saver';
import 'jquery';

import { createElement } from 'react';
import { createRoot } from 'react-dom/client';

import {
  locationUtil,
  monacoLanguageRegistry,
  setLocale,
  setTimeZoneResolver,
  setWeekStart,
  standardEditorsRegistry,
  standardFieldConfigEditorRegistry,
  standardTransformersRegistry,
} from '@grafana/data';
import {
  locationService,
  registerEchoBackend,
  setBackendSrv,
  setDataSourceSrv,
  setEchoSrv,
  setLocationSrv,
  setQueryRunnerFactory,
  setRunRequest,
  setPluginImportUtils,
  setPluginExtensionGetter,
  setEmbeddedDashboard,
  setAppEvents,
  setReturnToPreviousHook,
  setPluginExtensionsHook,
  setPluginComponentHook,
  setPluginComponentsHook,
  setCurrentUser,
  setChromeHeaderHeightHook,
  setPluginLinksHook,
  setCorrelationsService,
  setPluginFunctionsHook,
} from '@grafana/runtime';
import { setPanelDataErrorView } from '@grafana/runtime/src/components/PanelDataErrorView';
import { setPanelRenderer } from '@grafana/runtime/src/components/PanelRenderer';
import { setPluginPage } from '@grafana/runtime/src/components/PluginPage';
import config, { updateConfig } from 'app/core/config';
import { getStandardTransformers } from 'app/features/transformers/standardTransformers';

import getDefaultMonacoLanguages from '../lib/monaco-languages';

import { AppWrapper } from './AppWrapper';
import appEvents from './core/app_events';
import { AppChromeService } from './core/components/AppChrome/AppChromeService';
import { getAllOptionEditors, getAllStandardFieldConfigs } from './core/components/OptionsUI/registry';
import { PluginPage } from './core/components/Page/PluginPage';
import { GrafanaContextType, useChromeHeaderHeight, useReturnToPreviousInternal } from './core/context/GrafanaContext';
import { initializeCrashDetection } from './core/crash';
import { initializeI18n } from './core/internationalization';
import { setMonacoEnv } from './core/monacoEnv';
import { interceptLinkClicks } from './core/navigation/patch/interceptLinkClicks';
import { CorrelationsService } from './core/services/CorrelationsService';
import { NewFrontendAssetsChecker } from './core/services/NewFrontendAssetsChecker';
import { backendSrv } from './core/services/backend_srv';
import { contextSrv, RedirectToUrlKey } from './core/services/context_srv';
import { Echo } from './core/services/echo/Echo';
import { reportPerformance } from './core/services/echo/EchoSrv';
import { KeybindingSrv } from './core/services/keybindingSrv';
import { startMeasure, stopMeasure } from './core/utils/metrics';
import { initDevFeatures } from './dev';
import { initAlerting } from './features/alerting/unified/initAlerting';
import { initAuthConfig } from './features/auth-config';
import { getTimeSrv } from './features/dashboard/services/TimeSrv';
import { EmbeddedDashboardLazy } from './features/dashboard-scene/embedding/EmbeddedDashboardLazy';
import { initGrafanaLive } from './features/live';
import { PanelDataErrorView } from './features/panel/components/PanelDataErrorView';
import { PanelRenderer } from './features/panel/components/PanelRenderer';
import { DatasourceSrv } from './features/plugins/datasource_srv';
import { createPluginExtensionsGetter } from './features/plugins/extensions/getPluginExtensions';
import { pluginExtensionRegistries } from './features/plugins/extensions/registry/setup';
import { usePluginComponent } from './features/plugins/extensions/usePluginComponent';
import { usePluginComponents } from './features/plugins/extensions/usePluginComponents';
import { createUsePluginExtensions } from './features/plugins/extensions/usePluginExtensions';
import { usePluginFunctions } from './features/plugins/extensions/usePluginFunctions';
import { usePluginLinks } from './features/plugins/extensions/usePluginLinks';
import { getAppPluginsToAwait, getAppPluginsToPreload } from './features/plugins/extensions/utils';
import { importPanelPlugin, syncGetPanelPlugin } from './features/plugins/importPanelPlugin';
import { preloadPlugins } from './features/plugins/pluginPreloader';
import { QueryRunner } from './features/query/state/QueryRunner';
import { runRequest } from './features/query/state/runRequest';
import { initWindowRuntime } from './features/runtime/init';
import { initializeScopes } from './features/scopes';
import { cleanupOldExpandedFolders } from './features/search/utils';
import { variableAdapters } from './features/variables/adapters';
import { createAdHocVariableAdapter } from './features/variables/adhoc/adapter';
import { createConstantVariableAdapter } from './features/variables/constant/adapter';
import { createCustomVariableAdapter } from './features/variables/custom/adapter';
import { createDataSourceVariableAdapter } from './features/variables/datasource/adapter';
import { getVariablesUrlParams } from './features/variables/getAllVariableValuesForUrl';
import { createIntervalVariableAdapter } from './features/variables/interval/adapter';
import { setVariableQueryRunner, VariableQueryRunner } from './features/variables/query/VariableQueryRunner';
import { createQueryVariableAdapter } from './features/variables/query/adapter';
import { createSystemVariableAdapter } from './features/variables/system/adapter';
import { createTextBoxVariableAdapter } from './features/variables/textbox/adapter';
import { configureStore } from './store/configureStore';

// import symlinked extensions
const extensionsIndex = require.context('.', true, /extensions\/index.ts/);
const extensionsExports = extensionsIndex.keys().map((key) => {
  return extensionsIndex(key);
});

if (process.env.NODE_ENV === 'development') {
  initDevFeatures();
}

export class GrafanaApp {
  context!: GrafanaContextType;

  async init() {
    try {
      // Let iframe container know grafana has started loading
      parent.postMessage('GrafanaAppInit', '*');

      const initI18nPromise = initializeI18n(config.bootData.user.language);
      initI18nPromise.then(({ language }) => updateConfig({ language }));

      setBackendSrv(backendSrv);
      await initEchoSrv();
      // This needs to be done after the `initEchoSrv` since it is being used under the hood.
      startMeasure('frontend_app_init');

      setLocale(config.bootData.user.locale);
      setWeekStart(config.bootData.user.weekStart);
      setPanelRenderer(PanelRenderer);
      setPluginPage(PluginPage);
      setPanelDataErrorView(PanelDataErrorView);
      setLocationSrv(locationService);
      setCorrelationsService(new CorrelationsService());
      setEmbeddedDashboard(EmbeddedDashboardLazy);
      setTimeZoneResolver(() => config.bootData.user.timezone);
      initGrafanaLive();
      setCurrentUser(contextSrv.user);

      initAuthConfig();

      // Expose the app-wide eventbus
      setAppEvents(appEvents);

      // We must wait for translations to load because some preloaded store state requires translating
      await initI18nPromise;

      // Important that extension reducers are initialized before store
      addExtensionReducers();
      configureStore();
      initExtensions();

      initAlerting();

      standardEditorsRegistry.setInit(getAllOptionEditors);
      standardFieldConfigEditorRegistry.setInit(getAllStandardFieldConfigs);
      standardTransformersRegistry.setInit(getStandardTransformers);
      variableAdapters.setInit(() => [
        createQueryVariableAdapter(),
        createCustomVariableAdapter(),
        createTextBoxVariableAdapter(),
        createConstantVariableAdapter(),
        createDataSourceVariableAdapter(),
        createIntervalVariableAdapter(),
        createAdHocVariableAdapter(),
        createSystemVariableAdapter(),
      ]);

      monacoLanguageRegistry.setInit(getDefaultMonacoLanguages);
      setMonacoEnv();

      setQueryRunnerFactory(() => new QueryRunner());
      setVariableQueryRunner(new VariableQueryRunner());

      // Provide runRequest implementation to packages, @grafana/scenes in particular
      setRunRequest(runRequest);

      // Privide plugin import utils to packages, @grafana/scenes in particular
      setPluginImportUtils({
        importPanelPlugin,
        getPanelPluginFromCache: syncGetPanelPlugin,
      });

      if (config.featureToggles.useSessionStorageForRedirection) {
        handleRedirectTo();
      }

      locationUtil.initialize({
        config,
        getTimeRangeForUrl: getTimeSrv().timeRangeForUrl,
        getVariablesUrlParams: getVariablesUrlParams,
      });

      // intercept anchor clicks and forward it to custom history instead of relying on browser's history
      document.addEventListener('click', interceptLinkClicks);

      // Init DataSourceSrv
      const dataSourceSrv = new DatasourceSrv();
      dataSourceSrv.init(config.datasources, config.defaultDatasource);
      setDataSourceSrv(dataSourceSrv);
      initWindowRuntime();

      if (contextSrv.user.orgRole !== '') {
        const appPluginsToAwait = getAppPluginsToAwait();
        const appPluginsToPreload = getAppPluginsToPreload();

        preloadPlugins(appPluginsToPreload);
        await preloadPlugins(appPluginsToAwait);
      }

      setPluginExtensionGetter(createPluginExtensionsGetter(pluginExtensionRegistries));
      setPluginExtensionsHook(createUsePluginExtensions(pluginExtensionRegistries));
      setPluginLinksHook(usePluginLinks);
      setPluginComponentHook(usePluginComponent);
      setPluginComponentsHook(usePluginComponents);
      setPluginFunctionsHook(usePluginFunctions);

      // initialize chrome service
      const queryParams = locationService.getSearchObject();
      const chromeService = new AppChromeService();
      const keybindingsService = new KeybindingSrv(locationService, chromeService);
      const newAssetsChecker = new NewFrontendAssetsChecker();
      newAssetsChecker.start();

      // Read initial kiosk mode from url at app startup
      chromeService.setKioskModeFromUrl(queryParams.kiosk);

      // Clean up old search local storage values
      try {
        cleanupOldExpandedFolders();
      } catch (err) {
        console.warn('Failed to clean up old expanded folders', err);
      }

      this.context = {
        backend: backendSrv,
        location: locationService,
        chrome: chromeService,
        keybindings: keybindingsService,
        newAssetsChecker,
        config,
      };

      setReturnToPreviousHook(useReturnToPreviousInternal);
      setChromeHeaderHeightHook(useChromeHeaderHeight);

      initializeScopes();

      if (config.featureToggles.crashDetection) {
        initializeCrashDetection();
      }

      const root = createRoot(document.getElementById('reactRoot')!);
      root.render(
        createElement(AppWrapper, {
          app: this,
        })
      );
    } catch (error) {
      console.error('Failed to start Grafana', error);
      window.__grafana_load_failed();
    } finally {
      stopMeasure('frontend_app_init');
    }
  }
}

function addExtensionReducers() {
  if (extensionsExports.length > 0) {
    extensionsExports[0].addExtensionReducers();
  }
}

function initExtensions() {
  if (extensionsExports.length > 0) {
    extensionsExports[0].init();
  }
}

async function initEchoSrv() {
  setEchoSrv(new Echo({ debug: process.env.NODE_ENV === 'development' }));

  window.addEventListener('load', (e) => {
    const loadMetricName = 'frontend_boot_load_time_seconds';
    // Metrics below are marked in public/views/index.html
    const jsLoadMetricName = 'frontend_boot_js_done_time_seconds';
    const cssLoadMetricName = 'frontend_boot_css_time_seconds';

    if (performance) {
      performance.mark(loadMetricName);
      reportMetricPerformanceMark('first-paint', 'frontend_boot_', '_time_seconds');
      reportMetricPerformanceMark('first-contentful-paint', 'frontend_boot_', '_time_seconds');
      reportMetricPerformanceMark(loadMetricName);
      reportMetricPerformanceMark(jsLoadMetricName);
      reportMetricPerformanceMark(cssLoadMetricName);
    }
  });

  if (contextSrv.user.orgRole !== '') {
    const { PerformanceBackend } = await import('./core/services/echo/backends/PerformanceBackend');
    registerEchoBackend(new PerformanceBackend({}));
  }

  if (config.grafanaJavascriptAgent.enabled) {
    // Ignore Rudderstack URLs
    const rudderstackUrls = [
      config.rudderstackConfigUrl,
      config.rudderstackDataPlaneUrl,
      config.rudderstackIntegrationsUrl,
    ]
      .filter(Boolean)
      .map((url) => new RegExp(`${url}.*.`));

    const { GrafanaJavascriptAgentBackend } = await import(
      './core/services/echo/backends/grafana-javascript-agent/GrafanaJavascriptAgentBackend'
    );

    registerEchoBackend(
      new GrafanaJavascriptAgentBackend({
        ...config.grafanaJavascriptAgent,
        app: {
          version: config.buildInfo.version,
          environment: config.buildInfo.env,
        },
        buildInfo: config.buildInfo,
        user: {
          id: String(config.bootData.user?.id),
          email: config.bootData.user?.email,
        },
        ignoreUrls: rudderstackUrls,
      })
    );
  }

  if (config.googleAnalyticsId) {
    const { GAEchoBackend } = await import('./core/services/echo/backends/analytics/GABackend');
    registerEchoBackend(
      new GAEchoBackend({
        googleAnalyticsId: config.googleAnalyticsId,
      })
    );
  }

  if (config.googleAnalytics4Id) {
    const { GA4EchoBackend } = await import('./core/services/echo/backends/analytics/GA4Backend');
    registerEchoBackend(
      new GA4EchoBackend({
        googleAnalyticsId: config.googleAnalytics4Id,
        googleAnalytics4SendManualPageViews: config.googleAnalytics4SendManualPageViews,
      })
    );
  }

  if (config.rudderstackWriteKey && config.rudderstackDataPlaneUrl) {
    const { RudderstackBackend } = await import('./core/services/echo/backends/analytics/RudderstackBackend');
    registerEchoBackend(
      new RudderstackBackend({
        writeKey: config.rudderstackWriteKey,
        dataPlaneUrl: config.rudderstackDataPlaneUrl,
        user: config.bootData.user,
        sdkUrl: config.rudderstackSdkUrl,
        configUrl: config.rudderstackConfigUrl,
        integrationsUrl: config.rudderstackIntegrationsUrl,
        buildInfo: config.buildInfo,
      })
    );
  }

  if (config.applicationInsightsConnectionString) {
    const { ApplicationInsightsBackend } = await import(
      './core/services/echo/backends/analytics/ApplicationInsightsBackend'
    );
    registerEchoBackend(
      new ApplicationInsightsBackend({
        connectionString: config.applicationInsightsConnectionString,
        endpointUrl: config.applicationInsightsEndpointUrl,
      })
    );
  }

  if (config.analyticsConsoleReporting) {
    const { BrowserConsoleBackend } = await import('./core/services/echo/backends/analytics/BrowseConsoleBackend');
    registerEchoBackend(new BrowserConsoleBackend());
  }
}

/**
 * Report when a metric of a given name was marked during the document lifecycle. Works for markers with no duration,
 * like PerformanceMark or PerformancePaintTiming (e.g. created with performance.mark, or first-contentful-paint)
 */
function reportMetricPerformanceMark(metricName: string, prefix = '', suffix = ''): void {
  const metric = performance.getEntriesByName(metricName).at(0);
  if (metric) {
    const metricName = metric.name.replace(/-/g, '_');
    reportPerformance(`${prefix}${metricName}${suffix}`, Math.round(metric.startTime) / 1000);
  }
}

function handleRedirectTo(): void {
  const queryParams = locationService.getSearch();
  const redirectToParamKey = 'redirectTo';

  if (queryParams.has('auth_token')) {
    // URL Login should not be redirected
    window.sessionStorage.removeItem(RedirectToUrlKey);
    return;
  }

  if (queryParams.has(redirectToParamKey) && window.location.pathname !== '/') {
    const rawRedirectTo = queryParams.get(redirectToParamKey)!;
    window.sessionStorage.setItem(RedirectToUrlKey, encodeURIComponent(rawRedirectTo));
    queryParams.delete(redirectToParamKey);
    window.history.replaceState({}, '', `${window.location.pathname}${queryParams.size > 0 ? `?${queryParams}` : ''}`);
    return;
  }

  if (!contextSrv.user.isSignedIn) {
    return;
  }

  const redirectTo = window.sessionStorage.getItem(RedirectToUrlKey);
  if (!redirectTo) {
    return;
  }

  window.sessionStorage.removeItem(RedirectToUrlKey);
  const decodedRedirectTo = decodeURIComponent(redirectTo);
  if (decodedRedirectTo.startsWith('/goto/')) {
    // In this case there should be a request to the backend
    window.location.replace(decodedRedirectTo);
  } else {
    locationService.replace(decodedRedirectTo);
  }
}

export default new GrafanaApp();
