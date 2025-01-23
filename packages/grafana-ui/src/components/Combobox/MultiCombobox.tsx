import { cx } from '@emotion/css';
import { useVirtualizer } from '@tanstack/react-virtual';
import { useCombobox, useMultipleSelection } from 'downshift';
import { useCallback, useMemo, useState } from 'react';

import { useStyles2 } from '../../themes';
import { t } from '../../utils/i18n';
import { Checkbox } from '../Forms/Checkbox';
import { Box } from '../Layout/Box/Box';
import { Stack } from '../Layout/Stack/Stack';
import { Portal } from '../Portal/Portal';
import { ScrollContainer } from '../ScrollContainer/ScrollContainer';
import { Text } from '../Text/Text';
import { Tooltip } from '../Tooltip';

import { ComboboxBaseProps, AutoSizeConditionals, VIRTUAL_OVERSCAN_ITEMS } from './Combobox';
import { NotFoundError } from './MessageRows';
import { OptionListItem } from './OptionListItem';
import { SuffixIcon } from './SuffixIcon';
import { ValuePill } from './ValuePill';
import { itemToString } from './filter';
import { getComboboxStyles, MENU_OPTION_HEIGHT, MENU_OPTION_HEIGHT_DESCRIPTION } from './getComboboxStyles';
import { getMultiComboboxStyles } from './getMultiComboboxStyles';
import { ALL_OPTION_VALUE, ComboboxOption } from './types';
import { useComboboxFloat } from './useComboboxFloat';
import { MAX_SHOWN_ITEMS, useMeasureMulti } from './useMeasureMulti';
import { useMultiInputAutoSize } from './useMultiInputAutoSize';
import { useOptions } from './useOptions';

interface MultiComboboxBaseProps<T extends string | number> extends Omit<ComboboxBaseProps<T>, 'value' | 'onChange'> {
  value?: T[] | Array<ComboboxOption<T>>;
  onChange: (items?: T[]) => void;
  enableAllOption?: boolean;
}

export type MultiComboboxProps<T extends string | number> = MultiComboboxBaseProps<T> & AutoSizeConditionals;

