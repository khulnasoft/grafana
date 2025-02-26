import { HttpResponse } from 'msw';
import { Route, Routes } from 'react-router-dom-v5-compat';
import { render, screen } from 'test/test-utils';
import { byRole, byTestId } from 'testing-library-selector';

import { AccessControlAction } from 'app/types';

import { setupMswServer } from '../mockApi';
import { grantUserPermissions, mockRulerGrafanaRule, mockRulerRuleGroup } from '../mocks';
import {
  mimirDataSource,
  setFolderResponse,
  setPrometheusRules,
  setRulerRuleGroupHandler,
  setRulerRuleGroupResolver,
} from '../mocks/server/configure';
import { alertingFactory } from '../mocks/server/db';

import GroupDetailsPage from './GroupDetailsPage';

const ui = {
  header: byRole('heading', { level: 1 }),
  editLink: byRole('link', { name: 'Edit' }),
  tableRow: byTestId('row'),
  rowsTable: byTestId('dynamic-table'),
};

setupMswServer();

describe('GroupDetailsPage', () => {
  beforeEach(() => {
    grantUserPermissions([
      AccessControlAction.AlertingRuleRead,
      AccessControlAction.AlertingRuleUpdate,
      AccessControlAction.AlertingRuleExternalRead,
      AccessControlAction.AlertingRuleExternalWrite,
    ]);
  });

  describe('Grafana managed rules', () => {
    it('should render grafana rules group based on the Ruler API', async () => {
      const group = mockRulerRuleGroup({
        name: 'test-group-cpu',
        interval: '5m',
        rules: [mockRulerGrafanaRule(), mockRulerGrafanaRule()],
      });
      setRulerRuleGroupHandler({ response: HttpResponse.json(group) });
      setFolderResponse({ uid: 'test-folder-uid', canSave: true });

      // Act
      renderGroupDetailsPage('grafana', 'test-folder-uid', 'test-group-cpu');

      const header = await ui.header.find();
      const editLink = await ui.editLink.find();

      // Assert
      expect(header).toHaveTextContent('test-group-cpu');
      expect(await screen.findByText(/test-folder-uid/)).toBeInTheDocument();
      expect(await screen.findByText(/5m/)).toBeInTheDocument();
      expect(editLink).toHaveAttribute(
        'href',
        '/alerting/grafana/namespaces/test-folder-uid/groups/test-group-cpu/edit'
      );
    });

    it('should render error alert when API returns an error', async () => {
      // Mock an error response from the API
      setRulerRuleGroupResolver((req) => {
        return HttpResponse.json({ error: 'Failed to fetch rule group' }, { status: 500 });
      });

      // Act
      renderGroupDetailsPage('grafana', 'test-folder-uid', 'test-group-cpu');

      // Assert
      expect(await screen.findByText('Error loading the group')).toBeInTheDocument();
      expect(await screen.findByText('Failed to fetch rule group')).toBeInTheDocument();
    });

    it('should render "not found" when group does not exist', async () => {
      // Mock a 404 response
      setRulerRuleGroupResolver((req) => {
        return HttpResponse.json({ error: 'rule group does not exist' }, { status: 404 });
      });

      // Act
      renderGroupDetailsPage('grafana', 'test-folder-uid', 'non-existent-group');

      const notFoundAlert = await screen.findByRole('alert', { name: /Error loading the group/ });

      // Assert
      expect(notFoundAlert).toBeInTheDocument();
      expect(notFoundAlert).toHaveTextContent(/rule group does not exist/);
      expect(screen.getByTestId('data-testid entity-not-found')).toHaveTextContent(
        'test-folder-uid/non-existent-group'
      );
    });

    it('should not show edit button when user lacks edit permissions', async () => {
      // Remove edit permissions
      grantUserPermissions([AccessControlAction.AlertingRuleRead, AccessControlAction.AlertingRuleExternalRead]);

      const group = mockRulerRuleGroup({
        name: 'test-group-cpu',
        interval: '5m',
        rules: [mockRulerGrafanaRule()],
      });
      setRulerRuleGroupHandler({ response: HttpResponse.json(group) });
      setFolderResponse({ uid: 'test-folder-uid', canSave: true });

      // Act
      renderGroupDetailsPage('grafana', 'test-folder-uid', 'test-group-cpu');

      const tableRows = await ui.tableRow.findAll(await ui.rowsTable.find());

      // Assert
      expect(tableRows).toHaveLength(1);
      expect(ui.editLink.query()).not.toBeInTheDocument(); // Edit button should not be present
    });

    it('should not show edit button when folder cannot be saved', async () => {
      const group = mockRulerRuleGroup({
        name: 'test-group-cpu',
        interval: '5m',
        rules: [mockRulerGrafanaRule()],
      });
      setRulerRuleGroupHandler({ response: HttpResponse.json(group) });
      setFolderResponse({ uid: 'test-folder-uid', canSave: false }); // Folder cannot be saved

      // Act
      renderGroupDetailsPage('grafana', 'test-folder-uid', 'test-group-cpu');

      const tableRows = await ui.tableRow.findAll(await ui.rowsTable.find());

      // Assert
      expect(tableRows).toHaveLength(1);
      expect(ui.editLink.query()).not.toBeInTheDocument(); // Edit button should not be present
    });

    it('should render rules with correct details in the table', async () => {
      // Create rules with specific properties to test table rendering
      const rule1 = mockRulerGrafanaRule({ for: '10m' }, { title: 'High CPU Usage' });
      const rule2 = mockRulerGrafanaRule({ for: '5m' }, { title: 'Memory Pressure' });

      const group = mockRulerRuleGroup({
        name: 'test-group-resources',
        interval: '3m',
        rules: [rule1, rule2],
      });

      setRulerRuleGroupHandler({ response: HttpResponse.json(group) });
      setFolderResponse({ uid: 'test-folder-uid', canSave: true });

      // Act
      renderGroupDetailsPage('grafana', 'test-folder-uid', 'test-group-resources');

      const tableRows = await ui.tableRow.findAll(await ui.rowsTable.find());

      // Assert
      expect(tableRows).toHaveLength(2);

      expect(tableRows[0]).toHaveTextContent('High CPU Usage');
      expect(tableRows[0]).toHaveTextContent('10m');
      expect(tableRows[0]).toHaveTextContent('5');

      expect(tableRows[1]).toHaveTextContent('Memory Pressure');
      expect(tableRows[1]).toHaveTextContent('5m');
      expect(tableRows[1]).toHaveTextContent('3');
    });
  });

  describe('Prometheus rules', () => {
    it('should render vanilla prometheus rules group', async () => {
      const promDs = alertingFactory.dataSource.build({ uid: 'prometheus', name: 'Prometheus' });
      const group = alertingFactory.prometheus.group.build({ name: 'test-group-cpu', interval: 500 });
      setPrometheusRules({ uid: promDs.uid }, [group]);

      // Act
      renderGroupDetailsPage(promDs.uid, 'test-prom-namespace', 'test-group-cpu');

      // Assert
      const header = await ui.header.find();

      expect(header).toHaveTextContent('test-group-cpu');
      expect(await screen.findByText(/test-group-cpu/)).toBeInTheDocument();
      expect(await screen.findByText(/8m20s/)).toBeInTheDocument();
      expect(ui.editLink.query()).not.toBeInTheDocument();
    });
  });

  describe('Mimir rules', () => {
    it('should render mimir rules group', async () => {
      const { dataSource: mimirDs } = mimirDataSource();

      const group = alertingFactory.ruler.group.build({ name: 'test-group-cpu', interval: '11m40s' });
      setRulerRuleGroupResolver((req) => {
        if (req.params.namespace === 'test-mimir-namespace' && req.params.groupName === 'test-group-cpu') {
          return HttpResponse.json(group);
        }
        return HttpResponse.json({ error: 'Group not found' }, { status: 404 });
      });

      renderGroupDetailsPage(mimirDs.uid, 'test-mimir-namespace', 'test-group-cpu');

      const header = await ui.header.find();
      const editLink = await ui.editLink.find();

      expect(header).toHaveTextContent('test-group-cpu');
      expect(await screen.findByText(/test-mimir-namespace/)).toBeInTheDocument();
      expect(await screen.findByText(/11m40s/)).toBeInTheDocument();
      expect(editLink).toHaveAttribute(
        'href',
        `/alerting/${mimirDs.uid}/namespaces/test-mimir-namespace/groups/test-group-cpu/edit`
      );
    });
  });
});

function renderGroupDetailsPage(dsUid: string, namespaceId: string, groupName: string) {
  render(
    <Routes>
      <Route path="/alerting/:sourceId/namespaces/:namespaceId/groups/:groupName/view" element={<GroupDetailsPage />} />
    </Routes>,
    {
      historyOptions: { initialEntries: [`/alerting/${dsUid}/namespaces/${namespaceId}/groups/${groupName}/view`] },
    }
  );
}
