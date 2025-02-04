import { css } from '@emotion/css';
import Skeleton from 'react-loading-skeleton';

import type { GrafanaTheme2 } from '@grafana/data';
import { getInputStyles, useStyles2 } from '@grafana/ui';

// Aim to keep the bundle-size of the as small as possible so it doesn't
// impact tree shaking too much
export function FolderPickerSkeleton() {
  const styles = useStyles2(getStyles);

  return (
    <div className={styles.wrapper}>
      <div className={styles.inputWrapper}>
        <button type="button" className={styles.fakeInput} aria-disabled>
          <Skeleton width={100} />
        </button>
      </div>
    </div>
  );
}

const getStyles = (theme: GrafanaTheme2) => {
  const baseStyles = getInputStyles({ theme });

  return {
    wrapper: baseStyles.wrapper,
    inputWrapper: baseStyles.inputWrapper,
    fakeInput: css([
      baseStyles.input,
      {
        textAlign: 'left',
      },
    ]),
  };
};
