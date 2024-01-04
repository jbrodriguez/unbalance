import React from 'react';

import { Button } from '@/components/ui/button';
import dayjs from 'dayjs';

import {
  Operation as IOperation,
  Op,
  CommandStatus,
  ConfirmationParams,
  ConfirmationKind,
} from '~/types';
import { Icon } from '~/shared/icons/icon';
import { operationKindToName } from '~/helpers/operation';
import { formatTime } from '~/helpers/units';

interface Props {
  operation: IOperation;
  first: boolean;
  onConfirm: (params: ConfirmationParams) => void;
}

export const OperationHeader: React.FunctionComponent<Props> = ({
  operation,
  first,
  onConfirm,
}) => {
  const safe = first;
  const replay = !operation.dryRun && safe;
  const validate =
    !operation.dryRun && operation.opKind === Op.ScatterCopy && safe;

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

  // calculate runtime from started and finished dates, using dayjs
  const started = dayjs(operation.started);
  const finished = dayjs(operation.finished);
  const diff = finished.diff(started, 'seconds');
  const runtime = formatTime(diff);

  const canBeFlagged =
    safe &&
    (operation.opKind === Op.ScatterMove || operation.opKind === Op.GatherMove);

  const onValidate = () =>
    onConfirm({
      kind: ConfirmationKind.ScatterValidate,
      operation: operation,
    });

  const onReplay = () =>
    onConfirm({
      kind: ConfirmationKind.Replay,
      operation: operation,
    });

  return (
    <div className="flex flex-col pt-2 px-2">
      <div className="px-2 pb-2">
        <div className="flex items-center justify-between">
          <div className="flex flex-row items-center">
            {operationStatus}
            {operation.dryRun ? (
              <>
                <span className="pr-2" />
                <span className="bg-indigo-600 text-indigo-100 dark:bg-indigo-900 dark:text-indigo-300 text-xs font-medium px-2.5 py-0.5 rounded ">
                  dry run
                </span>
              </>
            ) : null}
            <span className="pr-2" />
            <span className="text-slate-600 dark:text-gray-400 text-lg">
              {operationKindToName[operation.opKind]}
            </span>
            <span className="pr-2" />
            <span className="text-slate-800 dark:text-gray-200 text-xl">
              {runtime}
            </span>
          </div>
          <div>
            {validate && (
              <Button variant="secondary" onClick={onValidate}>
                validate
              </Button>
            )}
            {replay && (
              <>
                <span className="pr-2" />
                <Button variant="secondary" onClick={onReplay}>
                  replay
                </Button>
              </>
            )}
          </div>
        </div>
      </div>
      {canBeFlagged && flagged ? (
        <div className="text-sm text-slate-500 dark:text-gray-500">
          <p>
            One or more commands had an execution warning/error. Check
            /var/log/unbalance.log for additional details.
          </p>
          <p>
            Due to this, the plugin hasn&apos;t deleted the source files/folders
            for that/those commands.
          </p>
          <p>
            Once you&apos;ve checked/solved the issue(s), click on the{' '}
            <span className="bg-red-100 text-red-800 text-xs font-medium px-2.5 py-0.5 rounded dark:bg-red-900 dark:text-red-300">
              rmsrc
            </span>{' '}
            button to remove the source files/folders, if you wish to do so.
          </p>
          <div className="pb-2" />
        </div>
      ) : null}
      <hr className="border-slate-300 dark:border-gray-700" />
    </div>
  );
};
