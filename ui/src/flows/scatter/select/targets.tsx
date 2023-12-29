import React from 'react';

import AutoSizer from 'react-virtualized-auto-sizer';

import { useUnraidDisks } from '~/state/unraid';
import {
  useScatterSource,
  useScatterSelected,
  useScatterTargets,
  useScatterActions,
} from '~/state/scatter';
import { Checkbox } from '~/shared/checkbox/checkbox';
import { Disk } from '~/shared/disk/base-disk';
import { Disk as IDisk, Targets as ITargets } from '~/types';

const isChecked = (name: string, targets: ITargets) => targets[name] || false;

export const Targets: React.FunctionComponent = () => {
  const disks = useUnraidDisks();
  const source = useScatterSource();
  const selected = useScatterSelected();
  const targets = useScatterTargets();
  const { toggleTarget } = useScatterActions();

  const visible = source !== '' && selected.length > 0;
  const elegible = disks.filter((disk) => disk.name !== source);

  const onCheck = (disk: IDisk) => () => toggleTarget(disk.name);

  return (
    <div className="h-full bg-neutral-100 dark:bg-gray-950">
      <div className="flex flex-col pt-2 px-2">
        <h1 className="text-lg text-slate-500 dark:text-gray-500 pb-2">
          Target Disk(s)
        </h1>
        <hr className="border-slate-300 dark:border-gray-700" />
      </div>
      <AutoSizer disableWidth>
        {({ height }) => (
          <div className="flex flex-1 bg-neutral-100 dark:bg-gray-950">
            <div
              className="flex flex-1 flex-col overflow-y-auto px-2 pt-2"
              style={{ height: `${height}px` }}
            >
              {visible &&
                elegible.map((disk) => (
                  <div className="flex flex-row items-center">
                    <Checkbox
                      checked={isChecked(disk.name, targets)}
                      onCheck={onCheck(disk)}
                    />
                    <div className="pr-4" />
                    <Disk disk={disk} />
                  </div>
                ))}
            </div>
          </div>
        )}
      </AutoSizer>
    </div>
  );
};
