import React from 'react';

import AutoSizer from 'react-virtualized-auto-sizer';

import { Dashboard } from './dashboard';

export const Transfer: React.FunctionComponent = () => {
  return (
    <>
      <Dashboard />

      <div style={{ flex: '1 1 auto' }}>
        <AutoSizer disableWidth>
          {({ height }) => (
            <div className="flex flex-1 flex-col ">
              <div
                className={`overflow-y-auto`}
                style={{ height: `${height}px` }}
              ></div>
            </div>
          )}
        </AutoSizer>
      </div>
    </>
  );
};
