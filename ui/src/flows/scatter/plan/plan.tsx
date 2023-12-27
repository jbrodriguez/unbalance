import React from 'react';

import { Log } from '~/shared/log/log';
import { useScatterLogs } from '~/state/scatter';

export const Plan: React.FunctionComponent = () => {
  const logs = useScatterLogs();
  return <Log logs={logs} />;
};
