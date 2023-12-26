import React from 'react';

import { useUnraidOperation } from '~/state/unraid';
import { CommandStatus } from '~/types';
import { Icon } from '~/shared/icons/icon';

const getCommandStatus = (status: CommandStatus): React.ReactNode => {
  switch (status) {
    case CommandStatus.Complete:
      return (
        <Icon
          name="check-circle"
          size={20}
          fill="fill-green-600 dark:fill-green-600"
        />
      );
    case CommandStatus.Pending:
      return (
        <Icon
          name="minus-circle"
          size={20}
          fill="fill-blue-600 dark:fill-blue-600"
        />
      );
    case CommandStatus.Flagged:
      return (
        <Icon
          name="check-circle"
          size={20}
          fill="fill-yellow-600 dark:fill-yellow-600"
        />
      );
    case CommandStatus.Stopped:
      return (
        <Icon
          name="minus-circle"
          size={20}
          fill="fill-red-600 dark:fill-red-600"
        />
      );
    case CommandStatus.SourceRemoval:
      return (
        <Icon
          name="loading"
          size={20}
          fill="fill-yellow-600 dark:fill-yellow-600 animate-spin"
        />
      );
    default:
      return (
        <Icon
          name="loading"
          size={20}
          fill="fill-slate-600 dark:fill-slate-600 animate-spin"
        />
      );
  }
};

export const Commands: React.FunctionComponent = () => {
  const operation = useUnraidOperation();

  if (!operation) {
    return null;
  }

  return (
    <div className="overflow-x-auto">
      <table className="w-full text-sm text-left rtl:text-right text-gray-500 dark:text-gray-400">
        <thead className="text-xs text-gray-700 uppercase bg-gray-200 dark:bg-gray-700 dark:text-gray-400">
          <tr>
            <th scope="col" className="p-4"></th>
            <th scope="col" className="px-6 py-3">
              Source
            </th>
            <th scope="col" className="px-6 py-3">
              Command
            </th>
            <th scope="col" className="px-6 py-3">
              Progress
            </th>
          </tr>
        </thead>
        <tbody>
          {operation.commands.map((command) => {
            const progress = (
              (command.transferred / command.size) *
              100
            ).toFixed(0);
            return (
              <tr className="bg-gray-300 border-b dark:bg-gray-800 dark:border-gray-700 hover:bg-gray-50 dark:hover:bg-gray-600">
                <td className="w-4 p-4">{getCommandStatus(command.status)}</td>
                <th scope="row" className="px-6 py-4">
                  {command.src}
                </th>
                <td className="px-6 py-4  font-medium text-gray-900 whitespace-nowrap dark:text-white">
                  rsync {operation.rsyncStrArgs} &quot;{command.entry}&quot;
                  &quot;{command.dst}&quot;
                </td>
                <td className="flex flex-1 px-6 py-4">
                  <div className="w-full rounded bg-gray-400 dark:bg-gray-800">
                    <div
                      className="p-0.5 leading-none rounded bg-blue-900 "
                      style={{ width: `${progress}%` }}
                    ></div>
                  </div>
                </td>
              </tr>
            );
          })}
        </tbody>
      </table>
    </div>
  );
};
