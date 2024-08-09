import React from 'react';

import { Panel } from '~/shared/panel/panel';
import { useUnraidPlan, useUnraidDisks } from '~/state/unraid';
import { useGatherLocation } from '~/state/gather';
import { Selectable } from '~/shared/selectable/selectable';
import { Icon } from '~/shared/icons/icon';
import { Disk } from '~/shared/disk/disk';
import { FreePanel } from '~/shared/disk/free-panel';
import { humanBytes } from '~/helpers/units';
import { useGatherTarget, useGatherActions } from '~/state/gather';
import { Disk as IDisk } from '~/types';

const getPresence = (location: Record<string, string[]>, id: string) => {
  for (const key in location) {
    if (location[key].includes(id)) {
      return true;
    }
  }
  return false;
};

export const Target: React.FunctionComponent = () => {
  const plan = useUnraidPlan();
  const disks = useUnraidDisks();
  const location = useGatherLocation();
  const target = useGatherTarget();
  const { setTarget } = useGatherActions();

  if (!plan) {
    return (
      <div className="h-full flex flex-col bg-neutral-100 dark:bg-gray-950">
        <div className="flex flex-col p-2">
          <h1 className="text-lg text-slate-500 dark:text-gray-500 pb-2">
            Target
          </h1>
          <hr className="border-slate-300 dark:border-gray-700" />
        </div>
        <h1>no plan</h1>
      </div>
    );
  }

  // if no bin then this disk isn't elegible as a target
  const elegible = disks.filter(
    (disk) => plan.vdisks[disk.path].bin && plan.vdisks[disk.path].bin.items,
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

  const onDiskClick = (disk: IDisk) => () => setTarget(disk.path);

  return (
    <Panel title="Target">
      {targets.map((disk) => {
        const present = getPresence(location, disk.name);
        const fill = present
          ? 'fill-green-600 dark:fill-green-600'
          : 'fill-neutral-200 dark:fill-gray-950';
        return (
          <Selectable
            key={disk.id}
            onClick={onDiskClick(disk)}
            selected={disk.path === target}
          >
            <div className="grid grid-cols-12 gap-1 items-center">
              <div className="col-span-2 flex flex-row items-center">
                <Icon name="star" size={20} style={fill} />
                <span className="pr-2" />
                <span className="text-slate-500 dark:text-gray-500">
                  {humanBytes(disk.free - plan.vdisks[disk.path].plannedFree)}
                </span>
              </div>
              <div className="col-span-5 flex flex-row items-center">
                <Disk disk={disk} />
                <div className="pr-2" />
              </div>
              <div className="col-span-5">
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
    </Panel>
  );
};
