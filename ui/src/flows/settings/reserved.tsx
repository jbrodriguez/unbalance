import React from 'react';

import { Input } from '@/components/ui/input';
import { Button } from '@/components/ui/button';
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuLabel,
  DropdownMenuRadioGroup,
  DropdownMenuRadioItem,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from '@/components/ui/dropdown-menu';
import { useToast } from '@/components/ui/use-toast';

import { useConfigActions, useConfigReserved } from '~/state/config';

export const Reserved: React.FunctionComponent = () => {
  const { amount, unit } = useConfigReserved();
  const [amountValue, setAmountValue] = React.useState(amount);
  const [unitValue, setUnitValue] = React.useState(unit);
  const { toast } = useToast();
  const { setReservedSpace } = useConfigActions();

  const onAmountChange = (e: React.ChangeEvent<HTMLInputElement>) =>
    setAmountValue(+e.target.value);

  const onApply = () => {
    if (!Number.isInteger(amountValue) || amountValue <= 0) {
      toast({
        title: 'Amount value must be a positive integer',
        variant: 'destructive',
      });
      return;
    }

    if (unitValue === 'Gb' && amountValue < 1) {
      toast({
        title: 'Gb value must be greater or equal to 1',
        variant: 'destructive',
      });
      return;
    }

    // if we get here, we can save the value
    setReservedSpace(amountValue, unitValue);
  };

  return (
    <div className="p-4">
      <h1>
        unbalanced uses the threshold defined here as the minimum free space
        that should be kept available in a target disk, when planning how much
        the disk can be filled. <br />
        This threshold cannot be less than 1 Gb (hard limit set by this app).
      </h1>
      <div className="pb-4" />

      <div className="flex flex-row items-center">
        <Input
          type="number"
          placeholder="Size"
          className="w-40"
          defaultValue={amount}
          value={amountValue}
          onChange={onAmountChange}
        />
        <div className="pr-4" />
        <DropdownMenu>
          <DropdownMenuTrigger asChild>
            <Button>{unitValue}</Button>
          </DropdownMenuTrigger>
          <DropdownMenuContent className="">
            <DropdownMenuLabel>Unit</DropdownMenuLabel>
            <DropdownMenuSeparator />
            <DropdownMenuRadioGroup
              value={unitValue}
              onValueChange={setUnitValue}
            >
              <DropdownMenuRadioItem value="%">%</DropdownMenuRadioItem>
              <DropdownMenuRadioItem value="Gb">Gb</DropdownMenuRadioItem>
            </DropdownMenuRadioGroup>
          </DropdownMenuContent>
        </DropdownMenu>
      </div>
      <div className="pb-4" />
      <Button onClick={onApply} variant="secondary">
        Apply
      </Button>
    </div>
  );
};
