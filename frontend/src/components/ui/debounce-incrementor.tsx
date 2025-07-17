import { MinusIcon, PlusIcon } from 'lucide-react';
import {
  type ChangeEvent,
  forwardRef,
  useEffect,
  useImperativeHandle,
  useRef,
  useState,
} from 'react';

import { cn } from '@/lib/utils';

import { Button } from './button';
import { Input } from './input';

export interface InputProps
  extends React.InputHTMLAttributes<HTMLInputElement> {
  symbol?: string;
  onValueChange?: (value: number) => void;
}

type DebouncedInputProps = {
  delay?: number;
};

const DebounceIncrementorInput = forwardRef<
  HTMLInputElement,
  InputProps & DebouncedInputProps
>(({ symbol, className, delay = 500, ...props }, ref) => {
  const [hitMax, setHitMax] = useState(false);
  const [hitMin, setHitMin] = useState(false);
  const incrementInput = useRef<HTMLInputElement>(null);

  const [localChange, setLocalChange] =
    useState<ChangeEvent<HTMLInputElement>>();

  useEffect(() => {
    const handler = setTimeout(() => {
      if (!props.onChange || !localChange || !props.onValueChange) return;

      props.onChange(localChange);

      const value = Number(incrementInput.current?.value);

      if (!isNaN(value)) {
        props.onValueChange(value);
      }

      setLocalChange(undefined);
    }, delay);

    return () => {
      clearTimeout(handler);
    };
  }, [localChange, delay, props.onChange, props.onValueChange]);

  useImperativeHandle(ref, () => incrementInput.current!, []);

  const increment = () => {
    incrementInput.current?.stepUp();
    // Supports onchange events
    incrementInput.current?.dispatchEvent(
      new Event('change', { bubbles: true })
    );
    // Disbale when hitting max
    setHitMax(incrementInput.current?.value === incrementInput.current?.max);
    setHitMin(incrementInput.current?.value === incrementInput.current?.min);
  };

  const decrement = () => {
    incrementInput.current?.stepDown();
    // Supports onchange events
    incrementInput.current?.dispatchEvent(
      new Event('change', { bubbles: true })
    );
    // Disbale when hitting min
    setHitMax(incrementInput.current?.value === incrementInput.current?.max);
    setHitMin(incrementInput.current?.value === incrementInput.current?.min);
  };

  return (
    <div className="flex items-center w-full">
      <Button
        variant="outline"
        type="button"
        size="icon"
        disabled={hitMax}
        onClick={increment}
        aria-label="increase"
        className="border-r-0 rounded-r-none"
      >
        <PlusIcon className="size-4" />
      </Button>

      <div className="relative w-full">
        <Input
          type="number"
          className={cn(
            'border-x-0 w-full rounded-none text-center [appearance:textfield] [&::-webkit-outer-spin-button]:appearance-none [&::-webkit-inner-spin-button]:appearance-none',
            className
          )}
          ref={incrementInput}
          {...props}
          onChange={(event) => setLocalChange(event)}
        />
        {symbol && <span className="absolute right-4 top-0">{symbol}</span>}
      </div>

      <Button
        variant="outline"
        type="button"
        size="icon"
        disabled={hitMin}
        onClick={decrement}
        aria-label="decrease"
        className="border-l-0 rounded-l-none"
      >
        <MinusIcon className="size-4" />
      </Button>
    </div>
  );
});
DebounceIncrementorInput.displayName = 'DebounceIncrementorInput';

export { DebounceIncrementorInput };
