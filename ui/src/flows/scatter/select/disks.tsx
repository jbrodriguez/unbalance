import React from 'react';

import AutoSizer from 'react-virtualized-auto-sizer';

import { useUnraidDisks } from '~/state/unraid';
import { Disk as IDisk } from '~/types';
import { useScatterActions, useScatterSource } from '~/state/scatter';
import { Selectable } from '~/shared/disk/selectable-disk';
import { Disk } from '~/shared/disk/base-disk';

export const Disks: React.FunctionComponent = () => {
  const disks = useUnraidDisks();
  const selected = useScatterSource();
  const { setSource } = useScatterActions();

  const onDiskClick = (disk: IDisk) => setSource(disk.name);

  return (
    <div className="h-full bg-neutral-100 dark:bg-gray-950">
      <div className="flex flex-col p-2">
        <span className="text-lg font-bold text-slate-500 dark:text-gray-500">
          Source Disk
        </span>
        <span className="border-b pt-2 px-2 border-slate-300 dark:border-gray-700" />
      </div>
      <AutoSizer disableWidth>
        {({ height }) => (
          <div className="flex flex-1 flex-col">
            <div
              className="overflow-y-auto px-2 pt-2"
              style={{ height: `${height}px` }}
            >
              {disks.map((disk) => (
                <Selectable
                  disk={disk}
                  onSelectDisk={onDiskClick}
                  selected={disk.name === selected}
                >
                  <Disk disk={disk} />
                </Selectable>
              ))}
            </div>
          </div>
        )}
      </AutoSizer>
    </div>
  );
};
