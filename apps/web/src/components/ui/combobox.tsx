import { type FC, type ReactNode, useMemo, useState } from 'react';
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

export type ComboboxOption = {
  value: string;
  label: string;
  render?: () => ReactNode;
};

type ComboboxProps = {
  options: ComboboxOption[];
  placeholder?: string;
  searchPlaceholder?: string;
  emptyPlaceholder?: string;
  selected?: string;
  onChange?: (value: string) => void;
  className?: string;
};

export const Combobox: FC<ComboboxProps> = ({
  options,
  placeholder = 'Select an option...',
  searchPlaceholder = 'Search...',
  emptyPlaceholder,
  selected = '',
  onChange,
  className,
}) => {
  const [open, setOpen] = useState(false);
  const [value, setValue] = useState(selected);
  const [searchTerm, setSearchTerm] = useState('');
  const [scrollElement, setScrollElement] = useState<HTMLElement | null>(null);

  const filteredOptions = useMemo(() => {
    return options.filter((opt) =>
      opt.label.toLowerCase().includes(searchTerm.toLowerCase())
    );
  }, [searchTerm, options]);

  const handleSelect = (selectedValue: string) => {
    const newValue = selectedValue === value ? '' : selectedValue;
    setValue(newValue);
    setOpen(false);
    onChange?.(newValue);
  };

  const selectedLabel = options.find((opt) => opt.value === value)?.label;

  const virtualizer = useVirtualizer({
    count: filteredOptions.length,
    getScrollElement: () => scrollElement,
    estimateSize: () => 32,
    overscan: 5,
  });

  const virtualItems = scrollElement ? virtualizer.getVirtualItems() : [];
  const totalSize = scrollElement ? virtualizer.getTotalSize() : 0;

  return (
    <Popover
      open={open}
      onOpenChange={(nextOpen) => {
        setOpen(nextOpen);
        if (!nextOpen) setSearchTerm('');
      }}
    >
      <PopoverTrigger asChild>
        <Button
          variant="outline"
          role="combobox"
          aria-expanded={open}
          className={cn('w-full justify-between', className)}
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
            value={searchTerm}
            onValueChange={setSearchTerm}
          />
          <CommandList
            ref={(el) => {
              if (el) setScrollElement(el);
            }}
            style={{ maxHeight: '300px', overflow: 'auto' }}
          >
            <CommandEmpty>{emptyPlaceholder}</CommandEmpty>
            <CommandGroup>
              <div
                style={{
                  height: `${totalSize}px`,
                  position: 'relative',
                }}
              >
                {virtualItems.map((virtualRow) => {
                  const option = filteredOptions[virtualRow.index];
                  return (
                    <CommandItem
                      key={option.value}
                      value={option.label}
                      onSelect={() => handleSelect(option.value)}
                      style={{
                        position: 'absolute',
                        top: 0,
                        left: 0,
                        width: '100%',
                        transform: `translateY(${virtualRow.start}px)`,
                      }}
                    >
                      <span>
                        {option?.render ? option.render() : option.label}
                      </span>
                      <Check
                        className={cn(
                          'ml-auto h-4 w-4',
                          value === option.value ? 'opacity-100' : 'opacity-0'
                        )}
                        aria-hidden="true"
                      />
                    </CommandItem>
                  );
                })}
              </div>
            </CommandGroup>
          </CommandList>
        </Command>
      </PopoverContent>
    </Popover>
  );
};
