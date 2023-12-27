import React from 'react';

import AutoSizer from 'react-virtualized-auto-sizer';

import { useScatterBinDisk } from '~/state/scatter';
import { Panels } from '~/shared/panels/panels';
import { Origin } from './origin';
import { Destination } from './destination';
import { Bin } from '~/shared/bin/bin';

export const Validation: React.FunctionComponent = () => {
  const disk = useScatterBinDisk();

  return (
    <div style={{ flex: '1 1 auto' }}>
      <AutoSizer disableWidth>
        {({ height }) => (
          <>
            <Panels
              type="3col"
              left={<Origin height={height} />}
              middle={<Destination height={height} />}
              right={<Bin height={height} disk={disk} />}
              dimensions={{ left: 30, middle: 30 }}
            />
          </>
        )}
      </AutoSizer>
    </div>
  );
};
