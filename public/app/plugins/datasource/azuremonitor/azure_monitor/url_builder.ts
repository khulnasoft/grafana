import { TemplateSrv } from '@grafana/runtime';

import { AzureMonitorResource, GetMetricNamespacesQuery, GetMetricNamesQuery } from '../types';

export default class UrlBuilder {
  static buildResourceUri(templateSrv: TemplateSrv, resource: AzureMonitorResource, multipleResources?: boolean) {
    const urlArray = [];
    const { subscription, resourceGroup, metricNamespace, resourceName } = resource;
    if (subscription) {
      urlArray.push('/subscriptions', subscription);
      if (resourceGroup && !multipleResources) {
        urlArray.push('resourceGroups', resourceGroup);

        if (metricNamespace && resourceName) {
          const metricNamespaceProcessed = templateSrv.replace(metricNamespace);
          const metricNamespaceArray = metricNamespace.split('/');
          const resourceNameProcessed = templateSrv.replace(resourceName);
          const resourceNameArray = resourceName.split('/');
          const provider = metricNamespaceArray.shift();
          if (provider) {
            urlArray.push('providers', provider);
          }

          if (
            metricNamespaceProcessed.toLowerCase().startsWith('microsoft.storage/storageaccounts/') &&
            !resourceNameProcessed.endsWith('default')
          ) {
            resourceNameArray.push('default');
          }

          if (resourceNameArray.length > metricNamespaceArray.length) {
            const parentResource = resourceNameArray.shift();
            if (parentResource) {
              urlArray.push(parentResource);
            }
          }

          for (const i in metricNamespaceArray) {
            urlArray.push(metricNamespaceArray[i]);
            urlArray.push(resourceNameArray[i]);
          }
        }
      }
    }

    return urlArray.join('/');
  }

  static buildAzureMonitorGetMetricNamespacesUrl(
    baseUrl: string,
    apiVersion: string,
    query: GetMetricNamespacesQuery,
    templateSrv: TemplateSrv,
    globalRegion?: boolean,
    region?: string
  ) {
    let resourceUri: string;

    if ('resourceUri' in query) {
      resourceUri = query.resourceUri;
    } else {
      const { subscription, resourceGroup, metricNamespace, resourceName } = query;
      resourceUri = UrlBuilder.buildResourceUri(templateSrv, {
        subscription,
        resourceGroup, // if resourceGroup is not present use the original endpoint
        metricNamespace,
        resourceName,
      });
    }

    if (resourceUri.includes('resourceGroups')) {
      return `${baseUrl}${resourceUri}/resources?api-version=${apiVersion}${
        region ? `&region=${region}` : globalRegion ? '&region=global' : ''
      }`;
    } else {
      apiVersion = "2017-12-01-preview"
      return `${baseUrl}${resourceUri}/providers/microsoft.insights/metricNamespaces?api-version=${apiVersion}${
        region ? `&region=${region}` : globalRegion ? '&region=global' : ''
      }`;
    }
    // JUST MAKE RESOURCE GROUPS BEHAVIOR CHANGE NOT AT THE SUB OR RESOURCE LEVEL
    // toggle this on or off?
    // need distinct to filter out namespaces with metrics? resources endpoint supports filtering? use resourcegraph query
  }

  static buildAzureMonitorGetMetricNamesUrl(
    baseUrl: string,
    apiVersion: string,
    query: GetMetricNamesQuery,
    templateSrv: TemplateSrv,
    multipleResources?: boolean,
    region?: string
  ) {
    let resourceUri: string;
    const { customNamespace, metricNamespace } = query;
    if ('resourceUri' in query) {
      resourceUri = query.resourceUri;
    } else {
      const { subscription, resourceGroup, metricNamespace, resourceName } = query;
      resourceUri = UrlBuilder.buildResourceUri(
        templateSrv,
        {
          subscription,
          resourceGroup,
          metricNamespace,
          resourceName,
        },
        multipleResources
      );
    }
    let url = `${baseUrl}${resourceUri}/providers/microsoft.insights/metricdefinitions?api-version=${apiVersion}`;
    if (customNamespace) {
      url += `&metricnamespace=${encodeURIComponent(customNamespace)}`;
    }

    if (multipleResources && !customNamespace && metricNamespace) {
      url += `&metricnamespace=${encodeURIComponent(metricNamespace)}`;
    }

    if (region && multipleResources) {
      url += `&region=${region}`;
    }

    return url;
  }
}
