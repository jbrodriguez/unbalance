import React from 'react';

import AutoSizer from 'react-virtualized-auto-sizer';

import { useUnraidPlan } from '~/state/unraid';
import { humanBytes } from '~/helpers/units';

interface BinProps {
  height?: number;
  disk: string;
}

export const Bin: React.FunctionComponent<BinProps> = ({
  height,
  disk = '',
}) => {
  const plan = useUnraidPlan();

  if (!plan || disk === '') {
    return (
      <div className="flex flex-1 flex-col bg-neutral-200 dark:bg-gray-950">
        <div className="p-2 overflow-y-auto" style={{ height: `${height}px` }}>
          <span />
        </div>
      </div>
    );
  }

  console.log('disk, vdisk ########## ', disk, plan.vdisks);

  const bin = plan.vdisks[disk].bin;

  if (!bin) {
    return (
      <div className="flex flex-1 flex-col bg-neutral-200 dark:bg-gray-950">
        <div className="p-2 overflow-y-auto" style={{ height: `${height}px` }}>
          <span className="text-gray-500 dark:text-gray-700 text-sm">
            No items in the bin.
          </span>
        </div>
      </div>
    );
  }

  return (
    <div className="flex flex-1 bg-neutral-200 dark:bg-gray-950">
      <AutoSizer disableHeight>
        {({ width }) => (
          <div
            className="p-2 overflow-y-auto overflow-x-auto text-nowrap"
            style={{
              height: `${height}px`,
              width: `${width}px`,
            }}
          >
            {bin.items.map((item) => (
              <div>
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
  );
};
