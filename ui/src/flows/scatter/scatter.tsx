import React from 'react';

import AutoSizer from 'react-virtualized-auto-sizer';

import { Panels } from '~/shared/panels/panels';
import { Disks } from './disks';
import { FileSystem } from './filesystem';
import { Targets } from './targets';
import { Navbar } from './navbar';
import { Pane } from './pane';
import { Ticker } from './ticker';

export const Scatter: React.FC = () => {
  return (
    <div className="flex flex-col h-full">
      <Navbar />
      <Pane>
        <Ticker />
      </Pane>
      <div style={{ flex: '1 1 auto' }}>
        <AutoSizer disableWidth>
          {({ height }) => (
            <Panels
              type="3col"
              left={<Disks height={height} />}
              middle={<FileSystem height={height} />}
              right={<Targets height={height} />}
            />
          )}
        </AutoSizer>
      </div>
    </div>
  );
};
