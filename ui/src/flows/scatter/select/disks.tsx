import React from 'react';

import { useUnraidDisks } from '~/state/unraid';
import { Disk as IDisk } from '~/types';
import { useScatterActions, useScatterSource } from '~/state/scatter';
import { Selectable } from '~/shared/disk/selectable-disk';
import { Disk } from '~/shared/disk/base-disk';

interface Props {
  height?: number;
}

// const selectedBackground = (selected: boolean) =>
//   selected ? 'rounded dark:bg-gray-900 bg-neutral-300' : '';

export const Disks: React.FunctionComponent<Props> = ({ height = 0 }) => {
  const disks = useUnraidDisks();
  const selected = useScatterSource();
  const { setSource } = useScatterActions();

  const onDiskClick = (disk: IDisk) => setSource(disk.name);

  return (
    <div className="flex flex-1 flex-col bg-neutral-200 dark:bg-gray-950">
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
  );
};
