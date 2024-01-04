import React from 'react';

import AutoSizer from 'react-virtualized-auto-sizer';

import { useUnraidActions, useUnraidLogs } from '~/state/unraid';

export const Logs: React.FunctionComponent = () => {
  const { getLog } = useUnraidActions();
  const logs = useUnraidLogs();

  React.useEffect(() => {
    getLog();
  }, [getLog]);

  return (
    <div className="h-full flex flex-col bg-neutral-100 dark:bg-gray-950">
      <div className="flex-auto">
        <AutoSizer disableWidth>
          {({ height }) => (
            <div
              className="overflow-y-auto p-4 text-slate-700 dark:text-gray-300"
              style={{ height: `${height}px` }}
            >
              {logs.map((log) => (
                <p>{log}</p>
              ))}
            </div>
          )}
        </AutoSizer>
      </div>
    </div>
  );
};
