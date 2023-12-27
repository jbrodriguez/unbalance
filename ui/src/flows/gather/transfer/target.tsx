import React from 'react';

import { useUnraidPlan, useUnraidDisks } from '~/state/unraid';
import { useGatherLocation } from '~/state/gather';
import { Selectable } from '~/shared/disk/selectable-disk';
import { Icon } from '~/shared/icons/icon';
import { Disk } from '~/shared/disk/base-disk';
import { FreePanel } from '~/shared/disk/free-panel';
import { humanBytes } from '~/helpers/units';
import { useGatherTarget, useGatherActions } from '~/state/gather';
import { Disk as IDisk } from '~/types';

interface Props {
  height?: number;
}

const getPresence = (location: Record<string, string[]>, id: string) => {
  for (const key in location) {
    if (location[key].includes(id)) {
      return true;
    }
  }
  return false;
};

export const Target: React.FunctionComponent<Props> = ({ height = 0 }) => {
  const plan = useUnraidPlan();
  const disks = useUnraidDisks();
  const location = useGatherLocation();
  const target = useGatherTarget();
  const { setTarget } = useGatherActions();

  if (!plan) {
    return (
      <div className="flex flex-1 flex-col bg-neutral-200 dark:bg-gray-950">
        <div className="p-2 overflow-y-auto" style={{ height: `${height}px` }}>
          <h1>no plan</h1>
        </div>
      </div>
    );
  }

  // if free === plannedFree then this disk isn't elegible as a target
  const elegible = disks.filter(
    (disk) => disk.free !== plan.vdisks[disk.path].plannedFree,
  );

  // sort elegible disks by least amount of data transfer
  const targets = elegible.sort((a, b) => {
    const xferA = a.free - plan.vdisks[a.path].plannedFree;
    const xferB = b.free - plan.vdisks[b.path].plannedFree;
    if (xferA < xferB) return -1;
    if (xferA > xferB) return 1;
    if (a.id < b.id) return -1;
    if (a.id > b.id) return 1;
    return 0;
  });

  const onDiskClick = (disk: IDisk) => {
    console.log('disk clicked', disk);
    setTarget(disk.path);
  };

  return (
    <div className="flex flex-1 flex-col bg-neutral-200 dark:bg-gray-950">
      <div className="p-2 overflow-y-auto" style={{ height: `${height}px` }}>
        {targets.map((disk) => {
          const present = getPresence(location, disk.name);
          const fill = present
            ? 'fill-green-600 dark:fill-green-600'
            : 'fill-neutral-200 dark:fill-gray-950';
          return (
            <Selectable
              disk={disk}
              onSelectDisk={onDiskClick}
              selected={disk.path === target}
            >
              <div className="flex flex-row items-center justify-between">
                <Icon name="star" size={20} fill={fill} />
                <span className="pr-4" />
                <span>
                  {humanBytes(disk.free - plan.vdisks[disk.path].plannedFree)}
                </span>
                <span className="pr-4" />
                <Disk disk={disk} />
                <span className="pr-4" />
                <div className="flex flex-col flex-1">
                  <FreePanel
                    size={disk.size}
                    currentFree={plan.vdisks[disk.path].currentFree}
                    plannedFree={plan.vdisks[disk.path].plannedFree}
                  />
                </div>
              </div>
            </Selectable>
          );
        })}
      </div>
    </div>
  );
};
