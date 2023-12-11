import React from 'react';

import AutoSizer from 'react-virtualized-auto-sizer';

import { Targets } from './targets';

export const Plan: React.FunctionComponent = () => {
  return (
    <div style={{ flex: '1 1 auto' }}>
      <AutoSizer disableWidth>
        {({ height }) => <Targets height={height} />}
      </AutoSizer>
    </div>
  );
};