import React from 'react';

import { useGatherSelected } from '~/state/gather';

interface Props {
  height?: number;
}

export const Shares: React.FunctionComponent<Props> = ({ height = 0 }) => {
  const selected = useGatherSelected();

  return (
    <div className="flex flex-1 flex-col bg-neutral-200 dark:bg-gray-950">
      <div className="p-2 overflow-y-auto" style={{ height: `${height}px` }}>
        {Object.values(selected).map((share) => (
          <div>{share}</div>
        ))}
      </div>
    </div>
  );
};