export const MultiCombobox = <T extends string | number>(props: MultiComboboxProps<T>) => {
  const { placeholder, onChange, value, width, enableAllOption, invalid, loading, disabled, minWidth, maxWidth } =
    props;

  const { options: newOptions, loadOptionsWhenUserTypes /*, asyncLoading, asyncError */ } = useOptions(props.options);

  const selectedItems = useMemo(() => {
    if (!value) {
      return [];
    }

    return getSelectedItemsFromValue<T>(value, newOptions);
  }, [value, newOptions]);

  const styles = useStyles2(getComboboxStyles);
  const [inputValue, setInputValue] = useState('');

  const allOptionItem = useMemo(() => {
    return {
      label:
        inputValue === ''
          ? t('multicombobox.all.title', 'All')
          : t('multicombobox.all.title-filtered', 'All (filtered)'),
      // Type casting needed to make this work when T is a number
      value: ALL_OPTION_VALUE,
    } as ComboboxOption<T>;
  }, [inputValue]);

  const baseItems = useMemo(() => {
    return enableAllOption ? [allOptionItem, ...newOptions] : newOptions;
  }, [newOptions, enableAllOption, allOptionItem]);

  const items = useMemo(() => {
    if (enableAllOption && baseItems.length === 1 && baseItems[0] === allOptionItem) {
      return [];
    }

    return baseItems;
  }, [baseItems, enableAllOption, allOptionItem]);

  const { measureRef, counterMeasureRef, suffixMeasureRef, shownItems } = useMeasureMulti(
    selectedItems,
    width,
    disabled
  );

  const isOptionSelected = useCallback(
    (item: ComboboxOption<T>) => selectedItems.some((opt) => opt.value === item.value),
    [selectedItems]
  );

  const { getSelectedItemProps, getDropdownProps, removeSelectedItem } = useMultipleSelection({
    selectedItems, //initally selected items,
    onStateChange: ({ type, selectedItems: newSelectedItems }) => {
      switch (type) {
        case useMultipleSelection.stateChangeTypes.SelectedItemKeyDownBackspace:
        case useMultipleSelection.stateChangeTypes.SelectedItemKeyDownDelete:
        case useMultipleSelection.stateChangeTypes.DropdownKeyDownBackspace:
        case useMultipleSelection.stateChangeTypes.FunctionRemoveSelectedItem:
          if (newSelectedItems) {
            onChange(getComboboxOptionsValues(newSelectedItems));
          }
          break;

        default:
          break;
      }
    },
    stateReducer: (state, actionAndChanges) => {
      const { changes } = actionAndChanges;
      return {
        ...changes,

        /**
         * TODO: Fix Hack!
         * This prevents the menu from closing when the user unselects an item in the dropdown at the expense
         * of breaking keyboard navigation.
         *
         * Downshift isn't really designed to keep selected items in the dropdown menu, so when you unselect an item
         * in a multiselect, the stateReducer tries to move focus onto another item which causes the menu to be closed.
         * This only seems to happen when you deselect the last item in the selectedItems list.
         *
         * Check out:
         *  - FunctionRemoveSelectedItem in the useMultipleSelection reducer https://github.com/downshift-js/downshift/blob/master/src/hooks/useMultipleSelection/reducer.js#L75
         *  - The activeIndex useEffect in useMultipleSelection https://github.com/downshift-js/downshift/blob/master/src/hooks/useMultipleSelection/index.js#L68-L72
         *
         * Forcing the activeIndex to -999 both prevents the useEffect that changes the focus from triggering (value never changes)
         * and prevents the if statement in useMultipleSelection from focusing anything.
         */
        activeIndex: -999,
      };
    },
  });

  const {
    getToggleButtonProps,
    //getLabelProps,
    isOpen,
    highlightedIndex,
    getMenuProps,
    getInputProps,
    getItemProps,
  } = useCombobox({
    items,
    itemToString,
    inputValue,
    selectedItem: null,
    stateReducer: (state, actionAndChanges) => {
      const { type } = actionAndChanges;
      let { changes } = actionAndChanges;
      const menuBeingOpened = state.isOpen === false && changes.isOpen === true;
      // const menuBeingClosed = state.isOpen === true && changes.isOpen === false;

      // Reset the input value when the menu is opened. If the menu is opened due to an input change
      // then make sure we keep that.
      // This will trigger onInputValueChange to load async options
      if (menuBeingOpened && changes.inputValue === state.inputValue) {
        changes = {
          ...changes,
          inputValue: '',
        };
      }

      // TODO: handle menuBeingClosed?

      switch (type) {
        case useCombobox.stateChangeTypes.InputKeyDownEnter:
        case useCombobox.stateChangeTypes.ItemClick:
          return {
            ...changes,
            isOpen: true,
            highlightedIndex: state.highlightedIndex,
          };
        case useCombobox.stateChangeTypes.InputBlur:
          setInputValue('');
        default:
          return changes;
      }
    },

    onIsOpenChange: ({ isOpen, inputValue }) => {
      console.log('onIsOpenChange', { isOpen, inputValue });
      if (isOpen && inputValue === '') {
        loadOptionsWhenUserTypes(inputValue);
      }
    },

    onStateChange: ({ inputValue: newInputValue, type, selectedItem: newSelectedItem }) => {
      console.log('state change', type);
      switch (type) {
        case useCombobox.stateChangeTypes.InputKeyDownEnter:
        case useCombobox.stateChangeTypes.ItemClick:
          // Handle All functionality
          if (newSelectedItem?.value === ALL_OPTION_VALUE) {
            const allFilteredSelected = selectedItems.length === items.length - 1;
            let newSelectedItems = allFilteredSelected && inputValue === '' ? [] : baseItems.slice(1);

            if (!allFilteredSelected && inputValue !== '') {
              // Select all currently filtered items and deduplicate
              newSelectedItems = [...new Set([...selectedItems, ...items.slice(1)])];
            }

            if (allFilteredSelected && inputValue !== '') {
              // Deselect all currently filtered items
              const filteredSet = new Set(items.slice(1).map((item) => item.value));
              newSelectedItems = selectedItems.filter((item) => !filteredSet.has(item.value));
            }

            onChange(getComboboxOptionsValues(newSelectedItems));
            break;
          }
          if (newSelectedItem) {
            if (!isOptionSelected(newSelectedItem)) {
              onChange(getComboboxOptionsValues([...selectedItems, newSelectedItem]));
              break;
            }
            removeSelectedItem(newSelectedItem); // onChange is handled by multiselect here
          }
          break;
        case useCombobox.stateChangeTypes.InputChange:
          setInputValue(newInputValue ?? '');
          loadOptionsWhenUserTypes(newInputValue ?? '');

          break;
        default:
          break;
      }
    },
  });

  const { inputRef: containerRef, floatingRef, floatStyles, scrollRef } = useComboboxFloat(items, isOpen);
  const multiStyles = useStyles2(getMultiComboboxStyles, isOpen, invalid, disabled, width, minWidth, maxWidth);

  const virtualizerOptions = {
    count: items.length,
    getScrollElement: () => scrollRef.current,
    estimateSize: (index: number) =>
      'description' in items[index] ? MENU_OPTION_HEIGHT_DESCRIPTION : MENU_OPTION_HEIGHT,
    overscan: VIRTUAL_OVERSCAN_ITEMS,
  };

  const rowVirtualizer = useVirtualizer(virtualizerOptions);

  // Selected items that show up in the input field
  const visibleItems = isOpen ? selectedItems.slice(0, MAX_SHOWN_ITEMS) : selectedItems.slice(0, shownItems);

  const { inputRef, inputWidth } = useMultiInputAutoSize(inputValue);
  return (
    <div className={multiStyles.container} ref={containerRef}>
      <div className={cx(multiStyles.wrapper, { [multiStyles.disabled]: disabled })} ref={measureRef}>
        <span className={multiStyles.pillWrapper}>
          {visibleItems.map((item, index) => (
            <ValuePill
              disabled={disabled}
              onRemove={() => {
                removeSelectedItem(item);
              }}
              key={`${item.value}${index}`}
              {...getSelectedItemProps({ selectedItem: item, index })}
            >
              {itemToString(item)}
            </ValuePill>
          ))}
          {selectedItems.length > visibleItems.length && (
            <Box display="flex" direction="row" marginLeft={0.5} gap={1} ref={counterMeasureRef}>
              {/* eslint-disable-next-line @grafana/no-untranslated-strings */}
              <Text>...</Text>
              <Tooltip
                interactive
                content={
                  <>
                    {selectedItems.slice(visibleItems.length).map((item) => (
                      <div key={item.value}>{itemToString(item)}</div>
                    ))}
                  </>
                }
              >
                <div className={multiStyles.restNumber}>{selectedItems.length - shownItems}</div>
              </Tooltip>
            </Box>
          )}
          <input
            className={multiStyles.input}
            {...getInputProps(
              getDropdownProps({
                disabled,
                preventKeyAction: isOpen,
                placeholder,
                ref: inputRef,
                style: { width: inputWidth },
              })
            )}
          />

          <div className={multiStyles.suffix} ref={suffixMeasureRef} {...getToggleButtonProps()}>
            <SuffixIcon isLoading={loading || false} isOpen={isOpen} />
          </div>
        </span>
      </div>
      <Portal>
        <div
          className={cx(styles.menu, !isOpen && styles.menuClosed)}
          style={{ ...floatStyles }}
          {...getMenuProps({ ref: floatingRef })}
        >
          {isOpen && (
            <ScrollContainer showScrollIndicators maxHeight="inherit" ref={scrollRef}>
              <ul style={{ height: rowVirtualizer.getTotalSize() }} className={styles.menuUlContainer}>
                {rowVirtualizer.getVirtualItems().map((virtualRow) => {
                  const index = virtualRow.index;
                  const item = items[index];
                  const itemProps = getItemProps({ item, index });
                  const isSelected = isOptionSelected(item);
                  const id = 'multicombobox-option-' + item.value.toString();
                  const isAll = item.value === ALL_OPTION_VALUE;
                  const allItemsSelected =
                    items[0]?.value === ALL_OPTION_VALUE && selectedItems.length === items.length - 1;

                  return (
                    <li
                      key={`${item.value}-${index}`}
                      data-index={index}
                      {...itemProps}
                      className={cx(styles.option, { [styles.optionFocused]: highlightedIndex === index })}
                      style={{ height: virtualRow.size, transform: `translateY(${virtualRow.start}px)` }}
                    >
                      <Stack direction="row" alignItems="center">
                        <Checkbox
                          key={id}
                          value={allItemsSelected || isSelected}
                          indeterminate={isAll && selectedItems.length > 0 && !allItemsSelected}
                          aria-labelledby={id}
                          onClick={(e) => {
                            e.stopPropagation();
                          }}
                        />
                        <OptionListItem
                          label={
                            isAll
                              ? (item.label ?? item.value.toString()) +
                                (isAll && inputValue !== '' ? ` (${items.length - 1})` : '')
                              : (item.label ?? item.value.toString())
                          }
                          description={item?.description}
                          id={id}
                        />
                      </Stack>
                    </li>
                  );
                })}
              </ul>
              <div aria-live="polite">{items.length === 0 && <NotFoundError />}</div>
            </ScrollContainer>
          )}
        </div>
      </Portal>
    </div>
  );
};

function getSelectedItemsFromValue<T extends string | number>(
  value: T[] | Array<ComboboxOption<T>>,
  options: Array<ComboboxOption<T>>
) {
  if (isComboboxOptions(value)) {
    return value;
  }
  const valueMap = new Map(value.map((val, index) => [val, index]));
  const resultingItems: Array<ComboboxOption<T>> = [];

  for (const option of options) {
    const index = valueMap.get(option.value);
    if (index !== undefined) {
      resultingItems[index] = option;
      valueMap.delete(option.value);
    }
    if (valueMap.size === 0) {
      // We found all values
      break;
    }
  }

  // Handle items that are not in options
  for (const [val, index] of valueMap) {
    resultingItems[index] = { value: val };
  }
  return resultingItems;
}

function isComboboxOptions<T extends string | number>(
  value: T[] | Array<ComboboxOption<T>>
): value is Array<ComboboxOption<T>> {
  return typeof value[0] === 'object';
}

function getComboboxOptionsValues<T extends string | number>(optionArray: Array<ComboboxOption<T>>) {
  return optionArray.map((option) => option.value);
}
