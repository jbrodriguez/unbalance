import React from 'react';

import AutoSizer from 'react-virtualized-auto-sizer';

import { useGatherSelected } from '~/state/gather';

export const Shares: React.FunctionComponent = () => {
  const selected = useGatherSelected();

  return (
    <div className="h-full flex flex-col bg-neutral-100 dark:bg-gray-950">
      <div className="flex flex-col p-2">
        <h1 className="text-lg text-slate-500 dark:text-gray-500 pb-2">
          Shares
        </h1>
        <hr className="border-slate-300 dark:border-gray-700" />
      </div>
      <div className="flex-auto">
        <AutoSizer disableWidth>
          {({ height }) => (
            <div
              className="p-2 overflow-y-auto text-slate-500 dark:text-gray-500"
              style={{ height: `${height}px` }}
            >
              {Object.values(selected).map((share) => (
                <div className="whitespace-nowrap">{share}</div>
              ))}
            </div>
          )}
        </AutoSizer>
      </div>
    </div>
  );
};
