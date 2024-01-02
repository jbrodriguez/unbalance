import React from 'react';

import AutoSizer from 'react-virtualized-auto-sizer';

import { Operation as IOperation, Op, CommandStatus } from '~/types';
import { Icon } from '~/shared/icons/icon';
import { operationKindToName } from '~/helpers/operation';

interface Props {
  operation: IOperation;
  first: boolean;
}

export const OperationHeader: React.FunctionComponent<Props> = ({
  operation,
  first,
}) => {
  const safe = first;
  const replay = !operation.dryRun && safe;
  // const validate =
  //   !operation.dryRun && operation.opKind === Op.ScatterCopy && safe;
  const validate = operation.opKind === Op.ScatterCopy && safe;

  const flagged = operation.commands.some(
    (command) => command.status === CommandStatus.Flagged,
  );
  const operationStatus = flagged ? (
    <Icon
      name="check-circle"
      size={14}
      style="fill-yellow-600 dark:fill-yellow-600"
    />
  ) : operation.bytesTransferred === operation.bytesToTransfer ? (
    <Icon
      name="check-circle"
      size={14}
      style="fill-green-600 dark:fill-green-600"
    />
  ) : (
    <Icon
      name="minus-circle"
      size={14}
      style="fill-red-600 dark:fill-red-600"
    />
  );

  return (
    <div className="h-full flex flex-col bg-neutral-100 dark:bg-gray-950">
      <div className="flex flex-col pt-2 px-2">
        <div className="text-slate-500 dark:text-gray-500 pb-2">
          <div className="flex items-center justify-between">
            <div className="flex flex-row items-center">
              {operationStatus}
              {operation.dryRun ? (
                <>
                  <span className="pr-2" />
                  <span className="bg-indigo-100 text-indigo-800 text-xs font-medium me-2 px-2.5 py-0.5 rounded dark:bg-indigo-900 dark:text-indigo-300">
                    dry run
                  </span>
                </>
              ) : null}
              <div className="text-lg">
                {operationKindToName[operation.opKind]}
              </div>
            </div>
            <div>
              {validate && (
                <button
                  // className={cx('btn', 'btn-primary', 'rspacer')}
                  // onClick={this.confirm('validate', op.id)}
                  disabled={!validate}
                >
                  VALIDATE
                </button>
              )}
              {replay && (
                <button
                  // className={cx('btn', 'btn-primary', 'rspacer')}
                  // onClick={this.confirm('replay', op.id)}
                  disabled={!replay}
                >
                  REPLAY
                </button>
              )}
            </div>
          </div>
        </div>
        <hr className="border-slate-300 dark:border-gray-700" />
      </div>
      <div className="flex-auto">
        <AutoSizer disableWidth>
          {({ height }) => (
            <div
              className="overflow-y-auto overflow-x-auto p-2 text-slate-700 dark:text-gray-300"
              style={{ height: `${height}px` }}
            ></div>
          )}
        </AutoSizer>
      </div>
    </div>
  );
};
