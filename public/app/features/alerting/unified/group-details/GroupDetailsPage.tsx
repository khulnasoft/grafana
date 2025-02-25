import { skipToken } from '@reduxjs/toolkit/query';
import { useMemo } from 'react';
import { useParams } from 'react-router-dom-v5-compat';

import { Alert, Badge, LinkButton, Text, withErrorBoundary } from '@grafana/ui';
import { EntityNotFound } from 'app/core/components/PageNotFound/EntityNotFound';
import { t } from 'app/core/internationalization';
import { FolderDTO } from 'app/types';
import { GrafanaRulesSourceSymbol, RuleGroup } from 'app/types/unified-alerting';
import { PromRuleType, RulerRuleGroupDTO } from 'app/types/unified-alerting-dto';

import { alertRuleApi } from '../api/alertRuleApi';
import { RulesSourceFeatures, featureDiscoveryApi } from '../api/featureDiscoveryApi';
import { AlertingPageWrapper } from '../components/AlertingPageWrapper';
import { DynamicTable, DynamicTableColumnProps } from '../components/DynamicTable';
import { useFolder } from '../hooks/useFolder';
import { DEFAULT_GROUP_EVALUATION_INTERVAL } from '../rule-editor/formDefaults';
import { useRulesAccess } from '../utils/accessControlHooks';
import { GRAFANA_RULES_SOURCE_NAME } from '../utils/datasource';
import { stringifyErrorLike } from '../utils/misc';
import { groups } from '../utils/navigation';
import { getEvaluationsToStartAlerting, isAlertingRulerRule, isGrafanaAlertingRule } from '../utils/rules';
import { formatPrometheusDuration, safeParsePrometheusDuration } from '../utils/time';

type GroupPageRouteParams = {
  sourceId?: string;
  namespaceId?: string;
  groupName?: string;
};

const { useDiscoverDsFeaturesQuery } = featureDiscoveryApi;
const { usePrometheusRuleNamespacesQuery, useGetRuleGroupForNamespaceQuery } = alertRuleApi;

function GroupDetailsPage() {
  const { sourceId = '', namespaceId = '', groupName = '' } = useParams<GroupPageRouteParams>();

  const { folder, loading: isFolderLoading } = useFolder(sourceId === 'grafana' ? namespaceId : '');
  const {
    data: dsFeatures,
    isLoading: isDsFeaturesLoading,
    error: dsFeaturesError,
  } = useDiscoverDsFeaturesQuery({ uid: sourceId === 'grafana' ? GrafanaRulesSourceSymbol : sourceId });

  const {
    data: promGroup,
    isLoading: isRuleNamespacesLoading,
    error: ruleNamespacesError,
  } = usePrometheusRuleNamespacesQuery(
    !dsFeatures?.rulerConfig
      ? { ruleSourceName: dsFeatures?.name ?? '', namespace: namespaceId, groupName: groupName }
      : skipToken,
    {
      selectFromResult: (result) => ({
        ...result,
        data: result.data?.[0]?.groups.find((g) => g.name === groupName),
      }),
    }
  );

  const {
    data: rulerGroup,
    isLoading: isRuleGroupLoading,
    error: ruleGroupError,
  } = useGetRuleGroupForNamespaceQuery(
    dsFeatures?.rulerConfig
      ? { rulerConfig: dsFeatures?.rulerConfig, namespace: namespaceId, group: groupName }
      : skipToken
  );

  const isLoading = isFolderLoading || isDsFeaturesLoading || isRuleNamespacesLoading || isRuleGroupLoading;

  const groupInterval = promGroup?.interval
    ? formatPrometheusDuration(promGroup.interval * 1000)
    : (rulerGroup?.interval ?? DEFAULT_GROUP_EVALUATION_INTERVAL);

  const namespaceName = folder?.title ?? namespaceId;

  return (
    <AlertingPageWrapper
      pageNav={{ text: groupName }}
      title={groupName}
      info={[
        { label: 'Namespace', value: namespaceName },
        { label: 'Interval', value: groupInterval },
      ]}
      navId="alert-list"
      isLoading={isLoading}
      actions={
        <>
          {dsFeatures && (
            <GroupActions dsFeatures={dsFeatures} namespaceId={namespaceId} groupName={groupName} folder={folder} />
          )}
        </>
      }
    >
      <>
        {Boolean(dsFeaturesError) && (
          <Alert
            title={t('alerting.group-details.ds-features-error', 'Error loading data source details')}
            bottomSpacing={0}
            topSpacing={2}
          >
            <div>{stringifyErrorLike(dsFeaturesError)}</div>
          </Alert>
        )}
        {Boolean(ruleNamespacesError || ruleGroupError) && (
          <Alert
            title={t('alerting.group-details.group-loading-error', 'Error loading the group')}
            bottomSpacing={0}
            topSpacing={2}
          >
            <div>{stringifyErrorLike(ruleNamespacesError || ruleGroupError)}</div>
          </Alert>
        )}
        {promGroup && <GroupDetails group={promRuleGroupToRuleGroupDetails(promGroup)} />}
        {rulerGroup && <GroupDetails group={rulerRuleGroupToRuleGroupDetails(rulerGroup)} />}
        {(!promGroup || !rulerGroup) && <EntityNotFound entity={`${namespaceId}/${groupName}`} />}
      </>
    </AlertingPageWrapper>
  );
}

