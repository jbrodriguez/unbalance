import React from 'react';

import {
  ResizableHandle,
  ResizablePanel,
  ResizablePanelGroup,
} from '@/components/ui/resizable';

import { Shares } from './shares';
import { Target } from './target';
import { Bin } from '~/shared/bin/bin';
import { useGatherTarget } from '~/state/gather';

export const Targets: React.FunctionComponent = () => {
  const disk = useGatherTarget();

  return (
    <div className="flex flex-1">
      <ResizablePanelGroup direction="horizontal" className="flex flex-1">
        <ResizablePanel defaultSize={20}>
          <Shares />
        </ResizablePanel>
        <ResizableHandle withHandle />
        <ResizablePanel defaultSize={50}>
          <Target />
        </ResizablePanel>
        <ResizableHandle withHandle />
        <ResizablePanel defaultSize={30}>
          <Bin disk={disk} />
        </ResizablePanel>
      </ResizablePanelGroup>
    </div>
  );
};
