import { useMemo, useState } from 'react';
import { Check, ChevronsUpDown } from 'lucide-react';
import { useDebouncedCallback } from 'use-debounce';
import { useQuery, type UseQueryOptions } from '@tanstack/react-query';

import { cn } from '@/lib/utils';
import { Button } from '@/components/ui/button';
import {
  Command,
  CommandEmpty,
  CommandGroup,
  CommandInput,
  CommandItem,
  CommandList,
} from '@/components/ui/command';
import {
  Popover,
  PopoverContent,
  PopoverTrigger,
} from '@/components/ui/popover';
import { ALIGN_OPTIONS } from '@radix-ui/react-popper';

export type ComboBoxItem = {
  value: string;
  label: string;
};

export type ApiComboboxProps<
  TQueryFnData,
  TError,
  TData,
  TQueryKey extends unknown[] = unknown[],
> = {
  selectedItem?: ComboBoxItem;
  onSelect: (item: ComboBoxItem) => void;
  queryOptionsFactory: (
    search: string
  ) => UseQueryOptions<TQueryFnData, TError, TData, TQueryKey>;
  mapDataToItems: (data: TData) => ComboBoxItem[];
  searchPlaceholder?: string;
  selectPlaceholder?: string;
  className?: string;
  disabled?: boolean;
  align?: (typeof ALIGN_OPTIONS)[number];
};

export function ApiComboBox<
  TQueryFnData,
  TError,
  TData,
  TQueryKey extends unknown[],
>({
  selectedItem,
  onSelect,
  queryOptionsFactory,
  mapDataToItems,
  searchPlaceholder = 'Search...',
  selectPlaceholder = 'Select an item',
  className,
  disabled = false,
  align,
}: ApiComboboxProps<TQueryFnData, TError, TData, TQueryKey>) {
  const [open, setOpen] = useState(false);
  const [search, setSearch] = useState('');

  const debouncedSearch = useDebouncedCallback((value: string) => {
    setSearch(value);
  }, 300);

  const queryOptions = queryOptionsFactory(search);
  const { data, isLoading } = useQuery(queryOptions);

  const items = useMemo(() => {
    if (!data) return [];
    return mapDataToItems(data as TData);
  }, [data, mapDataToItems]);

  return (
    <Popover open={open} onOpenChange={setOpen} modal>
      <PopoverTrigger asChild>
        <Button
          variant="outline"
          role="combobox"
          aria-expanded={open}
          className={cn('justify-between', className)}
          disabled={disabled}
        >
          <span className="truncate">
            {selectedItem?.label || selectPlaceholder}
          </span>
          <ChevronsUpDown className="ml-2 h-4 w-4 shrink-0 opacity-50" />
        </Button>
      </PopoverTrigger>
      <PopoverContent
        style={{ width: 'var(--radix-popover-trigger-width)' }}
        className="p-0"
        align={align}
      >
        <Command shouldFilter={false}>
          <CommandInput
            placeholder={searchPlaceholder}
            onValueChange={debouncedSearch}
          />
          <CommandList>
            {isLoading && (
              <div className="text-muted-foreground p-4 text-sm">
                Nalagam...
              </div>
            )}
            <CommandEmpty>Ni rezultatov.</CommandEmpty>
            <CommandGroup>
              {items.map((item) => {
                const isSelected = selectedItem?.value === item.value;
                return (
                  <CommandItem
                    key={item.value}
                    value={item.value}
                    onSelect={() => {
                      onSelect(item);
                      setOpen(false);
                    }}
                  >
                    {item.label}
                    <Check
                      className={cn(
                        'ml-auto h-4 w-4',
                        isSelected ? 'opacity-100' : 'opacity-0'
                      )}
                    />
                  </CommandItem>
                );
              })}
            </CommandGroup>
          </CommandList>
        </Command>
      </PopoverContent>
    </Popover>
  );
}
