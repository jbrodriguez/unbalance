import React from 'react';

import { useUnraidDisks } from '~/state/unraid';
import {
  useScatterSource,
  useScatterSelected,
  useScatterTargets,
  useScatterActions,
} from '~/state/scatter';
import { Disk } from '~/shared/disk/disk';
import { Disk as IDisk, Targets as ITargets } from '~/types';

interface Props {
  height?: number;
}

const isChecked = (name: string, targets: ITargets) => targets[name] || false;

export const Targets: React.FC<Props> = ({ height = 0 }) => {
  const disks = useUnraidDisks();
  const source = useScatterSource();
  const selected = useScatterSelected();
  const targets = useScatterTargets();
  const { toggleTarget } = useScatterActions();

  const visible = source !== '' && selected.length > 0;
  const elegible = disks.filter((disk) => disk.name !== source);

  const onCheck = (disk: IDisk) => () => toggleTarget(disk.name);

  return (
    <div className="flex flex-1 bg-neutral-200 dark:bg-gray-950">
      <div
        className="flex flex-1 flex-col overflow-y-auto px-2 pt-2"
        style={{ height: `${height}px` }}
      >
        {visible &&
          elegible.map((disk) => (
            <Disk
              disk={disk}
              checkable
              checked={isChecked(disk.name, targets)}
              onCheck={onCheck(disk)}
            />
          ))}
      </div>
    </div>
  );
};
