import { cloneDeep } from 'lodash';
import { useCallback, useEffect, useMemo, useReducer, useState } from 'react';
import { Controller, useFormContext } from 'react-hook-form';
import { useEffectOnce } from 'react-use';

import { getDefaultRelativeTimeRange } from '@grafana/data';
import { selectors } from '@grafana/e2e-selectors';
import { config, getDataSourceSrv } from '@grafana/runtime';
import { Alert, Button, ConfirmModal, Dropdown, Field, Icon, Menu, MenuItem, Stack, Tooltip } from '@grafana/ui';
import { Text } from '@grafana/ui/src/components/Text/Text';
import { Trans, t } from 'app/core/internationalization';
import { isExpressionQuery } from 'app/features/expressions/guards';
import { ExpressionDatasourceUID, ExpressionQueryType, expressionTypes } from 'app/features/expressions/types';
import { AlertQuery } from 'app/types/unified-alerting-dto';

import { useRulesSourcesWithRuler } from '../../../hooks/useRuleSourcesWithRuler';
import {
  areQueriesTransformableToSimpleCondition,
  isExpressionQueryInAlert,
} from '../../../rule-editor/formProcessing';
import { RuleFormType, RuleFormValues } from '../../../types/rule-form';
import { getDefaultOrFirstCompatibleDataSource } from '../../../utils/datasource';
import { PromOrLokiQuery, isPromOrLokiQuery } from '../../../utils/rule-form';
import {
  isAlertingRuleByType,
  isDataSourceManagedRuleByType,
  isGrafanaManagedRuleByType,
  isRecordingRuleByType,
} from '../../../utils/rules';
import { ExpressionEditor } from '../ExpressionEditor';
import { ExpressionsEditor } from '../ExpressionsEditor';
import { QueryEditor } from '../QueryEditor';
import { RecordingRuleEditor } from '../RecordingRuleEditor';
import { RuleEditorSection, RuleEditorSubSection } from '../RuleEditorSection';
import { errorFromCurrentCondition, errorFromPreviewData, findRenamedDataQueryReferences, refIdExists } from '../util';

import { CloudDataSourceSelector } from './CloudDataSourceSelector';
import { SimpleConditionEditor, SimpleConditionIdentifier, getSimpleConditionFromExpressions } from './SimpleCondition';
import { SmartAlertTypeDetector } from './SmartAlertTypeDetector';
import { DESCRIPTIONS } from './descriptions';
import {
  addExpressions,
  addNewDataQuery,
  addNewExpression,
  duplicateQuery,
  optimizeReduceExpression,
  queriesAndExpressionsReducer,
  removeExpression,
  removeExpressions,
  resetToSimpleCondition,
  rewireExpressions,
  setDataQueries,
  setRecordingRulesQueries,
  updateExpression,
  updateExpressionRefId,
  updateExpressionTimeRange,
  updateExpressionType,
} from './reducer';
import { useAdvancedMode } from './useAdvancedMode';
import { useAlertQueryRunner } from './useAlertQueryRunner';

interface Props {
  editingExistingRule: boolean;
  onDataChange: (error: string) => void;
}

