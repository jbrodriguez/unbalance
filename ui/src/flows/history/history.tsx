import React from 'react';

import AutoSizer from 'react-virtualized-auto-sizer';

import { Panels } from '~/shared/panels/panels';
import { Operations } from './operations';
import { Operation } from './operation';

export const History: React.FunctionComponent = () => {
  return (
    <div className="flex flex-col h-full">
      <div style={{ flex: '1 1 auto' }}>
        <AutoSizer disableWidth>
          {({ height }) => (
            <Panels
              type="2col"
              left={<Operations height={height} />}
              middle={<Operation height={height} />}
            />
          )}
        </AutoSizer>
      </div>
    </div>
  );
};
