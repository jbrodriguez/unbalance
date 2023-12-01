import React from 'react';

import AutoSizer from 'react-virtualized-auto-sizer';

import { Panels } from '~/shared/panels/panels';
import { Disks } from './select/disks';
import { FileSystem } from './select/filesystem';
import { Targets } from './select/targets';
import { Origin } from './plan/origin';
import { Destination } from './plan/destination';

export const Content: React.FunctionComponent = () => {
  const step: string = 'plan';

  return (
    <div style={{ flex: '1 1 auto' }}>
      <AutoSizer disableWidth>
        {({ height }) => (
          <>
            {step === 'select' && (
              <Panels
                type="3col"
                left={<Disks height={height} />}
                middle={<FileSystem height={height} />}
                right={<Targets height={height} />}
              />
            )}
            {step === 'plan' && (
              <Panels
                type="2col"
                left={<Origin height={height} />}
                middle={<Destination height={height} />}
              />
            )}
          </>
        )}
      </AutoSizer>
    </div>
  );
};
