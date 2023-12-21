import React from 'react';

import { Disk as IDisk } from '~/types';
// import { humanBytes } from '~/helpers/units';
// import { Icon } from '~/shared/icons/icon';

type DiskProps = {
  disk: IDisk;
  onSelectDisk?: (disk: IDisk) => void;
  selected?: boolean;
  children?: React.ReactNode;
};

const selectedBackground = (selected: boolean) =>
  selected ? 'rounded dark:bg-gray-900 bg-neutral-300' : '';

export const Selectable: React.FunctionComponent<DiskProps> = ({
  disk,
  onSelectDisk,
  selected = false,
  children,
}) => {
  const onClick = (disk: IDisk) => () => {
    console.log('onSelectClick ', disk);
    // @ts-expect-error ts-migrate(2531) FIXME: Object is possibly 'null'.
    onSelectDisk(disk);
  };

  return (
    <div
      className={`py-2 px-3 ${selectedBackground(selected)}`}
      onClick={onClick(disk)}
    >
      {children}
      {/* 
      <div className="flex flex-col">
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
      </div> */}
    </div>
  );
};
