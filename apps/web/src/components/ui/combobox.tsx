'use client';

import * as React from 'react';
import { Check, ChevronsUpDown } from 'lucide-react';
import { useVirtualizer } from '@tanstack/react-virtual';

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

type ComboboxOption = {
  value: string;
  label: string;
};

type ComboboxProps = {
  options: ComboboxOption[];
  value: string;
  onChange: (value: string) => void;
  placeholder?: string;
  searchPlaceholder?: string;
  className?: string;
  itemHeight?: number;
  maxHeight?: number;
};

const Combobox: React.FC<ComboboxProps> = ({
  options,
  value,
  onChange,
  placeholder = 'Select...',
  searchPlaceholder = 'Search...',
  className,
  itemHeight = 40,
  maxHeight = 300,
}) => {
  const [open, setOpen] = React.useState(false);
  const [search, setSearch] = React.useState('');
  const listRef = React.useRef<HTMLDivElement | null>(null);
  const [scrollElement, setScrollElement] = React.useState<HTMLElement | null>(
    null
  );

  const filteredOptions = React.useMemo(() => {
    return options.filter((option) =>
      option.label.toLowerCase().includes(search.toLowerCase())
    );
  }, [options, search]);

  const virtualizer = useVirtualizer({
    count: filteredOptions.length,
    getScrollElement: () => scrollElement,
    estimateSize: () => itemHeight,
    overscan: 10,
  });

  const virtualItems = virtualizer.getVirtualItems();
  const selectedLabel = options.find((opt) => opt.value === value)?.label;

  return (
    <Popover open={open} onOpenChange={setOpen}>
      <PopoverTrigger asChild>
        <Button
          variant="outline"
          role="combobox"
          aria-expanded={open}
          className={cn('w-[200px] justify-between', className)}
        >
          {selectedLabel || placeholder}
          <ChevronsUpDown className="ml-2 h-4 w-4 opacity-50" />
        </Button>
      </PopoverTrigger>
      <PopoverContent
        style={{ width: 'var(--radix-popover-trigger-width)' }}
        className="p-0"
      >
        <Command>
          <CommandInput
            placeholder={searchPlaceholder}
            className="h-9"
            value={search}
            onValueChange={setSearch}
          />
          <CommandList>
            <CommandEmpty>No results found.</CommandEmpty>
            <CommandGroup>
              <div
                ref={(el) => {
                  listRef.current = el;
                  if (el) setScrollElement(el);
                }}
                style={{
                  maxHeight,
                  overflowY: 'auto',
                  position: 'relative',
                  height: Math.min(virtualizer.getTotalSize(), maxHeight),
                }}
              >
                <div
                  style={{
                    height: virtualizer.getTotalSize(),
                    width: '100%',
                    position: 'relative',
                  }}
                >
                  {virtualItems.map((virtualRow) => {
                    const option = filteredOptions[virtualRow.index];
                    if (!option) return null;

                    return (
                      <CommandItem
                        key={option.value}
                        onSelect={() => {
                          onChange(option.value === value ? '' : option.value);
                          setOpen(false);
                          setSearch('');
                        }}
                        style={{
                          position: 'absolute',
                          top: 0,
                          left: 0,
                          transform: `translateY(${virtualRow.start}px)`,
                          height: `${itemHeight}px`,
                          width: '100%',
                        }}
                      >
                        <Check
                          className={cn(
                            'mr-2 h-4 w-4',
                            value === option.value ? 'opacity-100' : 'opacity-0'
                          )}
                        />
                        {option.label}
                      </CommandItem>
                    );
                  })}
                </div>
              </div>
            </CommandGroup>
          </CommandList>
        </Command>
      </PopoverContent>
    </Popover>
  );
};

export default Combobox;
