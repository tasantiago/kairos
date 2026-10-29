import React from 'react';
import { Slot } from '@radix-ui/react-slot';
import { cva } from 'class-variance-authority';
import { clsx } from 'clsx';
import { twMerge } from 'tailwind-merge';
const variants = cva('button', { variants: { variant: { default: 'primary', outline: 'outline', ghost: 'ghost' }, size: { default: '', icon: 'icon-button', sm: 'small' } }, defaultVariants: { variant: 'default', size: 'default' } });
export const Button = React.forwardRef(({ className, variant, size, asChild = false, ...props }, ref) => {
  const Component = asChild ? Slot : 'button';
  return <Component ref={ref} className={twMerge(clsx(variants({ variant, size }), className))} {...props}/>;
});
Button.displayName = 'Button';
