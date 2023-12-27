import React from 'react';

import { useUnraidRoute } from '~/state/unraid';
import { Description as SelectDescription } from './select/description';
import { Description as PlanDescription } from './plan/description';
import { Description as TransferDescription } from './transfer/description';

export const Ticker: React.FunctionComponent = () => {
  const route = useUnraidRoute();

  return (
    <>
      {route === '/gather/select' ||
        (route === '/gather' && <SelectDescription />)}
      {route === '/gather/plan' && <PlanDescription />}
      {route === '/gather/transfer/targets' && <TransferDescription />}
    </>
  );
};
