import React from 'react';

import { Panel } from '~/shared/panel/panel';
import { useGatherSelected } from '~/state/gather';

export const Shares: React.FunctionComponent = () => {
  const selected = useGatherSelected();

  return (
    <Panel title="Shares">
      {Object.values(selected).map((share) => (
        <div className="whitespace-nowrap">{share}</div>
      ))}
    </Panel>
  );
};
