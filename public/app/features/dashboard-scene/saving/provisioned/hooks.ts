import { Chance } from 'chance';

import { dateTime } from '@grafana/data';
import { AnnoKeyRepoName, AnnoKeyRepoPath } from 'app/features/apiserver/types';
import { useGetResourceRepository } from 'app/features/provisioning/hooks';
import { DashboardMeta } from 'app/types';

import { getDefaultWorkflow } from './defaults';

interface UseDefaultValuesParams {
  meta: DashboardMeta;
  defaultTitle: string;
  defaultDescription?: string;
}

function generatePath(timestamp: string, pathFromAnnotation?: string, slug?: string) {
  if (pathFromAnnotation) {
    const hashIndex = pathFromAnnotation.indexOf('#');
    return hashIndex > 0 ? pathFromAnnotation.substring(0, hashIndex) : pathFromAnnotation;
  }

  const pathSlug = slug || `new-dashboard-${timestamp}`;
  return `${pathSlug}.json`;
}

export function useDefaultValues({ meta, defaultTitle, defaultDescription }: UseDefaultValuesParams) {
  const annotations = meta.k8s?.annotations;
  const annoName = annotations?.[AnnoKeyRepoName];
  const annoPath = annotations?.[AnnoKeyRepoPath];
  // Get config by resource name or folder UID for new resources
  const repositoryConfig = useGetResourceRepository({ name: annoName, folderUid: meta.folderUid });
  const repository = repositoryConfig?.spec;

  const random = Chance(1);
  const timestamp = `${dateTime().format('YYYY-MM-DD')}-${random.string({ length: 5, alpha: true })}`;

  return {
    values: {
      ref: `dashboard/${timestamp}`,
      path: generatePath(timestamp, annoPath, meta.slug),
      repo: annoName || repositoryConfig?.metadata?.name || '',
      comment: '',
      folder: {
        uid: meta.folderUid,
        title: '',
      },
      title: defaultTitle,
      description: defaultDescription ?? '',
      workflow: getDefaultWorkflow(repository),
    },
    isNew: !annoName,
    repositoryConfig: repository,
    isGitHub: repository?.type === 'github',
  };
}
