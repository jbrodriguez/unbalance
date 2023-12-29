import React from 'react';

import AutoSizer from 'react-virtualized-auto-sizer';

import { useUnraidPlan } from '~/state/unraid';
import { humanBytes } from '~/helpers/units';

interface BinProps {
  disk: string;
}

export const Bin: React.FunctionComponent<BinProps> = ({ disk = '' }) => {
  const plan = useUnraidPlan();

  if (!plan || disk === '') {
    return (
      <div className="h-full bg-neutral-100 dark:bg-gray-950">
        <div className="flex flex-col p-2">
          <h1 className="text-lg text-slate-500 dark:text-gray-500 pb-2">
            Items
          </h1>
          <hr className="border-slate-300 dark:border-gray-700" />
        </div>
      </div>
    );
  }

  console.log('disk, vdisk ########## ', disk, plan.vdisks);

  const bin = plan.vdisks[disk].bin;

  if (!bin) {
    return (
      <div className="h-full bg-neutral-100 dark:bg-gray-950">
        <div className="flex flex-col p-2">
          <h1 className="text-lg text-slate-500 dark:text-gray-500 pb-2">
            Source Disk
          </h1>
          <hr className="border-slate-300 dark:border-gray-700" />
          <div className="flex flex-1 flex-col bg-neutral-200 dark:bg-gray-950">
            <span className="text-gray-500 dark:text-gray-700 text-sm">
              No items in the bin.
            </span>
          </div>
        </div>
      </div>
    );
  }

  return (
    <div className="h-full flex flex-col bg-neutral-100 dark:bg-gray-950">
      <div className="flex flex-col pt-2 px-2">
        <h1 className="text-lg text-slate-500 dark:text-gray-500 pb-2">
          Items (per disk)
        </h1>
        <hr className="border-slate-300 dark:border-gray-700" />
      </div>
      <div className="flex-auto">
        <AutoSizer disableWidth>
          {({ height }) => (
            <div
              className="overflow-y-auto overflow-x-auto p-2"
              style={{ height: `${height}px` }}
            >
              {bin.items.map((item) => (
                <div className="whitespace-nowrap">
                  <span className="text-gray-700 dark:text-gray-500 text-sm">
                    ({humanBytes(item.size)}){' '}
                  </span>
                  <span className="text-gray-500 dark:text-gray-700 text-sm">
                    {item.path}
                  </span>
                </div>
              ))}
            </div>
          )}
        </AutoSizer>
      </div>
    </div>
  );
};
