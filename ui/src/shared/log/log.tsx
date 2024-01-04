import React from 'react';

import AutoSizer from 'react-virtualized-auto-sizer';

import { useUnraidLogs } from '~/state/unraid';

export const Log: React.FunctionComponent = () => {
  const logs = useUnraidLogs();

  return (
    <AutoSizer disableWidth>
      {({ height }) => (
        <div className="flex flex-1 flex-col bg-neutral-100 dark:bg-gray-950">
          <div
            className="overflow-y-auto p-2 text-base text-gray-700 dark:text-gray-500"
            style={{ height: `${height}px` }}
          >
            {logs.map((log) => (
              <p>{log}</p>
            ))}
          </div>
        </div>
      )}
    </AutoSizer>
  );
};
