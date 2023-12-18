import React from 'react';

import { Disk as IDisk } from '~/types';
import { humanBytes } from '~/helpers/units';
import { Icon } from '~/shared/icons/icon';

type DiskProps = {
  disk: IDisk;
  checkable?: boolean;
  checked?: boolean;
  onSelect?: (disk: IDisk) => void;
  onCheck?: (disk: IDisk) => void;
  selectedStyle?: string;
};

export const Disk: React.FunctionComponent<DiskProps> = ({
  disk,
  checkable = false,
  checked = false,
  onSelect,
  onCheck,
  selectedStyle = '',
}) => {
  const onCheckClick = (disk: IDisk) => () => {
    console.log('onCheckClick ', disk);
    onCheck?.(disk);
  };

  const onSelectClick = (disk: IDisk) => () => {
    console.log('onSelectClick ', disk);
    onSelect?.(disk);
  };

  return (
    <div
      className={`py-2 px-3 flex flex-row items-center text-blue-800 ${selectedStyle}`}
      onClick={onSelectClick(disk)}
    >
      {checkable ? (
        <span onClick={onCheckClick(disk)} className="mr-4">
          {checked ? (
            <Icon
              name="checked"
              size={20}
              fill="fill-slate-700 dark:fill-slate-200"
            />
          ) : (
            <Icon
              name="unchecked"
              size={20}
              fill="fill-slate-700 dark:fill-slate-200"
            />
          )}
        </span>
      ) : null}
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
      </div>
    </div>
  );
};