export const QueryAndExpressionsStep = ({ editingExistingRule, onDataChange }: Props) => {
  const {
    setValue,
    getValues,
    watch,
    formState: { errors },
    control,
  } = useFormContext<RuleFormValues>();

  const { queryPreviewData, runQueries, cancelQueries, isPreviewLoading, clearPreviewData } = useAlertQueryRunner();

  const initialState = {
    queries: getValues('queries'),
  };

  const [{ queries }, dispatch] = useReducer(queriesAndExpressionsReducer, initialState);
  const isOptimizeReducerEnabled = config.featureToggles.alertingUIOptimizeReducer ?? false;

  // data queries only
  const dataQueries = useMemo(() => {
    return queries.filter((query) => !isExpressionQuery(query.model));
  }, [queries]);

  // expression queries only
  const expressionQueries = useMemo(() => {
    return queries.filter((query) => isExpressionQueryInAlert(query));
  }, [queries]);

  useEffectOnce(() => {
    // we only remove or add the reducer(optimize reducer) expression when creating a new alert.
    // When editing an alert, we assume the user wants to manually adjust expressions and queries for more control and customization.

    if (!editingExistingRule && isOptimizeReducerEnabled) {
      dispatch(optimizeReduceExpression({ updatedQueries: dataQueries, expressionQueries }));
    }
  });

  const [type, condition, dataSourceName, editorSettings] = watch([
    'type',
    'condition',
    'dataSourceName',
    'editorSettings',
  ]);

  // differentiate Grafana-managed with Datasource-managed
  const isGrafanaManagedRule = isGrafanaManagedRuleByType(type);
  const isDataSourceManagedRule = isDataSourceManagedRuleByType(type);

  // differentiate alerting rule from recording rule
  const isRecordingRuleType = isRecordingRuleByType(type);
  const isAlertingRuleType = isAlertingRuleByType(type);

  // Grafana managed alert rule
  const isGrafanaAlertingType = isAlertingRuleType && isGrafanaManagedRule;
  const isDataSourceAlertingType = isAlertingRuleType && isDataSourceManagedRule;

  const [showResetModeModal, setShowResetModal] = useState(false);

  const simplifiedQueryInForm = editorSettings?.simplifiedQueryEditor;

  const { simpleCondition, setSimpleCondition, simplifiedQueryStep } = useAdvancedMode(
    simplifiedQueryInForm,
    isGrafanaAlertingType,
    dataQueries,
    expressionQueries
  );

  // If we switch to simple mode we need to update the simple condition with the data in the queries reducer
  useEffect(() => {
    if (simplifiedQueryStep && isGrafanaAlertingType) {
      setSimpleCondition(getSimpleConditionFromExpressions(expressionQueries));
    }
  }, [simplifiedQueryStep, expressionQueries, isGrafanaAlertingType, setSimpleCondition]);

  const { rulesSourcesWithRuler } = useRulesSourcesWithRuler();

  const q = watch('queries');

  const runQueriesPreview = useCallback(
    (condition?: string) => {
      if (isDataSourceManagedRule) {
        // we will skip preview for Data source managed rules, these do not have any time series preview
        // Grafana Managed rules and recording rules do
        return;
      }
      // we need to be sure the condition is set once we switch to simple mode
      if (simplifiedQueryStep) {
        setValue('condition', SimpleConditionIdentifier.thresholdId);
        runQueries(q, SimpleConditionIdentifier.thresholdId);
      } else {
        runQueries(q, condition || (getValues('condition') ?? ''));
      }
    },
    [isDataSourceManagedRule, simplifiedQueryStep, setValue, runQueries, q, getValues]
  );

  // whenever we update the queries we have to update the form too
  useEffect(() => {
    setValue('queries', queries, { shouldValidate: false });
  }, [queries, runQueries, setValue]);

  const noCompatibleDataSources = getDefaultOrFirstCompatibleDataSource() === undefined;

  const emptyQueries = queries.length === 0;

  // apply some validations and asserts to the results of the evaluation when creating or editing
  // Grafana-managed alert rules and Grafa-managed recording rules
  useEffect(() => {
    if (type && !isGrafanaManagedRuleByType(type)) {
      return;
    }

    const currentCondition = getValues('condition');
    if (!currentCondition) {
      return;
    }

    const previewData = queryPreviewData[currentCondition];
    if (!previewData) {
      return;
    }

    const error = errorFromPreviewData(previewData) ?? errorFromCurrentCondition(previewData);

    onDataChange(error?.message || '');
  }, [queryPreviewData, getValues, onDataChange, type]);

  // run preview once when everything for the rule definition is loaded
  useEffect(() => {
    runQueriesPreview();
  }, [runQueriesPreview]);

  const handleSetCondition = useCallback(
    (refId: string | null) => {
      if (!refId) {
        return;
      }

      runQueriesPreview(refId); //we need to run the queries to know if the condition is valid

      setValue('condition', refId);
    },
    [runQueriesPreview, setValue]
  );

  const onUpdateRefId = useCallback(
    (oldRefId: string, newRefId: string) => {
      const newRefIdExists = refIdExists(queries, newRefId);
      // TODO we should set an error and explain what went wrong instead of just refusing to update
      if (newRefIdExists) {
        return;
      }

      dispatch(updateExpressionRefId({ oldRefId, newRefId }));

      // update condition too if refId was updated
      if (condition === oldRefId) {
        setValue('condition', newRefId);
      }
    },
    [condition, queries, setValue]
  );

  const updateExpressionAndDatasource = useSetExpressionAndDataSource();

  const onChangeQueries = useCallback(
    (updatedQueries: AlertQuery[]) => {
      // Most data sources triggers onChange and onRunQueries consecutively
      // It means our reducer state is always one step behind when runQueries is invoked
      // Invocation cycle => onChange -> dispatch(setDataQueries) -> onRunQueries -> setDataQueries Reducer
      // As a workaround we update form values as soon as possible to avoid stale state
      // This way we can access up to date queries in runQueriesPreview without waiting for re-render
      const previousQueries = getValues('queries');
      const expressionQueries = previousQueries.filter((query) => isExpressionQuery(query.model));
      setValue('queries', [...updatedQueries, ...expressionQueries], { shouldValidate: false });
      updateExpressionAndDatasource(updatedQueries);

      // we only remove or add the reducer(optimize reducer) expression when creating a new alert.
      // When editing an alert, we assume the user wants to manually adjust expressions and queries for more control and customization.
      if (!editingExistingRule && isOptimizeReducerEnabled) {
        dispatch(optimizeReduceExpression({ updatedQueries, expressionQueries }));
      }

      dispatch(setDataQueries(updatedQueries));
      dispatch(updateExpressionTimeRange());

      // check if we need to rewire expressions (and which ones)
      const [oldRefId, newRefId] = findRenamedDataQueryReferences(queries, updatedQueries);
      if (oldRefId && newRefId) {
        dispatch(rewireExpressions({ oldRefId, newRefId }));
      }
    },
    [queries, updateExpressionAndDatasource, getValues, setValue, editingExistingRule, isOptimizeReducerEnabled]
  );

  const onChangeRecordingRulesQueries = useCallback(
    (updatedQueries: AlertQuery[]) => {
      const query = updatedQueries[0];

      if (!isPromOrLokiQuery(query.model)) {
        return;
      }

      const expression = query.model.expr;

      setValue('queries', updatedQueries, { shouldValidate: false });
      updateExpressionAndDatasource(updatedQueries);

      dispatch(setRecordingRulesQueries({ recordingRuleQueries: updatedQueries, expression }));
      runQueriesPreview();
    },
    [runQueriesPreview, setValue, updateExpressionAndDatasource]
  );

  // Using dataSourcesWithRuler[0] gives incorrect types - no undefined
  // Using at(0) provides a safe type with undefined
  const recordingRuleDefaultDatasource = rulesSourcesWithRuler.at(0);

  useEffect(() => {
    clearPreviewData();
    if (type === RuleFormType.cloudRecording) {
      const expr = getValues('expression');

      if (!recordingRuleDefaultDatasource) {
        return;
      }

      const datasourceUid =
        (editingExistingRule && getDataSourceSrv().getInstanceSettings(dataSourceName)?.uid) ||
        recordingRuleDefaultDatasource.uid;

      const defaultQuery = {
        refId: 'A',
        datasourceUid,
        queryType: '',
        relativeTimeRange: getDefaultRelativeTimeRange(),
        expr,
        instant: true,
        model: {
          refId: 'A',
          hide: false,
          expr,
        },
      };
      dispatch(setRecordingRulesQueries({ recordingRuleQueries: [defaultQuery], expression: expr }));
    }
  }, [type, recordingRuleDefaultDatasource, editingExistingRule, getValues, dataSourceName, clearPreviewData]);

  const onDuplicateQuery = useCallback((query: AlertQuery) => {
    dispatch(duplicateQuery(query));
  }, []);

  // update the condition if it's been removed
  useEffect(() => {
    if (!refIdExists(queries, condition)) {
      const lastRefId = queries.at(-1)?.refId ?? null;
      handleSetCondition(lastRefId);
    }
  }, [condition, queries, handleSetCondition]);

  const onClickType = useCallback(
    (type: ExpressionQueryType) => {
      dispatch(addNewExpression(type));
    },
    [dispatch]
  );

  // Cloud alerts load data from form values
  // whereas Grafana managed alerts load data from reducer
  //when data source is changed in the cloud selector we need to update the queries in the reducer

  const onChangeCloudDatasource = useCallback(
    (datasourceUid: string) => {
      const newQueries = cloneDeep(queries);
      newQueries[0].datasourceUid = datasourceUid;
      setValue('queries', newQueries, { shouldValidate: false });

      updateExpressionAndDatasource(newQueries);

      dispatch(setDataQueries(newQueries));
    },
    [queries, setValue, updateExpressionAndDatasource, dispatch]
  );

  // ExpressionEditor for cloud query needs to update queries in the reducer and in the form
  // otherwise the value is not updated for Grafana managed alerts

  const onChangeExpression = (value: string) => {
    const newQueries = cloneDeep(queries);

    if (newQueries[0].model) {
      if (isPromOrLokiQuery(newQueries[0].model)) {
        newQueries[0].model.expr = value;
      } else {
        // first time we come from grafana-managed type
        // we need to convert the model to PromOrLokiQuery
        const promLoki: PromOrLokiQuery = {
          ...cloneDeep(newQueries[0].model),
          expr: value,
        };
        newQueries[0].model = promLoki;
      }
    }

    setValue('queries', newQueries, { shouldValidate: false });

    updateExpressionAndDatasource(newQueries);

    dispatch(setDataQueries(newQueries));
    runQueriesPreview();
  };

  const removeExpressionsInQueries = useCallback(() => dispatch(removeExpressions()), [dispatch]);

  const addExpressionsInQueries = useCallback(
    (expressions: AlertQuery[]) => dispatch(addExpressions(expressions)),
    [dispatch]
  );

  // we need to keep track of the previous expressions and condition reference to be able to restore them when switching back to grafana managed
  const [prevExpressions, setPrevExpressions] = useState<AlertQuery[]>([]);
  const [prevCondition, setPrevCondition] = useState<string | null>(null);

  const restoreExpressionsInQueries = useCallback(() => {
    addExpressionsInQueries(prevExpressions);
  }, [prevExpressions, addExpressionsInQueries]);

  const onClickSwitch = useCallback(() => {
    const typeInForm = getValues('type');
    if (typeInForm === RuleFormType.cloudAlerting) {
      setValue('type', RuleFormType.grafana);
      setValue('dataSourceName', null); // set data source name back to "null"

      prevExpressions.length > 0 && restoreExpressionsInQueries();
      prevCondition && setValue('condition', prevCondition);
    } else {
      setValue('type', RuleFormType.cloudAlerting);
      // dataSourceName is used only by Mimir/Loki alerting and recording rules
      // It should be empty for Grafana managed alert rules
      const newDsName = getDataSourceSrv().getInstanceSettings(queries[0].datasourceUid)?.name;
      if (newDsName) {
        setValue('dataSourceName', newDsName);
      }

      updateExpressionAndDatasource(queries);

      const expressions = queries.filter((query) => query.datasourceUid === ExpressionDatasourceUID);
      setPrevExpressions(expressions);
      removeExpressionsInQueries();
      setPrevCondition(condition);
    }
  }, [
    getValues,
    setValue,
    prevExpressions.length,
    restoreExpressionsInQueries,
    prevCondition,
    updateExpressionAndDatasource,
    queries,
    removeExpressionsInQueries,
    condition,
  ]);

  const { sectionTitle, sectionDescription, helpContent, helpLink } = DESCRIPTIONS[type ?? RuleFormType.grafana];

  if (!type) {
    return null;
  }
  const switchMode = isGrafanaAlertingType
    ? {
        isAdvancedMode: !simplifiedQueryStep,
        setAdvancedMode: (isAdvanced: boolean) => {
          if (!getValues('editorSettings.simplifiedQueryEditor')) {
            if (!areQueriesTransformableToSimpleCondition(dataQueries, expressionQueries)) {
              setShowResetModal(true);
              return;
            }
          }
          setValue('editorSettings.simplifiedQueryEditor', !isAdvanced);
        },
      }
    : undefined;

  const expressionError = errors.expression?.message;

  return (
    <>
      <RuleEditorSection
        stepNo={2}
        title={sectionTitle}
        description={sectionDescription}
        helpInfo={{
          title: sectionDescription,
          contentText: helpContent,
          externalLink: helpLink,
          linkText: 'Read more on our documentation website',
        }}
        switchMode={switchMode}
      >
        <RuleEditorSubSection fullWidth>
          {/* This is the cloud data source selector */}
          {isDataSourceManagedRule && (
            <CloudDataSourceSelector onChangeCloudDatasource={onChangeCloudDatasource} disabled={editingExistingRule} />
          )}

          {/* This is the PromQL Editor for recording rules */}
          {isRecordingRuleType && isDataSourceManagedRule && dataSourceName && (
            <Field error={expressionError} invalid={Boolean(expressionError)} style={{ marginBottom: 0 }}>
              <RecordingRuleEditor
                dataSourceName={dataSourceName}
                queries={queries}
                runQueries={() => runQueriesPreview()}
                onChangeQuery={onChangeRecordingRulesQueries}
                panelData={queryPreviewData}
              />
            </Field>
          )}

          {/* This is the PromQL Editor for Cloud rules */}
          {/* @TODO why is dataSourceName "grafana" when creating a data source managed recording rule? */}
          {isDataSourceManagedRule && dataSourceName && dataSourceName !== 'grafana' && (
            <Field error={expressionError} invalid={Boolean(expressionError)} style={{ marginBottom: 0 }}>
              <Controller
                name="expression"
                render={({ field: { ref, ...field } }) => {
                  return (
                    <ExpressionEditor
                      {...field}
                      dataSourceName={dataSourceName}
                      showPreview={!isRecordingRuleType}
                      onChange={onChangeExpression}
                    />
                  );
                }}
                control={control}
                rules={{
                  required: { value: true, message: 'A valid expression is required' },
                }}
              />
            </Field>
          )}

          {/* This is the editor for Grafana managed rules and Grafana managed recording rules */}
          {isGrafanaManagedRule && (
            <Stack direction="column" alignItems="flex-start">
              {/* Data Queries */}
              <QueryEditor
                queries={dataQueries}
                expressions={expressionQueries}
                onRunQueries={() => runQueriesPreview()}
                onChangeQueries={onChangeQueries}
                onDuplicateQuery={onDuplicateQuery}
                panelData={queryPreviewData}
                condition={condition}
                onSetCondition={handleSetCondition}
              />
              {!simplifiedQueryStep && (
                <Tooltip content={'You appear to have no compatible data sources'} show={noCompatibleDataSources}>
                  <Button
                    type="button"
                    onClick={() => {
                      dispatch(addNewDataQuery());
                    }}
                    variant="secondary"
                    data-testid={selectors.components.QueryTab.addQuery}
                    disabled={noCompatibleDataSources}
                  >
                    Add query
                  </Button>
                </Tooltip>
              )}
            </Stack>
          )}
        </RuleEditorSubSection>

        {/* This is for switching from Grafana-managed to data-source managed */}
        {((isGrafanaAlertingType && !simplifiedQueryStep) || isDataSourceAlertingType) && (
          <RuleEditorSubSection
            title="Rule type"
            description="Select where the alert rule will be managed."
            helpInfo={{
              contentText: (
                <>
                  <Text color="primary" variant="h6">
                    Grafana-managed alert rules
                  </Text>
                  <p>
                    Grafana-managed alert rules allow you to create alerts that can act on data from any of our
                    supported data sources, including having multiple data sources in the same rule. You can also add
                    expressions to transform your data and set alert conditions. Using images in alert notifications is
                    also supported.
                  </p>
                  <Text color="primary" variant="h6">
                    Data source-managed alert rules
                  </Text>
                  <p>
                    Data source-managed alert rules can be used for Grafana Mimir or Grafana Loki data sources which
                    have been configured to support rule creation. The use of expressions or multiple queries is not
                    supported.
                  </p>
                </>
              ),
              externalLink:
                'https://grafana.com/docs/grafana/latest/alerting/fundamentals/alert-rules/alert-rule-types/',
              linkText: 'Read about alert rule types',
              title: 'Alert rule types',
            }}
          >
            <SmartAlertTypeDetector
              editingExistingRule={editingExistingRule}
              rulesSourcesWithRuler={rulesSourcesWithRuler}
              queries={queries}
              onClickSwitch={onClickSwitch}
            />
          </RuleEditorSubSection>
        )}

        {/* Expression Queries */}
        {!simplifiedQueryStep && isGrafanaManagedRule && (
          <RuleEditorSubSection
            title="Expressions"
            description="Manipulate data returned from queries with math and other operations."
          >
            <ExpressionsEditor
              queries={queries}
              panelData={queryPreviewData}
              condition={condition}
              onSetCondition={handleSetCondition}
              onRemoveExpression={(refId) => {
                dispatch(removeExpression(refId));
              }}
              onUpdateRefId={onUpdateRefId}
              onUpdateExpressionType={(refId, type) => {
                dispatch(updateExpressionType({ refId, type }));
              }}
              onUpdateQueryExpression={(model) => {
                dispatch(updateExpression(model));
              }}
            />

            {/* action buttons */}
            <Stack direction="column">
              {simplifiedQueryStep && (
                <SimpleConditionEditor
                  simpleCondition={simpleCondition}
                  onChange={setSimpleCondition}
                  expressionQueriesList={expressionQueries}
                  dispatch={dispatch}
                  previewData={queryPreviewData[condition ?? '']}
                />
              )}
              <Stack direction="row">
                {!simplifiedQueryStep && config.expressionsEnabled && <TypeSelectorButton onClickType={onClickType} />}

                {isPreviewLoading && (
                  <Button icon="spinner" type="button" variant="destructive" onClick={cancelQueries}>
                    <Trans i18nKey="alerting.common.cancel">Cancel</Trans>
                  </Button>
                )}
                {!isPreviewLoading && (
                  <Button
                    data-testid={selectors.components.AlertRules.previewButton}
                    icon="sync"
                    type="button"
                    onClick={() => runQueriesPreview()}
                    disabled={emptyQueries}
                  >
                    {!simplifiedQueryStep
                      ? t('alerting.queryAndExpressionsStep.preview', 'Preview')
                      : t('alerting.queryAndExpressionsStep.previewCondition', 'Preview alert rule condition')}
                  </Button>
                )}
              </Stack>
            </Stack>

            {/* No Queries */}
            {emptyQueries && (
              <Alert title="No queries or expressions have been configured" severity="warning">
                Create at least one query or expression to be alerted on
              </Alert>
            )}
          </RuleEditorSubSection>
        )}
      </RuleEditorSection>

      <ConfirmModal
        isOpen={showResetModeModal}
        title="Deactivate advanced options"
        body={
          <div>
            <Text element="p">
              <Trans i18nKey="alerting.queryAndExpressionsStep.disableAdvancedOptions.text">
                The selected queries and expressions cannot be converted to default. If you deactivate advanced options,
                your query and condition will be reset to default settings.
              </Trans>
            </Text>
            <br />
          </div>
        }
        confirmText="Deactivate"
        icon="exclamation-triangle"
        onConfirm={() => {
          setValue('editorSettings.simplifiedQueryEditor', true);
          setShowResetModal(false);
          dispatch(resetToSimpleCondition());
        }}
        onDismiss={() => setShowResetModal(false)}
      />
    </>
  );
};

