import React from 'react';

import AutoSizer from 'react-virtualized-auto-sizer';

// import { Panels } from '~/shared/panels/panels';
// import { Origin } from './origin';
// import { Destination } from './destination';
import { Log } from '~/shared/log/log';

export const Plan: React.FunctionComponent = () => {
  return (
    <div style={{ flex: '1 1 auto' }}>
      <AutoSizer disableWidth>
        {({ height }) => (
          <>
            {/* <Panels
              type="2col"
              left={<Origin height={height} />}
              middle={<Destination height={height} />}
            /> */}
            <Log height={height} />
          </>
        )}
      </AutoSizer>
    </div>
  );
};
