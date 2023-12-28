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
    <AutoSizer disableWidth>
      {({ height }) => (
        <div className="flex flex-1 flex-col bg-neutral-100 dark:bg-gray-950">
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
  );
};
