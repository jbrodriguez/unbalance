import React from 'react';

import { useUnraidStep } from '~/state/unraid';
import { Description as SelectDescription } from './select/description';
import { Description as PlanDescription } from './plan/description';
import { Description as TransferDescription } from './transfer/description';

export const Ticker: React.FunctionComponent = () => {
  const step = useUnraidStep();

  return (
    <>
      {(step === 'select' || step === 'idle') && <SelectDescription />}
      {step === 'plan' && <PlanDescription />}
      {step === 'transfer' && <TransferDescription />}
    </>
  );
};
