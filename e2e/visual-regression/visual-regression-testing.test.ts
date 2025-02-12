import { test, expect } from '@playwright/test';

test('BarGauge', { tag: '@viz_regression' }, async ({ page }) => {
  const barGaugeLocation = 'd/O6f11TZWk/panel-tests-bar-gauge?orgId=1&from=now-6h&to=now&timezone=Europe%2FLondon';
  await page.goto(`/${barGaugeLocation}`);

  const ledBarGauge = await page.locator('.react-grid-item').nth(6).boundingBox();

  await expect(page).toHaveScreenshot('bargauge.png', { animations: 'disabled', clip: ledBarGauge, fullPage: true });
});

test('TimeSeriesStacking', { tag: '@viz_regression' }, async ({ page }) => {
  const timeSeriesStackingLocation =
    'd/Ei86FP_Mx/panel-tests-timeseries-stacking?orgId=1&from=now-6h&to=now&timezone=Europe%2FLondon';
  await page.goto(`/${timeSeriesStackingLocation}`);

  const stackedTimeSeries = await page.locator('.react-grid-item').first().boundingBox();
  const yAxis = page.locator('.u-axis').nth(0);
  const xAxis = page.locator('.u-axis').nth(1);

  await expect(page).toHaveScreenshot('timeseriesstacking.png', {
    animations: 'disabled',
    clip: stackedTimeSeries,
    fullPage: true,
    mask: [yAxis, xAxis],
  });
});
