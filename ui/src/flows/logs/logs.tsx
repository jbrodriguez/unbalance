import React from 'react';

import { Panel } from '~/shared/panel/panel';
import { useUnraidActions, useUnraidLogs } from '~/state/unraid';

export const Logs: React.FunctionComponent = () => {
  const { getLog } = useUnraidActions();
  const logs = useUnraidLogs();

  React.useEffect(() => {
    getLog();
  }, [getLog]);

  return (
    <Panel>
      {logs.map((log) => (
        <p>{log}</p>
      ))}
    </Panel>
  );
};
