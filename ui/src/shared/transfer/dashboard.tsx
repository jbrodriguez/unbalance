import React from 'react';

import dayjs from 'dayjs';

import { useUnraidOperation } from '~/state/unraid';
import { formatBytes, formatTime } from '~/helpers/units';

export const Dashboard: React.FunctionComponent = () => {
  const operation = useUnraidOperation();

  if (!operation) {
    return null;
  }

  const completion = Math.round(operation.completed * 100) / 100;
  const completed = isNaN(completion) ? '0' : completion.toFixed(2);
  const velocity = Math.round(operation.speed * 100) / 100;
  const speed = isNaN(velocity) ? '0' : velocity.toFixed(2);

  let bytes = formatBytes(operation.bytesTransferred + operation.deltaTransfer);
  const transferredValue = bytes.value;
  const transferredUnit = bytes.unit;

  bytes = formatBytes(operation.bytesToTransfer);
  const totalValue = bytes.value;
  const totalUnit = bytes.unit;

  const started = dayjs(operation.started);
  const diff = dayjs().diff(started, 'second');
  const timeElapsed = formatTime(diff);

  const elapsed = timeElapsed === '' ? 'n/a' : timeElapsed;
  const remaining =
    !operation.remaining || operation.remaining === ''
      ? 'n/a'
      : operation.remaining;

  return (
    <div className="grid grid-cols-6 gap-4 text-blue-600 ">
      <div className="border border-stroke dark:border-gray-800 border-slate-300 px-4 py-3 shadow-default dark:border-strokedark dark:bg-boxdark">
        <div className="flex justify-between items-center">
          <span className="text-3xl">{completed}</span>
          <span className="text-sm font-medium dark:text-slate-400 text-slate-600">
            %
          </span>
        </div>
        <span className="text-sm font-medium dark:text-slate-600 text-slate-600">
          Completed
        </span>
      </div>

      <div className="border border-stroke dark:border-gray-800 border-slate-300 px-4 py-3 shadow-default dark:border-strokedark dark:bg-boxdark">
        <div className="flex justify-between items-center">
          <span className="text-3xl">{speed}</span>
          <span className="text-sm font-medium dark:text-slate-400 text-slate-600">
            MB/s
          </span>
        </div>
        <span className="text-sm font-medium dark:text-slate-600 text-slate-600">
          Speed
        </span>
      </div>

      <div className="border border-stroke dark:border-gray-800 border-slate-300 px-4 py-3 shadow-default dark:border-strokedark dark:bg-boxdark">
        <div className="flex justify-between items-center">
          <span className="text-3xl">{transferredValue}</span>
          <span className="text-sm font-medium dark:text-slate-400 text-slate-600">
            {transferredUnit}
          </span>
        </div>
        <span className="text-sm font-medium dark:text-slate-600 text-slate-600">
          Transferred
        </span>
      </div>

      <div className="border border-stroke dark:border-gray-800 border-slate-300 px-4 py-3 shadow-default dark:border-strokedark dark:bg-boxdark">
        <div className="flex justify-between items-center">
          <span className="text-3xl">{totalValue}</span>
          <span className="text-sm font-medium dark:text-slate-400 text-slate-600">
            {totalUnit}
          </span>
        </div>
        <span className="text-sm font-medium dark:text-slate-600 text-slate-600">
          Total
        </span>
      </div>

      <div className="border border-stroke dark:border-gray-800 border-slate-300 px-4 py-3 shadow-default dark:border-strokedark dark:bg-boxdark">
        <div className="flex justify-between items-center">
          <span className="text-3xl">{elapsed}</span>
          <span></span>
        </div>
        <span className="text-sm font-medium dark:text-slate-600 text-slate-600">
          Elapsed
        </span>
      </div>

      <div className="border border-stroke dark:border-gray-800 border-slate-300 px-4 py-3 shadow-default dark:border-strokedark dark:bg-boxdark">
        <div className="flex justify-between items-center">
          <span className="text-3xl">{remaining}</span>
          <span></span>
        </div>
        <span className="text-sm font-medium dark:text-slate-600 text-slate-600">
          Remaining
        </span>
      </div>
    </div>
  );
};
