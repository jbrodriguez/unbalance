import React from 'react';

import { Disk as IDisk } from '~/types';
import { humanBytes } from '~/helpers/units';

type DiskProps = {
  disk: IDisk;
};

export const Disk: React.FunctionComponent<DiskProps> = ({ disk }) => {
  return (
    <div className="flex flex-1 flex-col text-blue-800">
      <div className="flex flex-1 flex-row items-center justify-between">
        <div className="flex flex-1 flex-row items-center">
          <span className="font-bold">{disk.name}</span>
          <span className="dark:text-slate-700 text-slate-500 text-sm">
            &nbsp;({disk.fsType})
          </span>
          <span className="px-1">-</span>
          <span className="dark:text-slate-700 text-slate-500">
            {humanBytes(disk.size)}
          </span>
        </div>
        <span className="dark:text-slate-500 text-slate-700">
          {humanBytes(disk.free)}
        </span>
      </div>
      <p className="dark:text-indigo-500 text-indigo-500 text-sm">
        {disk.serial}
      </p>
    </div>
  );
};
