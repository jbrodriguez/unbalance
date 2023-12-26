import React from 'react';

import AutoSizer from 'react-virtualized-auto-sizer';

import { Dashboard } from './dashboard';
import { Commands } from './commands';

export const Transfer: React.FunctionComponent = () => {
  return (
    <>
      <Dashboard />

      <span className="pb-4" />

      <div style={{ flex: '1 1 auto' }}>
        <AutoSizer disableWidth>
          {({ height }) => (
            <div className="flex flex-1 flex-col ">
              <div
                className={`overflow-y-auto`}
                style={{ height: `${height}px` }}
              >
                <Commands />
              </div>
            </div>
          )}
        </AutoSizer>
      </div>
    </>
  );
};
