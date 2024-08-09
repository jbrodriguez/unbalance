import React from 'react';

import { Panel } from '~/shared/panel/panel';
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
            Items (per disk)
          </h1>
          <hr className="border-slate-300 dark:border-gray-700" />
        </div>
      </div>
    );
  }

  const bin = plan.vdisks[disk].bin;

  if (!bin || !bin.items) {
    return (
      <div className="h-full bg-neutral-100 dark:bg-gray-950">
        <div className="flex flex-col p-2">
          <h1 className="text-lg text-slate-500 dark:text-gray-500 pb-2">
            Items (per disk)
          </h1>
          <hr className="border-slate-300 dark:border-gray-700" />
          <span className="text-sm">No items in the bin.</span>
        </div>
      </div>
    );
  }

  return (
    <Panel title="Items (per disk)">
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
    </Panel>
  );
};
