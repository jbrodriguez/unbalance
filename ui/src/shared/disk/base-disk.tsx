import React from 'react';

import { Disk as IDisk } from '~/types';
import { humanBytes } from '~/helpers/units';

type DiskProps = {
  disk: IDisk;
};

export const Disk: React.FunctionComponent<DiskProps> = ({ disk }) => {
  return (
    <div className="flex flex-col text-blue-800">
      <div>
        <span className="font-bold">{disk.name}</span>
        <span className="dark:text-slate-700 text-slate-500 text-sm">
          &nbsp;({disk.fsType})
        </span>{' '}
        <span className="dark:text-gray-900 text-neutral-400">-</span>{' '}
        <span className="text-sm">
          <span className="dark:text-slate-500 text-slate-700 text-sm">
            {humanBytes(disk.free)}{' '}
          </span>
          <span className="dark:text-gray-900 text-neutral-400">/</span>{' '}
          <span className="dark:text-slate-700 text-slate-500">
            {humanBytes(disk.size)}
          </span>
        </span>
      </div>
      <p className="dark:text-indigo-500 text-indigo-500 text-sm">
        {disk.serial}
      </p>
    </div>
  );
};
