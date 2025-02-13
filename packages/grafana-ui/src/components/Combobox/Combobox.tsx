import { cx } from '@emotion/css';
import { useVirtualizer } from '@tanstack/react-virtual';
import { useCombobox } from 'downshift';
import { useId, useMemo } from 'react';

import { useStyles2 } from '../../themes';
import { t } from '../../utils/i18n';
import { Icon } from '../Icon/Icon';
import { AutoSizeInput } from '../Input/AutoSizeInput';
import { Input, Props as InputProps } from '../Input/Input';
import { Stack } from '../Layout/Stack/Stack';
import { Portal } from '../Portal/Portal';
import { ScrollContainer } from '../ScrollContainer/ScrollContainer';

import { AsyncError, NotFoundError } from './MessageRows';
import { OptionListItem } from './OptionListItem';
import { itemToString } from './filter';
import { getComboboxStyles, MENU_OPTION_HEIGHT, MENU_OPTION_HEIGHT_DESCRIPTION } from './getComboboxStyles';
import { ComboboxOption } from './types';
import { useComboboxFloat } from './useComboboxFloat';
import { useOptions } from './useOptions';
import { isNewGroup } from './utils';

// TODO: It would be great if ComboboxOption["label"] was more generic so that if consumers do pass it in (for async),
// then the onChange handler emits ComboboxOption with the label as non-undefined.
export interface ComboboxBaseProps<T extends string | number>
  extends Pick<
    InputProps,
    'placeholder' | 'autoFocus' | 'id' | 'aria-labelledby' | 'disabled' | 'loading' | 'invalid'
  > {
  /**
   * An `X` appears in the UI, which clears the input and sets the value to `null`. Do not use if you have no `null` case.
   */
  isClearable?: boolean;
  /**
   * Allows the user to set a value which is not in the list of options.
   */
  createCustomValue?: boolean;
  options: Array<ComboboxOption<T>> | ((inputValue: string) => Promise<Array<ComboboxOption<T>>>);
  onChange: (option: ComboboxOption<T>) => void;
  /**
   * Most consumers should pass value in as a scalar string | number. However, sometimes with Async because we don't
   * have the full options loaded to match the value to, consumers may also pass in an Option with a label to display.
   */
  value?: T | ComboboxOption<T> | null;
  /**
   * Defaults to 100%. Number is a multiple of 8px. 'auto' will size the input to the content.
   * */
  width?: number | 'auto';
  onBlur?: () => void;
}

type ClearableConditionals<T extends number | string> =
  | {
      isClearable: true;
      /**
       * The onChange handler is called with `null` when clearing the Combobox.
       */
      onChange: (option: ComboboxOption<T> | null) => void;
    }
  | { isClearable?: false; onChange: (option: ComboboxOption<T>) => void };

export type AutoSizeConditionals =
  | {
      width: 'auto';
      /**
       * Needs to be set when width is 'auto' to prevent the input from shrinking too much
       */
      minWidth: number;
      /**
       * Recommended to set when width is 'auto' to prevent the input from growing too much.
       */
      maxWidth?: number;
    }
  | {
      width?: number;
      minWidth?: never;
      maxWidth?: never;
    };

export type ComboboxProps<T extends string | number> = ComboboxBaseProps<T> &
  AutoSizeConditionals &
  ClearableConditionals<T>;

const noop = () => {};

export const VIRTUAL_OVERSCAN_ITEMS = 4;

/**
 * A performant Select replacement.
 *
 * @alpha
 */
