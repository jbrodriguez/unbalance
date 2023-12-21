import React from 'react';

// import AutoSizer from 'react-virtualized-auto-sizer';

interface BinProps {
  height?: number;
}

export const Bin: React.FunctionComponent<BinProps> = ({ height }) => {
  // const bin = useUnraidBin();
  return (
    // <div style={{ flex: '1 1 auto' }}>
    //   <AutoSizer disableWidth>
    //     {({ height }) => (
    <div className="flex flex-1 flex-col bg-neutral-200 dark:bg-gray-950">
      <div className="p-2 overflow-y-auto" style={{ height: `${height}px` }}>
        <h1>bins</h1>
      </div>
    </div>
    //     )}
    //   </AutoSizer>
    // </div>
  );
};
