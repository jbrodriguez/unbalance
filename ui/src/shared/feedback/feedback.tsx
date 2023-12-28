import React from 'react';

import {
  ResizableHandle,
  ResizablePanel,
  ResizablePanelGroup,
} from '@/components/ui/resizable';

// import { useScatterLogs } from '~/state/scatter';
import { Log } from '~/shared/log/log';
import { Issues } from '~/shared/issues/issues';

export const Feedback: React.FunctionComponent = () => {
  const logs = [
    'Planning started',
    'Planning complete',
    'Planning complete',
    'Planning complete',
    'Planning complete',
    'Planning complete',
    'Planning complete',
    'Planning complete',
    'Planning complete',
    'Planning complete',
    'Planning complete',
    'Planning complete',
    'Planning complete',
    'Planning complete',
    'Planning complete',
    'Planning complete',
    'Planning complete',
    'Planning complete',
    'Planning complete',
    'Planning complete',
    'Planning complete',
    'Planning complete',
    'Planning complete',
    'Planning complete',
    'Planning complete',
    'Planning complete',
    'Planning complete',
    'Planning complete',
    'Planning complete',
    'Planning complete',
    'Planning complete',
    'Planning complete',
    'Planning complete',
    'Planning complete',
    'Planning complete',
    'Planning complete',
  ];
  return (
    <div className="flex flex-1">
      <ResizablePanelGroup
        direction="horizontal"
        className="flex flex-1 border border-slate-300 dark:border-gray-700"
      >
        <ResizablePanel defaultSizePercentage={60}>
          <Log logs={logs} />
        </ResizablePanel>
        <ResizableHandle withHandle />
        <ResizablePanel defaultSizePercentage={40}>
          <Issues />
        </ResizablePanel>
      </ResizablePanelGroup>
    </div>
  );
};
