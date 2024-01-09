import React from 'react';

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

import { useUnraidActions, useUnraidError } from '~/state/unraid';
import { useConfigActions, useConfigRefreshRate } from '~/state/config';

const refreshDiscreteToSec: Record<number, string> = {
  1000: '1 sec',
  5000: '5 sec',
  15000: '15 sec',
  30000: '30 sec',
};

export const Actions: React.FunctionComponent = () => {
  const { setRefreshRate } = useConfigActions();
  const { stop } = useUnraidActions();
  const refresh = useConfigRefreshRate();
  const error = useUnraidError();
  const { toast } = useToast();

  const [pressed, setPressed] = React.useState(false);
  React.useEffect(() => {
    if (error !== '') {
      toast({
        title: 'Operation Error',
        description: error,
        variant: 'destructive',
      });
    }
  }, [toast, error]);

  const display = refreshDiscreteToSec[refresh];

  const onStop = () => {
    setPressed(true);
    stop();
  };
  const onChangeRefreshRate = (value: string) => setRefreshRate(+value);

  return (
    <div className="flex flex-row items-center justify-end pr-2">
      {/* <Button label="STOP" variant="primary" onClick={onMove} /> */}
      <Button variant="secondary" onClick={onStop} disabled={pressed}>
        STOP
      </Button>
      <span className="px-1">|</span>
      <span>Refresh Rate: </span>
      <span className="pr-2" />
      <DropdownMenu>
        <DropdownMenuTrigger asChild>
          <Button>{display}</Button>
        </DropdownMenuTrigger>
        <DropdownMenuContent className="">
          <DropdownMenuLabel>Unit</DropdownMenuLabel>
          <DropdownMenuSeparator />
          <DropdownMenuRadioGroup
            value={refresh.toString()}
            onValueChange={onChangeRefreshRate}
          >
            <DropdownMenuRadioItem value="1000">1 sec</DropdownMenuRadioItem>
            <DropdownMenuRadioItem value="5000">5 sec</DropdownMenuRadioItem>
            <DropdownMenuRadioItem value="15000">15 sec</DropdownMenuRadioItem>
            <DropdownMenuRadioItem value="30000">30 sec</DropdownMenuRadioItem>
          </DropdownMenuRadioGroup>
        </DropdownMenuContent>
      </DropdownMenu>
    </div>
  );
};
