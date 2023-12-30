import React from 'react';

import { useGatherSelected, useGatherLocation } from '~/state/gather';
import { Icon } from '~/shared/icons/icon';

interface Props {
  height?: number;
}

export const Presence: React.FunctionComponent<Props> = ({ height = 0 }) => {
  const selected = useGatherSelected();
  const location = useGatherLocation();

  return (
    <div className="flex flex-1 flex-col bg-neutral-200 dark:bg-gray-950">
      <div className={`overflow-y-auto`} style={{ height: `${height}px` }}>
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
                  <div className="font-bold text-slate-700 dark:text-slate-200">
                    {selected[key]}
                  </div>
                </div>
                <div className="pl-7 text-neutral-500 dark:text-gray-500">
                  {location[key].join(', ')}
                </div>
              </div>
            </div>
          );
        })}
      </div>
    </div>
  );
};
