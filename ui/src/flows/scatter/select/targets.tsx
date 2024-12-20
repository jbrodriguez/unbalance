import React from 'react';

import { Panel } from '~/shared/panel/panel';
import { useUnraidDisks } from '~/state/unraid';
import {
  useScatterSource,
  useScatterTargets,
  useScatterActions,
} from '~/state/scatter';
import { Checkbox } from '~/shared/checkbox/checkbox';
import { Disk } from '~/shared/disk/disk';
import { Disk as IDisk } from '~/types';
import { Toggle } from './toggle';

export const Targets: React.FunctionComponent = () => {
  const [allChecked, setAllChecked] = React.useState(false);

  const disks = useUnraidDisks();
  const source = useScatterSource();
  const targets = useScatterTargets();
  const { toggleTarget, toggleAll } = useScatterActions();

  const visible = source !== '';
  const elegible = disks.filter((disk) => disk.name !== source);

  const onCheck = (disk: IDisk) => () => toggleTarget(disk.name);
  const onToggleAll = (_checked: boolean) => () => {
    setAllChecked(!_checked);
    toggleAll(!_checked);
  };

  return (
    <Panel
      title="Target Disk(s)"
      subtitle={
        <Toggle allChecked={allChecked} onCheck={onToggleAll(allChecked)} />
      }
    >
      {visible &&
        elegible.map((disk) => (
          <div className="flex flex-row items-center">
            <Checkbox checked={targets[disk.name]} onCheck={onCheck(disk)} />
            <div className="pr-4" />
            <Disk disk={disk} />
          </div>
        ))}
    </Panel>
  );
};
