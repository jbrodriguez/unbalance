import React from 'react';

import AutoSizer from 'react-virtualized-auto-sizer';

import { useGatherSelected, useGatherLocation } from '~/state/gather';
import { Icon } from '~/shared/icons/icon';

export const Presence: React.FunctionComponent = () => {
  const selected = useGatherSelected();
  const location = useGatherLocation();

  return (
    <div className="h-full flex flex-col bg-neutral-100 dark:bg-gray-950">
      <div className="flex flex-col p-2">
        <h1 className="text-lg text-slate-500 dark:text-gray-500 pb-2">
          Presence
        </h1>
        <hr className="border-slate-300 dark:border-gray-700" />
      </div>
      <div className="flex-auto">
        <AutoSizer disableWidth>
          {({ height }) => (
            <div className="overflow-y-auto" style={{ height: `${height}px` }}>
              {Object.keys(selected).map((key) => {
                return (
                  <div key={key} className="flex flex-row items-center p-2">
                    <div className="flex flex-col flex-1">
                      <div className="flex flex-row items-center">
                        <Icon
                          name="file"
                          size={20}
                          style="fill-blue-400 dark:fill-gray-700"
                        />
                        <span className="pr-2" />
                        <div className="font-bold">{selected[key]}</div>
                      </div>
                      <div className="pl-7 text-neutral-500 dark:text-gray-500">
                        {location[key].join(', ')}
                      </div>
                    </div>
                  </div>
                );
              })}
            </div>
          )}
        </AutoSizer>
      </div>
    </div>
  );
};
