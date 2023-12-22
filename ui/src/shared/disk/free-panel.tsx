import React from 'react';

import { humanBytes } from '~/helpers/units';

interface FreePanelProps {
  size: number;
  currentFree: number;
  plannedFree: number;
}

export const FreePanel: React.FunctionComponent<FreePanelProps> = ({
  size,
  currentFree,
  plannedFree,
}) => {
  const current = ((currentFree / size) * 100).toFixed(0);
  const planned = ((plannedFree / size) * 100).toFixed(0);

  return (
    <>
      <span className="pt-2"></span>

      <div className="flex flex-row item-center font-mono text-xs text-sky-700 dark:text-slate-500">
        <span
          className="leading-none bg-slate-700 text-slate-300 dark:bg-gray-900 dark:text-gray-600 py-0.5 px-1"
          style={{ writingMode: 'vertical-lr', textOrientation: 'upright' }}
        >
          free
        </span>
        <span className="pr-2" />
        <div className="flex flex-1 flex-col justify-around">
          <div className="flex flex-1 items-center">
            <span className="pr-2">
              current {`${current}% (${humanBytes(currentFree)})`}
            </span>
            <div className="flex flex-1">
              <div className="w-full rounded bg-gray-400 dark:bg-gray-800">
                <div
                  className="p-0.5 leading-none rounded bg-red-900 "
                  style={{ width: `${current}%` }}
                ></div>
              </div>
            </div>
          </div>
          <div className="flex flex-1 items-center">
            <span className="pr-2">
              planned {`${planned}% (${humanBytes(plannedFree)})`}
            </span>
            <div className="flex flex-1">
              <div className="w-full rounded bg-gray-400 dark:bg-gray-800">
                <div
                  className="p-0.5 leading-none rounded bg-green-900 "
                  style={{ width: `${planned}%` }}
                ></div>
              </div>
            </div>
          </div>
        </div>
      </div>

      <span className="pb-2"></span>
    </>
  );
};
