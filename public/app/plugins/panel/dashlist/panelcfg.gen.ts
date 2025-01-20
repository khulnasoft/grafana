// Code generated - EDITING IS FUTILE. DO NOT EDIT.
//
// Generated by:
//     public/app/plugins/gen.go
// Using jennies:
//     PluginTsTypesJenny
//     PluginTsTypesJenny
//
// Run 'make gen-cue' from repository root to regenerate.

export interface Options {
  /**
   * folderId is deprecated, and migrated to folderUid on panel init
   */
  folderId?: number;
  folderUID?: string;
  includeVars: boolean;
  keepTime: boolean;
  maxItems: number;
  query: string;
  showFolderNames: boolean;
  showHeadings: boolean;
  showRecentlyViewed: boolean;
  showSearch: boolean;
  showStarred: boolean;
  tags: Array<string>;
}

export const defaultOptions: Partial<Options> = {
  includeVars: false,
  keepTime: false,
  maxItems: 10,
  query: '',
  showFolderNames: true,
  showHeadings: true,
  showRecentlyViewed: false,
  showSearch: false,
  showStarred: true,
  tags: [],
};
