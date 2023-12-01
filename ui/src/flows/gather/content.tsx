import React from 'react';

import AutoSizer from 'react-virtualized-auto-sizer';

import { Panels } from '~/shared/panels/panels';
import { Shares } from './select/shares';
import { Presence } from './select/presence';
import { Targets } from './plan/targets';
import { Transfer } from '~/shared/transfer/transfer';

export const Content: React.FunctionComponent = () => {
  const step: string = 'transfer';

  return (
    <>
      {step !== 'transfer' && (
        <div style={{ flex: '1 1 auto' }}>
          <AutoSizer disableWidth>
            {({ height }) => (
              <>
                {step === 'select' && (
                  <Panels
                    type="2col"
                    left={<Shares height={height} />}
                    middle={<Presence height={height} />}
                  />
                )}
                {step === 'plan' && <Targets height={height} />}
              </>
            )}
          </AutoSizer>
        </div>
      )}
      {step === 'transfer' && <Transfer />}
    </>
  );
};
