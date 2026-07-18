import React from 'react';

import { Panel } from '~/shared/panel/panel';
import {
  useGatherSelected,
  useGatherLocation,
  useGatherSizes,
  useGatherActions,
} from '~/state/gather';
import { Icon } from '~/shared/icons/icon';
import { humanBytes } from '~/helpers/units';

export const Presence: React.FunctionComponent = () => {
  const selected = useGatherSelected();
  const location = useGatherLocation();
  const sizes = useGatherSizes();
  const { loadSize } = useGatherActions();

  // sizes are computed in the background, selection never waits on them
  React.useEffect(() => {
    Object.keys(selected).forEach((key) => loadSize(key, selected[key]));
  }, [selected, loadSize]);

  const keys = Object.keys(selected);
  const pending = keys.some((key) => !sizes[key]);
  const total = keys.reduce((sum, key) => sum + (sizes[key]?.total ?? 0), 0);

  const subtitle =
    keys.length > 0 ? (
      <div className="flex flex-row items-center gap-1 text-slate-500 dark:text-gray-500">
        <span>{humanBytes(total)}</span>
        {pending && (
          <Icon
            name="loading"
            size={16}
            style="animate-spin fill-slate-500 dark:fill-gray-500"
          />
        )}
      </div>
    ) : undefined;

  return (
    <Panel title="Presence" subtitle={subtitle}>
      {keys.map((key) => {
        const size = sizes[key];

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
                <div className="flex-1" />
                {size ? (
                  <div>{humanBytes(size.total)}</div>
                ) : (
                  <Icon
                    name="loading"
                    size={16}
                    style="animate-spin fill-slate-500 dark:fill-gray-500"
                  />
                )}
              </div>
              <div className="pl-7 text-neutral-500 dark:text-gray-500">
                {location[key]
                  .map((disk) =>
                    size && size.disks[disk] !== undefined
                      ? `${disk} (${humanBytes(size.disks[disk])})`
                      : disk,
                  )
                  .join(', ')}
              </div>
            </div>
          </div>
        );
      })}
    </Panel>
  );
};
