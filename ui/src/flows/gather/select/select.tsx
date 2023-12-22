import React from 'react';

import AutoSizer from 'react-virtualized-auto-sizer';

import { Panels } from '~/shared/panels/panels';
import { Shares } from './shares';
import { Presence } from './presence';

export const Select: React.FunctionComponent = () => {
  return (
    <div style={{ flex: '1 1 auto' }}>
      <AutoSizer disableWidth>
        {({ height }) => (
          <>
            <Panels
              type="2col"
              left={<Shares height={height} />}
              middle={<Presence height={height} />}
            />
          </>
        )}
      </AutoSizer>
    </div>
  );
};
