import React from 'react';

import { Panel } from '~/shared/panel/panel';
import { useUnraidDisks } from '~/state/unraid';
import { Disk as IDisk } from '~/types';
import { useScatterActions, useScatterSource } from '~/state/scatter';
import { Selectable } from '~/shared/selectable/selectable';
import { Disk } from '~/shared/disk/disk';

export const Disks: React.FunctionComponent = () => {
  const disks = useUnraidDisks();
  const selected = useScatterSource();
  const { setSource } = useScatterActions();

  const onDiskClick = (disk: IDisk) => () => setSource(disk.name);

  return (
    <Panel title="Source Disk">
      {disks.map((disk) => (
        <Selectable
          key={disk.id}
          onClick={onDiskClick(disk)}
          selected={disk.name === selected}
        >
          <Disk disk={disk} />
        </Selectable>
      ))}
    </Panel>
  );
};
