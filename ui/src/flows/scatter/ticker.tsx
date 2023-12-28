import React from 'react';

import { useUnraidRoute } from '~/state/unraid';
import { Description as SelectDescription } from './select/description';
import { Description as PlanDescription } from './plan/description';
import { Description as TransferDescription } from './transfer/description';
import { Line } from '~/shared/transfer/line';

export const Ticker: React.FunctionComponent = () => {
  const route = useUnraidRoute();

  return (
    <>
      {(route === '/scatter/select' || route === '/scatter') && (
        <SelectDescription />
      )}
      {route === '/scatter/plan' && <PlanDescription />}
      {route === '/scatter/transfer/validation' && <TransferDescription />}
      {route === '/scatter/transfer/operation' && <Line />}
    </>
  );
};
