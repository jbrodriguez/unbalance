import React from 'react';

import * as dayjs from 'dayjs';

import { useUnraidOperation } from '~/state/unraid';
import { formatBytes, formatTime } from '~/helpers/units';

export const Dashboard: React.FunctionComponent = () => {
  const operation = useUnraidOperation();

  if (!operation) {
    return null;
  }

  const completed = (Math.round(operation.completed * 100) / 100).toFixed(2);
  const speed = (Math.round(operation.speed * 100) / 100).toFixed(2);

  let bytes = formatBytes(operation.bytesTransferred + operation.deltaTransfer);
  const transferredValue = bytes.value;
  const transferredUnit = ' ' + bytes.unit;

  bytes = formatBytes(operation.bytesToTransfer);
  const totalValue = bytes.value;
  const totalUnit = ' ' + bytes.unit;

  const started = dayjs(operation.started);
  const diff = dayjs().diff(started, 'second');
  const elapsed = formatTime(diff);

  return (
    <div className="grid grid-cols-5 gap-6 text-blue-600 ">
      <div className="border border-stroke dark:border-gray-800 border-slate-300 px-4 py-3 shadow-default dark:border-strokedark dark:bg-boxdark">
        <div className="flex justify-between items-center">
          <span className="text-3xl">{completed}</span>
          <span>%</span>
        </div>
        <span className="text-sm font-medium dark:text-slate-600 text-slate-600">
          Completed
        </span>
      </div>

      <div className="border border-stroke dark:border-gray-800 border-slate-300 px-4 py-3 shadow-default dark:border-strokedark dark:bg-boxdark">
        <div className="flex justify-between items-center">
          <span className="text-3xl">{speed}</span>
          <span>MB/s</span>
        </div>
        <span className="text-sm font-medium dark:text-slate-600 text-slate-600">
          Speed
        </span>
      </div>

      <div className="border border-stroke dark:border-gray-800 border-slate-300 px-4 py-3 shadow-default dark:border-strokedark dark:bg-boxdark">
        <div className="flex justify-between items-center">
          <span className="text-3xl">{transferredValue}</span>
          <span>{transferredUnit}</span>
        </div>
        <span className="text-sm font-medium dark:text-slate-600 text-slate-600">
          Transferred
        </span>
      </div>

      <div className="border border-stroke dark:border-gray-800 border-slate-300 px-4 py-3 shadow-default dark:border-strokedark dark:bg-boxdark">
        <div className="flex justify-between items-center">
          <span className="text-3xl">{totalValue}</span>
          <span>{totalUnit}</span>
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
          Total
        </span>
      </div>

      <div className="border border-stroke dark:border-gray-800 border-slate-300 px-4 py-3 shadow-default dark:border-strokedark dark:bg-boxdark">
        <div className="flex justify-between items-center">
          <span className="text-3xl">{operation.remaining}</span>
          <span></span>
        </div>
        <span className="text-sm font-medium dark:text-slate-600 text-slate-600">
          Remaining
        </span>
      </div>
    </div>
  );
};