export const Combobox = <T extends string | number>(props: ComboboxProps<T>) => {
  const {
    // options,
    onChange,
    value: valueProp,
    placeholder: placeholderProp,
    isClearable = false,
    createCustomValue = false,
    id,
    width,
    minWidth,
    maxWidth,
    'aria-labelledby': ariaLabelledBy,
    autoFocus,
    onBlur,
    disabled,
    loading,
    invalid,
  } = props;

  // Value can be an actual scalar Value (string or number), or an Option (value + label), so
  // get a consistent Value from it
  const value = typeof valueProp === 'object' ? valueProp?.value : valueProp;

  const allOptions = props.options;
  const {
    options: filteredOptions,
    updateOptions,
    asyncLoading,
    asyncError,
  } = useOptions(props.options, createCustomValue);
  const isAsync = typeof allOptions === 'function';

  const selectedItemIndex = useMemo(() => {
    if (isAsync) {
      return null;
    }

    if (valueProp === undefined || valueProp === null) {
      return null;
    }

    const index = allOptions.findIndex((option) => option.value === value);
    if (index === -1) {
      return null;
    }

    return index;
  }, [valueProp, allOptions, value, isAsync]);

  const selectedItem = useMemo(() => {
    if (valueProp === undefined || valueProp === null) {
      return null;
    }

    if (selectedItemIndex !== null && !isAsync) {
      return allOptions[selectedItemIndex];
    }

    return typeof valueProp === 'object' ? valueProp : { value: valueProp, label: valueProp.toString() };
  }, [selectedItemIndex, isAsync, valueProp, allOptions]);

  const menuId = `downshift-${useId().replace(/:/g, '--')}-menu`;
  const labelId = `downshift-${useId().replace(/:/g, '--')}-label`;

  const styles = useStyles2(getComboboxStyles);

  const virtualizerOptions = {
    count: filteredOptions.length,
    getScrollElement: () => scrollRef.current,
    // estimateSize: (index: number) =>
    //   filteredOptions[index].description ? MENU_OPTION_HEIGHT_DESCRIPTION : MENU_OPTION_HEIGHT,
    estimateSize: (index: number) => {
      const firstGroupItem = isNewGroup(filteredOptions[index], index > 0 ? filteredOptions[index - 1] : undefined);
      const hasDescription = 'description' in filteredOptions[index];
      let itemHeight = MENU_OPTION_HEIGHT;
      if (hasDescription) {
        itemHeight = MENU_OPTION_HEIGHT_DESCRIPTION;
      }
      if (firstGroupItem) {
        itemHeight += MENU_OPTION_HEIGHT;
      }
      return itemHeight;
    },
    overscan: VIRTUAL_OVERSCAN_ITEMS,
  };

  const rowVirtualizer = useVirtualizer(virtualizerOptions);

  const {
    isOpen,
    highlightedIndex,

    getInputProps,
    getMenuProps,
    getItemProps,

    selectItem,
  } = useCombobox({
    menuId,
    labelId,
    inputId: id,
    items: filteredOptions,
    itemToString,
    selectedItem,

    // Don't change downshift state in the onBlahChange handlers. Instead, use the stateReducer to make changes.
    // Downshift calls change handlers on the render after so you can get sync/flickering issues if you change its state
    // in them.
    // Instead, stateReducer is called in the same tick as state changes, before that state is committed and rendered.

    onSelectedItemChange: ({ selectedItem }) => {
      onChange(selectedItem);
    },

    defaultHighlightedIndex: selectedItemIndex ?? 0,

    scrollIntoView: () => {},

    onIsOpenChange: ({ isOpen, inputValue }) => {
      if (isOpen && inputValue === '') {
        updateOptions(inputValue);
      }
    },

    onHighlightedIndexChange: ({ highlightedIndex, type }) => {
      if (type !== useCombobox.stateChangeTypes.MenuMouseLeave) {
        rowVirtualizer.scrollToIndex(highlightedIndex);
      }
    },
    onStateChange: ({ inputValue: newInputValue, type, selectedItem: newSelectedItem }) => {
      switch (type) {
        case useCombobox.stateChangeTypes.InputChange:
          updateOptions(newInputValue ?? '');

          break;
        default:
          break;
      }
    },
    stateReducer(state, actionAndChanges) {
      let { changes } = actionAndChanges;
      const menuBeingOpened = state.isOpen === false && changes.isOpen === true;
      const menuBeingClosed = state.isOpen === true && changes.isOpen === false;

      // Reset the input value when the menu is opened. If the menu is opened due to an input change
      // then make sure we keep that.
      // This will trigger onInputValueChange to load async options
      if (menuBeingOpened && changes.inputValue === state.inputValue) {
        changes = {
          ...changes,
          inputValue: '',
        };
      }

      if (menuBeingClosed) {
        // Flush the selected item to the input when the menu is closed
        if (changes.selectedItem) {
          changes = {
            ...changes,
            inputValue: itemToString(changes.selectedItem),
          };
        } else if (changes.inputValue !== '') {
          // Otherwise if no selected value, clear any search from the input
          changes = {
            ...changes,
            inputValue: '',
          };
        }
      }

      return changes;
    },
  });

  const { inputRef, floatingRef, floatStyles, scrollRef } = useComboboxFloat(filteredOptions, isOpen);

  const isAutoSize = width === 'auto';

  const InputComponent = isAutoSize ? AutoSizeInput : Input;

  const suffixIcon = asyncLoading
    ? 'spinner'
    : // If it's loading, show loading icon. Otherwise, icon indicating menu state
      isOpen
      ? 'search'
      : 'angle-down';

  const placeholder = (isOpen ? itemToString(selectedItem) : null) || placeholderProp;

  return (
    <div className={isAutoSize ? styles.addaptToParent : undefined}>
      <InputComponent
        width={isAutoSize ? undefined : width}
        {...(isAutoSize ? { minWidth, maxWidth } : {})}
        autoFocus={autoFocus}
        onBlur={onBlur}
        disabled={disabled}
        loading={loading}
        invalid={invalid}
        className={styles.input}
        suffix={
          <>
            {!!value && value === selectedItem?.value && isClearable && (
              <Icon
                name="times"
                className={styles.clear}
                title={t('combobox.clear.title', 'Clear value')}
                tabIndex={0}
                role="button"
                onClick={() => {
                  selectItem(null);
                }}
                onKeyDown={(e) => {
                  if (e.key === 'Enter' || e.key === ' ') {
                    selectItem(null);
                  }
                }}
              />
            )}

            <Icon name={suffixIcon} />
          </>
        }
        {...getInputProps({
          ref: inputRef,
          /*  Empty onCall to avoid TS error
           *  See issue here: https://github.com/downshift-js/downshift/issues/718
           *  Downshift repo: https://github.com/downshift-js/downshift/tree/master
           */
          onChange: noop,
          'aria-labelledby': ariaLabelledBy, // Label should be handled with the Field component
          placeholder,
        })}
      />
      <Portal>
        <div
          className={cx(styles.menu, !isOpen && styles.menuClosed)}
          style={{
            ...floatStyles,
          }}
          {...getMenuProps({
            ref: floatingRef,
            'aria-labelledby': ariaLabelledBy,
          })}
        >
          {isOpen && (
            <ScrollContainer showScrollIndicators maxHeight="inherit" ref={scrollRef} padding={0.5}>
              {!asyncError && (
                <ul style={{ height: rowVirtualizer.getTotalSize() }} className={styles.menuUlContainer}>
                  {rowVirtualizer.getVirtualItems().map((virtualRow) => {
                    const item = filteredOptions[virtualRow.index];
                    const startingNewGroup = isNewGroup(item, filteredOptions[virtualRow.index - 1]);
                    const itemId = 'multicombobox-option-' + item.value.toString(); // TODO
                    const groupHeaderid = 'multicombobox-option-group-' + item.value.toString(); // TODO

                    return (
                      <li
                        key={`${item.value}-${virtualRow.index}`}
                        data-index={virtualRow.index}
                        className={styles.optionBasic}
                        style={{
                          height: virtualRow.size,
                          transform: `translateY(${virtualRow.start}px)`,
                        }}
                        {...getItemProps({
                          item: item,
                          index: virtualRow.index,
                        })}
                      >
                        <Stack direction="column" justifyContent="space-between" width="100%" height="100%" gap={0}>
                          {startingNewGroup && (
                            <div className={styles.optionGroup}>
                              <OptionListItem
                                label={item.group ?? t('combobox.group.undefined', 'No group')}
                                id={groupHeaderid}
                                isGroup={true}
                              />
                            </div>
                          )}

                          <div
                            className={cx(
                              styles.option,
                              selectedItem && item.value === selectedItem.value && styles.optionSelected,
                              highlightedIndex === virtualRow.index && styles.optionFocused
                            )}
                          >
                            <OptionListItem
                              label={item.label ?? item.value}
                              description={item.description}
                              id={itemId}
                              isGroup={false}
                            />
                          </div>
                        </Stack>
                      </li>
                    );
                  })}
                </ul>
              )}
              <div aria-live="polite">
                {asyncError && <AsyncError />}
                {filteredOptions.length === 0 && !asyncError && <NotFoundError />}
              </div>
            </ScrollContainer>
          )}
        </div>
      </Portal>
    </div>
  );
};
