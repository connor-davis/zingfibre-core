import {
  type ChangeEvent,
  type ComponentProps,
  useEffect,
  useState,
} from 'react';

import { cn } from '@/lib/utils';

type DebouncedInputProps = {
  delay?: number;
};

function DebounceInput({
  className,
  delay = 500,
  ...props
}: ComponentProps<'input'> & DebouncedInputProps) {
  const [localChange, setLocalChange] =
    useState<ChangeEvent<HTMLInputElement>>();

  useEffect(() => {
    const handler = setTimeout(() => {
      if (!props.onChange || !localChange) return;

      props.onChange(localChange);
      setLocalChange(undefined);
    }, delay);

    return () => {
      clearTimeout(handler);
    };
  }, [localChange, delay, props.onChange]);

  return (
    <input
      data-slot="input"
      className={cn(
        'file:text-foreground placeholder:text-muted-foreground selection:bg-primary selection:text-primary-foreground border-input flex h-9 w-full min-w-0 rounded-md border bg-input/80 dark:bg-input/30 px-3 py-1 text-base transition-[color,box-shadow] outline-none file:inline-flex file:h-7 file:border-0 file:bg-transparent file:text-sm file:font-medium disabled:pointer-events-none disabled:opacity-50 disabled:border-muted dark:disabled:border-input disabled:cursor-not-allowed md:text-sm',
        'focus-visible:border-ring focus-visible:ring-ring/50 focus-visible:ring-[3px]',
        'aria-invalid:ring-destructive/20 dark:aria-invalid:ring-destructive/40 aria-invalid:border-destructive',
        className
      )}
      {...props}
      onChange={(event) => setLocalChange(event)}
    />
  );
}

export { DebounceInput };
