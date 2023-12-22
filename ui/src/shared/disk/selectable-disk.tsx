import React from 'react';

import { Disk as IDisk } from '~/types';

type DiskProps = {
  disk: IDisk;
  onSelectDisk?: (disk: IDisk) => void;
  selected?: boolean;
  children?: React.ReactNode;
};

const selectedBackground = (selected: boolean) =>
  selected ? 'rounded dark:bg-gray-900 bg-neutral-300' : '';

export const Selectable: React.FunctionComponent<DiskProps> = ({
  disk,
  onSelectDisk,
  selected = false,
  children,
}) => {
  const onClick = (disk: IDisk) => () => onSelectDisk?.(disk);

  return (
    <div
      className={`py-2 px-3 ${selectedBackground(selected)}`}
      onClick={onClick(disk)}
    >
      {children}
    </div>
  );
};
