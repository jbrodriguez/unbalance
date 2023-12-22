import React from 'react';

import AutoSizer from 'react-virtualized-auto-sizer';

import { useUnraidPlan } from '~/state/unraid';
import { useScatterBinDisk } from '~/state/scatter';
import { humanBytes } from '~/helpers/units';

interface BinProps {
  height?: number;
}

export const Bin: React.FunctionComponent<BinProps> = ({ height }) => {
  const plan = useUnraidPlan();
  const binDisk = useScatterBinDisk();

  if (!plan || binDisk === '') {
    return (
      <div className="flex flex-1 flex-col bg-neutral-200 dark:bg-gray-950">
        <div className="p-2 overflow-y-auto" style={{ height: `${height}px` }}>
          <span />
        </div>
      </div>
    );
  }

  const bin = plan.vdisks[binDisk].bin;

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
