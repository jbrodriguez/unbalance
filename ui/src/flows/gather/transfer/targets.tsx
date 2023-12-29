import React from 'react';

import AutoSizer from 'react-virtualized-auto-sizer';

import { Panels } from '~/shared/panels/panels';
import { Shares } from './shares';
import { Target } from './target';
import { Bin } from '~/shared/bin/bin';
import { useGatherTarget } from '~/state/gather';

export const Targets: React.FunctionComponent = () => {
  const disk = useGatherTarget();

  return (
    <div style={{ flex: '1 1 auto' }}>
      <AutoSizer disableWidth>
        {({ height }) => (
          <>
            <Panels
              type="3col"
              left={<Shares height={height} />}
              middle={<Target height={height} />}
              right={<Bin disk={disk} />}
              dimensions={{ left: 20, middle: 50, right: 30 }}
            />
          </>
        )}
      </AutoSizer>
    </div>
  );
};
