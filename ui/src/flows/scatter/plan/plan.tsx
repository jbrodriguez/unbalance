import React from 'react';

import AutoSizer from 'react-virtualized-auto-sizer';
import { Outlet } from 'react-router-dom';

// import { Panels } from '~/shared/panels/panels';
// import { Origin } from './origin';
// import { Destination } from './destination';
// import { Log } from '~/shared/log/log';
import { Sidebar } from '~/flows/scatter/plan/sidebar';

export const Plan: React.FunctionComponent = () => {
  return (
    <div style={{ flex: '1 1 auto' }}>
      <AutoSizer disableWidth>
        {({ height }) => (
          <div className="flex flex-row">
            {/* <Panels
              type="2col"
              left={<Origin height={height} />}
              middle={<Destination height={height} />}
            /> */}
            <Sidebar height={height} />
            <Outlet />
          </div>
        )}
      </AutoSizer>
    </div>
  );
};
