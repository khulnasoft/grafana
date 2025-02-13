import { css } from '@emotion/css';

import { GrafanaTheme2 } from '@grafana/data';

import { useStyles2 } from '../../themes';
import { InlineField } from '../Forms/InlineField';
import { Stack } from '../Layout/Stack/Stack';
import { InlineSwitch } from '../Switch/Switch';

import { HttpSettingsBaseProps } from './types';

const LABEL_WIDTH = 26;

export const HttpProxySettings = ({
  dataSourceConfig,
  onChange,
  showForwardOAuthIdentityOption = true,
}: HttpSettingsBaseProps) => {
  const gridLayout = useStyles2(getGridLayout);
  return (
    <div className={gridLayout}>
      <Stack direction="row" gap={0.5}>
        <InlineField label="TLS Client Auth" labelWidth={LABEL_WIDTH} disabled={dataSourceConfig.readOnly}>
          <InlineSwitch
            id="http-settings-tls-client-auth"
            value={dataSourceConfig.jsonData.tlsAuth || false}
            onChange={(event) => onChange({ ...dataSourceConfig.jsonData, tlsAuth: event!.currentTarget.checked })}
          />
        </InlineField>
        <InlineField
          label="With CA Cert"
          tooltip="Needed for verifying self-signed TLS Certs"
          labelWidth={LABEL_WIDTH}
          disabled={dataSourceConfig.readOnly}
        >
          <InlineSwitch
            id="http-settings-ca-cert"
            value={dataSourceConfig.jsonData.tlsAuthWithCACert || false}
            onChange={(event) =>
              onChange({ ...dataSourceConfig.jsonData, tlsAuthWithCACert: event!.currentTarget.checked })
            }
          />
        </InlineField>
      </Stack>
      <InlineField label="Skip TLS Verify" labelWidth={LABEL_WIDTH} disabled={dataSourceConfig.readOnly}>
        <InlineSwitch
          id="http-settings-skip-tls-verify"
          value={dataSourceConfig.jsonData.tlsSkipVerify || false}
          onChange={(event) => onChange({ ...dataSourceConfig.jsonData, tlsSkipVerify: event!.currentTarget.checked })}
        />
      </InlineField>
      {showForwardOAuthIdentityOption && (
        <InlineField
          label="Forward OAuth Identity"
          tooltip="Forward the user's upstream OAuth identity to the data source (Their access token gets passed along)."
          labelWidth={LABEL_WIDTH}
          disabled={dataSourceConfig.readOnly}
        >
          <InlineSwitch
            id="http-settings-forward-oauth"
            value={dataSourceConfig.jsonData.oauthPassThru || false}
            onChange={(event) =>
              onChange({ ...dataSourceConfig.jsonData, oauthPassThru: event!.currentTarget.checked })
            }
          />
        </InlineField>
      )}
    </div>
  );
};

const getGridLayout = (theme: GrafanaTheme2) =>
  css({
    display: 'grid',
    gridTemplateColumns: 'auto',
    gap: 0, // Inline field has a margin
  });
