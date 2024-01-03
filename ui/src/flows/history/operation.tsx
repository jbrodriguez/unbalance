import React from 'react';

import AutoSizer from 'react-virtualized-auto-sizer';

import { Operation as IOperation } from '~/types';
import { OperationHeader } from './operation-header';
import { Command } from '~/shared/command/command';

interface Props {
  current: IOperation | null;
  first: boolean;
}

export const Operation: React.FunctionComponent<Props> = ({
  current,
  first,
}) => {
  if (!current) {
    return (
      <div className="h-full flex flex-col bg-neutral-100 dark:bg-gray-950">
        <div className="flex-auto">
          <AutoSizer disableWidth>
            {({ height }) => (
              <div
                className="overflow-y-auto overflow-x-auto p-2 text-slate-700 dark:text-gray-300"
                style={{ height: `${height}px` }}
              >
                <span />
              </div>
            )}
          </AutoSizer>
        </div>
      </div>
    );
  }

  return (
    <div className="h-full flex flex-col bg-neutral-100 dark:bg-gray-950">
      <OperationHeader operation={current} first={first} />
      <div className="flex-auto">
        <AutoSizer disableWidth>
          {({ height }) => (
            <div
              className="overflow-y-auto overflow-x-auto p-2 text-slate-700 dark:text-gray-300"
              style={{ height: `${height}px` }}
            >
              {current.commands.map((command) => (
                <Command
                  key={command.id}
                  command={command}
                  rsyncStrArgs={current.rsyncStrArgs}
                />
              ))}
            </div>
          )}
        </AutoSizer>
      </div>
    </div>
  );
};