function TypeSelectorButton({ onClickType }: { onClickType: (type: ExpressionQueryType) => void }) {
  const newMenu = (
    <Menu>
      {expressionTypes.map((type) => (
        <Tooltip key={type.value} content={type.description ?? ''} placement="right">
          <MenuItem
            key={type.value}
            onClick={() => onClickType(type.value ?? ExpressionQueryType.math)}
            label={type.label ?? ''}
          />
        </Tooltip>
      ))}
    </Menu>
  );

  return (
    <Dropdown overlay={newMenu}>
      <Button variant="secondary" data-testid={'add-expression-button'}>
        Add expression
        <Icon name="angle-down" />
      </Button>
    </Dropdown>
  );
}

const useSetExpressionAndDataSource = () => {
  const { setValue } = useFormContext<RuleFormValues>();

  return (updatedQueries: AlertQuery[]) => {
    // update data source name and expression if it's been changed in the queries from the reducer when prom or loki query
    const query = updatedQueries[0];
    if (!query) {
      return;
    }

    const dataSourceSettings = getDataSourceSrv().getInstanceSettings(query.datasourceUid);
    if (!dataSourceSettings) {
      throw new Error('The Data source has not been defined.');
    }

    if (isPromOrLokiQuery(query.model)) {
      const expression = query.model.expr;
      setValue('expression', expression);
    }
  };
};