interface GroupActionsProps {
  dsFeatures: RulesSourceFeatures;
  namespaceId: string;
  groupName: string;
  folder: FolderDTO | undefined;
}

function GroupActions({ dsFeatures, namespaceId, groupName, folder }: GroupActionsProps) {
  const { canEditRules } = useRulesAccess();

  const isGrafanaSource = dsFeatures.uid === GRAFANA_RULES_SOURCE_NAME;
  const canSaveInFolder = isGrafanaSource ? !!folder?.canSave : true;
  const canEdit = Boolean(dsFeatures.rulerConfig) && canEditRules(dsFeatures.name) && canSaveInFolder;

  if (!canEdit) {
    return null;
  }

  return (
    <LinkButton icon="pen" href={groups.editPageLink(dsFeatures.uid, namespaceId, groupName)} variant="secondary">
      Edit
    </LinkButton>
  );
}

/** An common interface for both Prometheus and Ruler rule groups */
interface RuleGroupDetails {
  name: string;
  interval: string;
  rules: RuleDetails[];
}

interface AlertingRuleDetails {
  name: string;
  type: 'alerting';
  pendingPeriod: string;
  evaluationsToFire: number;
}
interface RecordingRuleDetails {
  name: string;
  type: 'recording';
}

type RuleDetails = AlertingRuleDetails | RecordingRuleDetails;

interface GroupDetailsProps {
  group: RuleGroupDetails;
}

function GroupDetails({ group }: GroupDetailsProps) {
  return (
    <div>
      <RulesTable rules={group.rules} />
    </div>
  );
}

function RulesTable({ rules }: { rules: RuleDetails[] }) {
  const rows = rules.map((rule: RuleDetails, index) => ({
    id: index,
    data: rule,
  }));

  const columns: Array<DynamicTableColumnProps<RuleDetails>> = useMemo(() => {
    return [
      {
        id: 'alertName',
        label: 'Rule name',
        renderCell: ({ data }) => {
          return <Text truncate>{data.name}</Text>;
        },
        size: 0.4,
      },
      {
        id: 'for',
        label: 'Pending period',
        renderCell: ({ data }) => {
          switch (data.type) {
            case 'alerting':
              return <>{data.pendingPeriod}</>;
            case 'recording':
              return <Badge text="Recording" color="purple" />;
          }
        },
        size: 0.3,
      },
      {
        id: 'numberEvaluations',
        label: 'Evaluation cycles to fire',
        renderCell: ({ data }) => {
          switch (data.type) {
            case 'alerting':
              return <>{data.evaluationsToFire}</>;
            case 'recording':
              return null;
          }
        },
        size: 0.3,
      },
    ];
  }, []);

  return <DynamicTable items={rows} cols={columns} />;
}

function promRuleGroupToRuleGroupDetails(group: RuleGroup): RuleGroupDetails {
  const groupIntervalMs = group.interval * 1000;

  return {
    name: group.name,
    interval: formatPrometheusDuration(group.interval * 1000),
    rules: group.rules.map<RuleDetails>((rule) => {
      switch (rule.type) {
        case PromRuleType.Alerting:
          return {
            name: rule.name,
            type: 'alerting',
            pendingPeriod: formatPrometheusDuration(rule.duration ? rule.duration * 1000 : 0),
            evaluationsToFire: getEvaluationsToStartAlerting(rule.duration ? rule.duration * 1000 : 0, groupIntervalMs),
          };
        case PromRuleType.Recording:
          return { name: rule.name, type: 'recording' };
      }
    }),
  };
}

function rulerRuleGroupToRuleGroupDetails(group: RulerRuleGroupDTO): RuleGroupDetails {
  const groupIntervalMs = safeParsePrometheusDuration(group.interval ?? DEFAULT_GROUP_EVALUATION_INTERVAL);

  return {
    name: group.name,
    interval: group.interval ?? DEFAULT_GROUP_EVALUATION_INTERVAL,
    rules: group.rules.map<RuleDetails>((rule) => {
      if (isAlertingRulerRule(rule)) {
        return {
          name: rule.alert,
          type: 'alerting',
          pendingPeriod: rule.for ?? '0s',
          evaluationsToFire: getEvaluationsToStartAlerting(
            rule.for ? safeParsePrometheusDuration(rule.for) : 0,
            groupIntervalMs
          ),
        };
      }
      if (isGrafanaAlertingRule(rule)) {
        return {
          name: rule.grafana_alert.title,
          type: 'alerting',
          pendingPeriod: rule.for ?? '0s',
          evaluationsToFire: getEvaluationsToStartAlerting(
            rule.for ? safeParsePrometheusDuration(rule.for) : 0,
            groupIntervalMs
          ),
        };
      }

      return { name: rule.record, type: 'recording' };
    }),
  };
}

export default withErrorBoundary(GroupDetailsPage, { style: 'page' });
