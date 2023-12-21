import React from 'react';

import { useUnraidPlan, useUnraidDisks } from '~/state/unraid';
import { Selectable } from '~/shared/disk/selectable-disk';
import { Disk } from '~/shared/disk/base-disk';
import { FreePanel } from '~/shared/disk/free-panel';

interface Props {
  height?: number;
}

export const Destination: React.FunctionComponent<Props> = ({ height = 0 }) => {
  const plan = useUnraidPlan();
  const disks = useUnraidDisks();

  if (!plan) {
    return (
      <div className="flex flex-1 flex-col bg-neutral-200 dark:bg-gray-950">
        <div className="p-2 overflow-y-auto" style={{ height: `${height}px` }}>
          <h1>no plan</h1>
        </div>
      </div>
    );
  }

  const items = disks.filter((disk) => plan?.vdisks[disk.path].dst);

  return (
    <div className="flex flex-1 flex-col bg-neutral-200 dark:bg-gray-950">
      <div className="p-2 overflow-y-auto" style={{ height: `${height}px` }}>
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
    </div>
  );
};
