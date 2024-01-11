import React from 'react';

import {
  ResizableHandle,
  ResizablePanel,
  ResizablePanelGroup,
} from '@/components/ui/resizable';

import { Disks } from './disks';
import { FileSystem } from './filesystem';
import { Targets } from './targets';

export const Select: React.FunctionComponent = () => {
  return (
    <div className="flex flex-1">
      <ResizablePanelGroup direction="horizontal" className="flex flex-1">
        <ResizablePanel defaultSizePercentage={30}>
          <Disks />
        </ResizablePanel>
        <ResizableHandle withHandle />
        <ResizablePanel defaultSizePercentage={40}>
          <FileSystem />
        </ResizablePanel>
        <ResizableHandle withHandle />
        <ResizablePanel defaultSizePercentage={30}>
          <Targets />
        </ResizablePanel>
      </ResizablePanelGroup>
    </div>
  );
};
