import { Variant } from '~/types';

export const styles = {
  enabled: {
    variant: 'accent' as Variant,
    fill: 'fill-neutral-100',
  },
  disabled: {
    variant: 'secondary' as Variant,
    fill: 'fill-neutral-600',
  },
};

export const getStyles = (enabled: boolean) =>
  enabled ? 'enabled' : 'disabled';

export const getVariant = (enabled: boolean) =>
  styles[getStyles(enabled)].variant;

export const getFill = (enabled: boolean) => styles[getStyles(enabled)].fill;
