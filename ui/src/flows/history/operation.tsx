import React from 'react';

import AutoSizer from 'react-virtualized-auto-sizer';

import {
  ConfirmationParams,
  Operation as IOperation,
  Op,
  Command as ICommand,
  ConfirmationKind,
} from '~/types';
import { OperationHeader } from './operation-header';
import { Command } from '~/shared/command/command';
// import { Button } from '~/components/ui/button';

interface Props {
  current: IOperation | null;
  first: boolean;
  onConfirm: (params: ConfirmationParams) => void;
}

export const Operation: React.FunctionComponent<Props> = ({
  current,
  first,
  onConfirm,
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

  const safe = first;
  const canBeFlagged =
    safe &&
    (current.opKind === Op.ScatterMove || current.opKind === Op.GatherMove);

  const onRemoveSource = (command: ICommand) =>
    onConfirm({
      kind: ConfirmationKind.RemoveSource,
      operation: current,
      command,
    });

  return (
    <div className="h-full flex flex-col bg-neutral-100 dark:bg-gray-950">
      <OperationHeader
        operation={current}
        first={first}
        onConfirm={onConfirm}
      />
      <div className="flex-auto">
        <AutoSizer disableWidth>
          {({ height }) => (
            <div
              className="overflow-y-auto overflow-x-auto p-2"
              style={{ height: `${height}px` }}
            >
              {current.commands.map((command, index) => (
                <Command
                  key={command.id}
                  command={command}
                  rsyncStrArgs={current.rsyncStrArgs}
                  canBeFlagged={canBeFlagged && index === 3}
                  onFlag={onRemoveSource}
                />
              ))}
            </div>
          )}
        </AutoSizer>
      </div>
    </div>
  );
};
