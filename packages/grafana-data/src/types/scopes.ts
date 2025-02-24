export interface ScopeDashboardBindingSpec {
  dashboard: string;
  scope: string;
}

export interface ScopeDashboardBindingStatus {
  dashboardTitle: string;
  groups?: string[];
}

// TODO: Use Resource from apiserver when we export the types
export interface ScopeDashboardBinding {
  metadata: {
    name: string;
  };
  spec: ScopeDashboardBindingSpec;
  status: ScopeDashboardBindingStatus;
}

export type ScopeFilterOperator = 'equals' | 'not-equals' | 'regex-match' | 'regex-not-match' | 'one-of' | 'not-one-of';
export type ScopeFilterSingleOrMultiOperator = Extract<
  ScopeFilterOperator,
  'equals' | 'not-equals' | 'one-of' | 'not-one-of'
>;

export function isScopeFilterSingleOrMultiOperator(value: string): value is ScopeFilterSingleOrMultiOperator {
  return value === 'equals' || value === 'not-equals' || value === 'one-of' || value === 'not-one-of';
}

export const scopeFilterOperatorMap: Record<string, ScopeFilterOperator> = {
  '=': 'equals',
  '!=': 'not-equals',
  '=~': 'regex-match',
  '!~': 'regex-not-match',
  '=|': 'one-of',
  '!=|': 'not-one-of',
};

export const reverseScopeFilterOperatorMap: Record<ScopeFilterOperator, string> = Object.fromEntries(
  Object.entries(scopeFilterOperatorMap).map(([symbol, operator]) => [operator, symbol])
  // eslint-disable-next-line @typescript-eslint/consistent-type-assertions
) as Record<ScopeFilterOperator, string>;

export interface ScopeSpecFilter {
  key: string;
  value: string;
  // values is used for operators that support multiple values (e.g. one-of, not-one-of)
  values?: string[];
  operator: ScopeFilterOperator;
}

export interface ScopeSpec {
  title: string;
  type: string;
  description: string;
  category: string;
  filters: ScopeSpecFilter[];
}

// TODO: Use Resource from apiserver when we export the types
export interface Scope {
  metadata: {
    name: string;
  };
  spec: ScopeSpec;
}

export type ScopeNodeNodeType = 'container' | 'leaf';
export type ScopeNodeLinkType = 'scope';

export interface ScopeNodeSpec {
  nodeType: ScopeNodeNodeType;
  title: string;

  description?: string;
  disableMultiSelect?: boolean;
  linkId?: string;
  linkType?: ScopeNodeLinkType;
}

// TODO: Use Resource from apiserver when we export the types
export interface ScopeNode {
  metadata: {
    name: string;
  };
  spec: ScopeNodeSpec;
}
