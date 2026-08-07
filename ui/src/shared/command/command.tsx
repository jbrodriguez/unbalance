import React from 'react';

import { CommandStatus, Command as ICommand } from '~/types';
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
    <div className="grid grid-cols-1 gap-1 items-start lg:grid-cols-12 text-sm text-gray-500 dark:text-gray-500 p-2 border-b border-slate-300 dark:border-gray-700">
      <div className="col-span-1 lg:col-span-2 flex items-center flex-wrap">
        {(command.status === CommandStatus.SourceRemoval ||
          command.status === CommandStatus.Flagged) &&
        canBeFlagged ? (
          <button onClick={onClick}>
            <span className="bg-red-100 text-red-800 text-xs font-medium px-2.5 py-0.5 rounded dark:bg-red-900 dark:text-red-300">
              rmsrc
            </span>
          </button>
        ) : (
          getCommandStatus(command.status)
        )}
        <span className="pr-2" />
        <span className="text-xs lg:text-sm truncate">{command.src}</span>
      </div>
      <div className="col-span-1 lg:col-span-8 text-xs lg:text-sm break-all lg:break-normal">
        rsync {rsyncStrArgs} &quot;{command.entry}&quot; &quot;{command.dst}
        &quot;
      </div>
      <div className="col-span-1 lg:col-span-2 flex flex-row items-center gap-1">
        {progress !== '100' && <span className="text-xs">{progress} %</span>}
        <div className="flex-1 rounded bg-gray-400 dark:bg-gray-800 min-w-[50px]">
          <div
            className="p-0.5 leading-none rounded bg-blue-900"
            style={{ width: `${progress}%` }}
          ></div>
        </div>
      </div>
      {command.reason && (
        <div className="col-span-1 lg:col-span-12 text-xs text-yellow-600 dark:text-yellow-600">
          {command.reason}
        </div>
      )}
    </div>
  );
};
