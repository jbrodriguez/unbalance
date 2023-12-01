import React from 'react';

import { Description as SelectDescription } from './select/description';
import { Description as PlanDescription } from './plan/description';
import { Description as TransferDescription } from './transfer/description';

export const Ticker: React.FunctionComponent = () => {
  const step: string = 'plan';

  return (
    <>
      {step === 'select' && <SelectDescription />}
      {step === 'plan' && <PlanDescription />}
      {step === 'transfer' && <TransferDescription />}
    </>
  );
};
