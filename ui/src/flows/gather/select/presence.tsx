import React from 'react';

import { Panel } from '~/shared/panel/panel';
import { useGatherSelected, useGatherLocation } from '~/state/gather';
import { Icon } from '~/shared/icons/icon';

export const Presence: React.FunctionComponent = () => {
  const selected = useGatherSelected();
  const location = useGatherLocation();

  return (
    <Panel title="Presence">
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
    </Panel>
  );
};
