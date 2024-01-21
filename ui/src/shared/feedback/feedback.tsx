import React from 'react';

import {
  ResizableHandle,
  ResizablePanel,
  ResizablePanelGroup,
} from '@/components/ui/resizable';

import { Log } from '~/shared/log/log';
import { Issues } from '~/shared/issues/issues';

export const Feedback: React.FunctionComponent = () => {
  return (
    <div className="flex flex-1">
      <ResizablePanelGroup direction="horizontal" className="flex flex-1">
        <ResizablePanel defaultSize={60}>
          <Log />
        </ResizablePanel>
        <ResizableHandle withHandle />
        <ResizablePanel defaultSize={40}>
          <Issues />
        </ResizablePanel>
      </ResizablePanelGroup>
    </div>
  );
};
