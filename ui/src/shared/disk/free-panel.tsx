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

      <div className="flex flex-row font-mono text-sky-700 dark:text-slate-500">
        <span
          className="text-xs rounded leading-none bg-slate-700 text-slate-300 dark:bg-sky-900 dark:text-sky-500 py-0.5 px-1"
          style={{ writingMode: 'vertical-lr', textOrientation: 'upright' }}
        >
          free
        </span>
        <span className="pr-2" />
        <div className="flex flex-1 flex-col justify-around">
          <div className="grid grid-cols-12 gap-1 items-center">
            <span className="col-span-8 text-xs">
              current {`${current}% (${humanBytes(currentFree)})`}
            </span>
            <div className="col-span-4">
              <div className="w-full bg-gray-300 rounded dark:bg-gray-800">
                <div
                  className="text-xs text-center p-0.5 leading-none rounded bg-red-900 "
                  style={{ width: `${current}%` }}
                ></div>
              </div>
            </div>
          </div>
          <div className="grid grid-cols-12 gap-1 items-center">
            <span className="col-span-8 text-xs">
              planned {`${planned}% (${humanBytes(plannedFree)})`}
            </span>
            <div className="col-span-4">
              <div className="w-full bg-gray-300 rounded dark:bg-gray-800">
                <div
                  className="text-xs text-center p-0.5 leading-none rounded bg-green-900 "
                  style={{ width: `${planned}%` }}
                ></div>
              </div>
            </div>
          </div>
        </div>
      </div>

      <span className="pb-2"></span>
    </>
    // <>
    //   <span className="pt-2"></span>

    //   <div className="flex flex-row items-center">
    //     <span className="text-xs">current {`${currentFree}% (size)`}</span>
    //     <div className="w-full h-1 rounded bg-gray-200 dark:bg-gray-800">
    //       <div
    //         className="text-xs h-1 p-0.5 leading-none rounded bg-red-900"
    //         style={{ width: `${currentFree}%` }}
    //       ></div>
    //     </div>
    //   </div>

    //   <span className="pt-2"></span>

    //   <div className="grid grid-cols-12 gap-1 items-center">
    //     <span className="col-span-1">p</span>
    //     <span className="col-span-5 text-xs">
    //       planned {`${plannedFree}% (size)`}
    //     </span>
    //     <div className="col-span-6">
    //       <div className="w-full bg-gray-200 rounded dark:bg-gray-800">
    //         <div
    //           className="text-xs text-center p-0.5 leading-none rounded bg-green-900 "
    //           style={{ width: `${plannedFree}%` }}
    //         ></div>
    //       </div>
    //     </div>
    //   </div>

    //   <span className="pb-2"></span>
    // </>
    // <div className="flex flex-1">
    //   <div className="bg-neutral-200 dark:bg-gray-950">
    //     <div className="p-2">
    //       <div
    //         className="h-full bg-gray-200 rounded-full w-1.5 dark:bg-gray-700"
    //         style={{ height: `${currentFree}%` }}
    //       >
    //         <div className="bg-blue-600 w-1.5 rounded-full"></div>
    //       </div>
    //     </div>
    //     <div className="p-2">
    //       <div className="h-full bg-gray-200 rounded-full w-1.5 dark:bg-gray-700">
    //         <div
    //           className="bg-blue-600 w-1.5 rounded-full"
    //           style={{ height: `${plannedFree}%` }}
    //         ></div>
    //       </div>
    //     </div>
    //   </div>
    // </div>
  );

  // return (
  //   <div className="flex flex-1 flex-row items-center">
  //     <span className="text-xs">free</span>
  //     <div className="bg-neutral-200 dark:bg-gray-950">
  //       <div className="p-2 flex flex-row">
  //         <span className="text-xs">current</span>
  //         <div
  //           className="bg-gray-200 rounded-full h-1.5 dark:bg-gray-700"
  //           style={{ width: `${currentFree}%` }}
  //         >
  //           <div className="bg-blue-600 h-1.5 rounded-full"></div>
  //         </div>
  //         <span className="text-xs">size</span>
  //       </div>
  //       <div className="p-2">
  //         <span className="text-xs">current</span>
  //         <div className="w-full bg-gray-200 rounded-full h-1.5 dark:bg-gray-700">
  //           <div
  //             className="bg-blue-600 h-1.5 rounded-full"
  //             style={{ width: `${plannedFree}%` }}
  //           ></div>
  //         </div>
  //         <span className="text-xs">size</span>
  //       </div>
  //     </div>
  //   </div>
  // );
};
