import React from 'react';

import AutoSizer from 'react-virtualized-auto-sizer';

import { useUnraidDisks } from '~/state/unraid';
import { Disk as IDisk } from '~/types';
import { useScatterActions, useScatterSource } from '~/state/scatter';
import { Selectable } from '~/shared/selectable/selectable';
import { Disk } from '~/shared/disk/disk';

export const Disks: React.FunctionComponent = () => {
  const disks = useUnraidDisks();
  const selected = useScatterSource();
  const { setSource } = useScatterActions();

  const onDiskClick = (disk: IDisk) => () => setSource(disk.name);

  return (
    <div className="h-full flex flex-col bg-neutral-100 dark:bg-gray-950">
      <div className="flex flex-col p-2">
        <h1 className="text-lg text-slate-500 dark:text-gray-500 pb-2">
          Source Disk
        </h1>
        <hr className="border-slate-300 dark:border-gray-700" />
      </div>
      <div className="flex-auto">
        <AutoSizer disableWidth>
          {({ height }) => (
            <div
              className="overflow-y-auto p-2"
              style={{ height: `${height}px` }}
            >
              {disks.map((disk) => (
                <Selectable
                  key={disk.id}
                  onClick={onDiskClick(disk)}
                  selected={disk.name === selected}
                >
                  <Disk disk={disk} />
                </Selectable>
              ))}
            </div>
          )}
        </AutoSizer>
      </div>
    </div>
  );
};
