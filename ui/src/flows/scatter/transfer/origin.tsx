import React from 'react';

import AutoSizer from 'react-virtualized-auto-sizer';

import { useUnraidPlan, useUnraidDisks } from '~/state/unraid';
import { Selectable } from '~/shared/disk/selectable-disk';
import { Disk } from '~/shared/disk/base-disk';
import { FreePanel } from '~/shared/disk/free-panel';

export const Origin: React.FunctionComponent = () => {
  const plan = useUnraidPlan();
  const disks = useUnraidDisks();

  if (!plan) {
    return (
      <div className="h-full flex flex-col bg-neutral-100 dark:bg-gray-950">
        <div className="flex flex-col p-2">
          <h1 className="text-lg text-slate-500 dark:text-gray-500 pb-2">
            Origin
          </h1>
          <hr className="border-slate-300 dark:border-gray-700" />
        </div>
      </div>
    );
  }

  const items = disks.filter((disk) => plan?.vdisks[disk.path].src);

  return (
    <div className="h-full flex flex-col bg-neutral-100 dark:bg-gray-950">
      <div className="flex flex-col p-2">
        <h1 className="text-lg text-slate-500 dark:text-gray-500 pb-2">
          Origin
        </h1>
        <hr className="border-slate-300 dark:border-gray-700" />
      </div>
      <div className="flex-auto">
        <AutoSizer disableWidth>
          {({ height }) => (
            <div
              className="p-2 overflow-y-auto"
              style={{ height: `${height}px` }}
            >
              {items.map((disk) => (
                <Selectable disk={disk}>
                  <div className="flex flex-col">
                    <Disk disk={disk} />
                    <FreePanel
                      size={disk.size}
                      currentFree={plan.vdisks[disk.path].currentFree}
                      plannedFree={plan.vdisks[disk.path].plannedFree}
                    />
                  </div>
                </Selectable>
              ))}
            </div>
          )}
        </AutoSizer>
      </div>
    </div>
  );
};
