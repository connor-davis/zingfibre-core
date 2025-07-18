import { ChevronDown, ChevronUp } from 'lucide-react';
import { forwardRef, useCallback, useEffect, useRef, useState } from 'react';
import { NumericFormat, type NumericFormatProps } from 'react-number-format';

import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';

export interface DebounceNumberInputProps
  extends Omit<NumericFormatProps, 'value' | 'onValueChange'> {
  stepper?: number;
  thousandSeparator?: string;
  placeholder?: string;
  defaultValue?: number;
  min?: number;
  max?: number;
  value?: number; // Controlled value
  suffix?: string;
  prefix?: string;
  onValueChange?: (value: number | undefined) => void;
  fixedDecimalScale?: boolean;
  decimalScale?: number;
  debounceMs?: number; // Debounce delay in ms
}

export const DebounceNumberInput = forwardRef<
  HTMLInputElement,
  DebounceNumberInputProps
>(
  (
    {
      stepper,
      thousandSeparator,
      placeholder,
      defaultValue,
      min = -Infinity,
      max = Infinity,
      onValueChange,
      fixedDecimalScale = false,
      decimalScale = 0,
      suffix,
      prefix,
      value: controlledValue,
      debounceMs = 300,
      ...props
    },
    ref
  ) => {
    const [value, setValue] = useState<number | undefined>(
      controlledValue ?? defaultValue
    );
    const [debouncedValue, setDebouncedValue] = useState<number | undefined>(
      value
    );
    const debounceTimeout = useRef<ReturnType<typeof setTimeout>>(undefined);

    // Debounce effect
    useEffect(() => {
      if (debounceTimeout.current) clearTimeout(debounceTimeout.current);
      debounceTimeout.current = setTimeout(() => {
        setDebouncedValue(value);
      }, debounceMs);
      return () => {
        if (debounceTimeout.current) clearTimeout(debounceTimeout.current);
      };
    }, [value, debounceMs]);

    // Call onValueChange only when debouncedValue changes
    useEffect(() => {
      if (onValueChange) {
        onValueChange(debouncedValue);
      }
    }, [debouncedValue]);

    const handleIncrement = useCallback(() => {
      setValue((prev) =>
        prev === undefined
          ? (stepper ?? 1)
          : Math.min(prev + (stepper ?? 1), max)
      );
    }, [stepper, max]);

    const handleDecrement = useCallback(() => {
      setValue((prev) =>
        prev === undefined
          ? -(stepper ?? 1)
          : Math.max(prev - (stepper ?? 1), min)
      );
    }, [stepper, min]);

    useEffect(() => {
      const handleKeyDown = (e: KeyboardEvent) => {
        if (
          document.activeElement ===
          (ref as React.RefObject<HTMLInputElement>).current
        ) {
          if (e.key === 'ArrowUp') {
            handleIncrement();
          } else if (e.key === 'ArrowDown') {
            handleDecrement();
          }
        }
      };

      window.addEventListener('keydown', handleKeyDown);

      return () => {
        window.removeEventListener('keydown', handleKeyDown);
      };
    }, [handleIncrement, handleDecrement, ref]);

    useEffect(() => {
      if (controlledValue !== undefined) {
        setValue(controlledValue);
      }
    }, [controlledValue]);

    const handleChange = (values: {
      value: string;
      floatValue: number | undefined;
    }) => {
      const newValue =
        values.floatValue === undefined ? undefined : values.floatValue;
      setValue(newValue);
    };

    const handleBlur = () => {
      if (value !== undefined) {
        if (value < min) {
          setValue(min);
          (ref as React.RefObject<HTMLInputElement>).current!.value =
            String(min);
        } else if (value > max) {
          setValue(max);
          (ref as React.RefObject<HTMLInputElement>).current!.value =
            String(max);
        }
      }
    };

    return (
      <div className="flex items-center w-full">
        <NumericFormat
          value={value}
          onValueChange={handleChange}
          thousandSeparator={thousandSeparator}
          decimalScale={decimalScale}
          fixedDecimalScale={fixedDecimalScale}
          allowNegative={min < 0}
          valueIsNumericString
          onBlur={handleBlur}
          max={max}
          min={min}
          suffix={suffix}
          prefix={prefix}
          customInput={Input}
          placeholder={placeholder}
          className="[appearance:textfield] [&::-webkit-outer-spin-button]:appearance-none [&::-webkit-inner-spin-button]:appearance-none rounded-r-none w-auto h-8 relative"
          getInputRef={ref}
          {...props}
        />

        <div className="flex flex-col w-auto h-auto">
          <Button
            aria-label="Increase value"
            className="px-2 h-4 rounded-l-none rounded-br-none border-input border-l-0 border-b-[0.5px] focus-visible:relative"
            variant="outline"
            onClick={handleIncrement}
            disabled={value === max}
          >
            <ChevronUp size={15} />
          </Button>
          <Button
            aria-label="Decrease value"
            className="px-2 h-4 rounded-l-none rounded-tr-none border-input border-l-0 border-t-[0.5px] focus-visible:relative"
            variant="outline"
            onClick={handleDecrement}
            disabled={value === min}
          >
            <ChevronDown size={15} />
          </Button>
        </div>
      </div>
    );
  }
);
