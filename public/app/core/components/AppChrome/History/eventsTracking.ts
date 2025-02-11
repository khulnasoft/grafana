import { reportInteraction } from '@grafana/runtime';

const UNIFIED_HISTORY_ENTRY_CLICKED = 'grafana_unified_history_entry_clicked';
const UNIFIED_HISTORY_ENTRY_DUPLICATED = 'grafana_unified_history_entry_duplicated';

//Currently just 'timeRange' is supported
//in short term, we could add 'templateVariables' for example
type subEntryTypes = 'timeRange';

//Event triggered when a user clicks on an entry of the `HistoryDrawer`
interface UnifiedHistoryEntryClicked {
  //We will also work with the current URL but we will get this from Rudderstack data
  //URL to return to
  entryURL: string;
  //In the case we want to go back to a specific query param, currently just a specific time range
  subEntry?: subEntryTypes;
}

//Event triggered when history entry name matches the previous one
//so we keep track of duplicated entries and be able to analyze them
interface UnifiedHistoryEntryDuplicated {
  // Common name of the history entries
  entryName: string;
  // URL of the last entry
  lastEntryName: string;
  // URL of the new entry
  newEntryURL: string;
}

export const clickUnifiedHistoryEntry = ({ entryURL, subEntry }: UnifiedHistoryEntryClicked) => {
  reportInteraction(UNIFIED_HISTORY_ENTRY_CLICKED, {
    entryURL,
    subEntry,
  });
};

export const duplicateUnifiedHistoryEntry = ({
  entryName,
  lastEntryName,
  newEntryURL,
}: UnifiedHistoryEntryDuplicated) => {
  reportInteraction(UNIFIED_HISTORY_ENTRY_DUPLICATED, {
    entryName,
    lastEntryName,
    newEntryURL,
  });
};
