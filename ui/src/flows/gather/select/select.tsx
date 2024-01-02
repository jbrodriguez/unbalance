import React from 'react';

import {
  ResizableHandle,
  ResizablePanel,
  ResizablePanelGroup,
} from '@/components/ui/resizable';

import { Shares } from './shares';
import { Presence } from './presence';

export const Select: React.FunctionComponent = () => {
  return (
    <div className="flex flex-1">
      <ResizablePanelGroup direction="horizontal" className="flex flex-1">
        <ResizablePanel defaultSizePercentage={50}>
          <Shares />
        </ResizablePanel>
        <ResizableHandle withHandle />
        <ResizablePanel defaultSizePercentage={50}>
          <Presence />
        </ResizablePanel>
      </ResizablePanelGroup>
    </div>
  );
};
