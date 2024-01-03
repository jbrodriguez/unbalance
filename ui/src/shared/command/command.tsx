import React from 'react';

import { Command as ICommand } from '~/types';
import { getCommandStatus } from '~/helpers/operation';

interface Props {
  command: ICommand;
  rsyncStrArgs: string;
  canBeFlagged?: boolean;
  onFlag?: (command: ICommand) => void;
}

export const Command: React.FunctionComponent<Props> = ({
  command,
  rsyncStrArgs,
  canBeFlagged = false,
  onFlag,
}) => {
  const progress = ((command.transferred / command.size) * 100).toFixed(0);
  const onClick = () => onFlag?.(command);

  return (
    <div className="grid grid-cols-12 gap-1 items-center text-sm text-gray-700 dark:text-gray-500 p-2 border-b border-slate-300 dark:border-gray-700 ">
      <div className="col-span-2 flex items-center">
        {canBeFlagged ? (
          <button onClick={onClick}>
            <span className="bg-red-100 text-red-800 text-xs font-medium px-2.5 py-0.5 rounded dark:bg-red-900 dark:text-red-300">
              rmsrc
            </span>
          </button>
        ) : (
          getCommandStatus(command.status)
        )}
        <span className="pr-2" />
        {command.src}
      </div>
      <div className="col-span-8">
        rsync {rsyncStrArgs} &quot;{command.entry}&quot; &quot;{command.dst}
        &quot;
      </div>
      <div className="col-span-2 flex flex-1">
        <div className="w-full rounded bg-gray-400 dark:bg-gray-800">
          <div
            className="p-0.5 leading-none rounded bg-blue-900 "
            style={{ width: `${progress}%` }}
          ></div>
        </div>
      </div>
    </div>
  );
};
