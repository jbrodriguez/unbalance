import React from 'react';

import { Panel } from '~/shared/panel/panel';
import { useUnraidDisks } from '~/state/unraid';
import {
  useScatterSource,
  useScatterSelected,
  useScatterTargets,
  useScatterActions,
} from '~/state/scatter';
import { Checkbox } from '~/shared/checkbox/checkbox';
import { Disk } from '~/shared/disk/disk';
import { Disk as IDisk, Targets as ITargets } from '~/types';

const isChecked = (name: string, targets: ITargets) => targets[name] || false;

export const Targets: React.FunctionComponent = () => {
  const disks = useUnraidDisks();
  const source = useScatterSource();
  const selected = useScatterSelected();
  const targets = useScatterTargets();
  const { toggleTarget } = useScatterActions();

  const visible = source !== '' && selected.length > 0;
  const elegible = disks.filter((disk) => disk.name !== source);

  const onCheck = (disk: IDisk) => () => toggleTarget(disk.name);

  return (
    <Panel title="Target Disk(s)">
      {visible &&
        elegible.map((disk) => (
          <div className="flex flex-row items-center">
            <Checkbox
              checked={isChecked(disk.name, targets)}
              onCheck={onCheck(disk)}
            />
            <div className="pr-4" />
            <Disk disk={disk} />
          </div>
        ))}
    </Panel>
  );
};
