import { css } from '@emotion/css';
import { useMemo } from 'react';

import { GrafanaTheme2 } from '@grafana/data';
import { SceneObject } from '@grafana/scenes';
import { Stack, useStyles2 } from '@grafana/ui';
import { OptionsPaneCategory } from 'app/features/dashboard/components/PanelEditor/OptionsPaneCategory';

import { EditableDashboardElement, isEditableDashboardElement } from '../scene/types';

export interface Props {
  obj: SceneObject;
}

export function ElementEditPane({ obj }: Props) {
  const element = getEditableElementFor(obj);
  const categories = useMemo(() => element.getEditPaneOptions(), [element]);
  const styles = useStyles2(getStyles);

  return (
    <Stack direction="column" gap={0}>
      <OptionsPaneCategory id="selected-item" title={element.getTypeName()} isOpenDefault={true}>
        <div className={styles.actionsBox}>{element.renderActions()}</div>
      </OptionsPaneCategory>
      {categories.map((cat) => cat.render())}
    </Stack>
  );
}

function getEditableElementFor(obj: SceneObject): EditableDashboardElement {
  if (isEditableDashboardElement(obj)) {
    return obj;
  }

  for (const behavior of obj.state.$behaviors ?? []) {
    if (isEditableDashboardElement(behavior)) {
      return behavior;
    }
  }

  throw new Error("Can't find editable element for selected object");
}

function getStyles(theme: GrafanaTheme2) {
  return {
    actionsBox: css({
      display: 'flex',
      alignItems: 'center',
      gap: theme.spacing(1),
      paddingBottom: theme.spacing(1),
    }),
  };
}
