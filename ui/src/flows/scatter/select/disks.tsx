import React from 'react';

import { useUnraidDisks } from '~/state/unraid';
import { humanBytes } from '~/helpers/units';
import { useScatterActions, useScatterSource } from '~/state/scatter';

interface Props {
  height?: number;
}

const selectedBackground = (selected: boolean) =>
  selected ? 'rounded dark:bg-gray-900 bg-neutral-300' : '';

export const Disks: React.FunctionComponent<Props> = ({ height = 0 }) => {
  const disks = useUnraidDisks();
  const selected = useScatterSource();
  const { setSource } = useScatterActions();

  const onDiskClick = (disk: string) => () => {
    console.log('onDiskClick ', disk);
    setSource(disk);
  };

  return (
    <div className="flex flex-1 flex-col bg-neutral-200 dark:bg-gray-950">
      <div
        className="overflow-y-auto px-2 pt-2"
        style={{ height: `${height}px` }}
      >
        {disks.map((disk) => (
          <div
            className={`py-2 px-3 text-blue-800 ${selectedBackground(
              disk.name === selected,
            )}`}
            onClick={onDiskClick(disk.name)}
          >
            <div>
              <span className="font-bold">{disk.name}</span>
              <span className="dark:text-slate-700 text-slate-500 text-sm">
                &nbsp;({disk.fsType})
              </span>{' '}
              <span className="dark:text-gray-900 text-neutral-400">-</span>{' '}
              <span className="dark:text-slate-700 text-slate-500 text-sm">
                {humanBytes(disk.free)}{' '}
                <span className="dark:text-gray-900 text-neutral-400">/</span>{' '}
                {humanBytes(disk.size)}
              </span>
            </div>
            <p className="dark:text-indigo-500 text-indigo-500 text-sm">
              {disk.serial}
            </p>
          </div>
        ))}
      </div>
    </div>
  );
};
