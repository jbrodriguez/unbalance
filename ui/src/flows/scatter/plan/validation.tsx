import React from 'react';

import AutoSizer from 'react-virtualized-auto-sizer';

import { Panels } from '~/shared/panels/panels';
import { Origin } from './origin';
import { Destination } from './destination';

export const Validation: React.FunctionComponent = () => {
  return (
    <div style={{ flex: '1 1 auto' }}>
      <AutoSizer disableWidth>
        {({ height }) => (
          <>
            <Panels
              type="2col"
              left={<Origin height={height} />}
              middle={<Destination height={height} />}
            />
          </>
        )}
      </AutoSizer>
    </div>
  );
};
