import React from 'react';

import AutoSizer from 'react-virtualized-auto-sizer';

import { Panels } from '~/shared/panels/panels';
import { Disks } from './disks';
import { FileSystem } from './filesystem';
import { Targets } from './targets';

export const Select: React.FC = () => {
  return (
    <div style={{ flex: '1 1 auto' }}>
      <AutoSizer disableWidth>
        {({ height }) => (
          <>
            <Panels
              type="3col"
              left={<Disks height={height} />}
              middle={<FileSystem height={height} />}
              right={<Targets height={height} />}
            />
          </>
        )}
      </AutoSizer>
    </div>
  );
};
