import React from 'react';

import {
  ResizableHandle,
  ResizablePanel,
  ResizablePanelGroup,
} from '@/components/ui/resizable';

import { useScatterBinDisk } from '~/state/scatter';
import { Origin } from './origin';
import { Destination } from './destination';
import { Bin } from '~/shared/bin/bin';

export const Validation: React.FunctionComponent = () => {
  const disk = useScatterBinDisk();

  return (
    <div className="flex flex-1">
      <ResizablePanelGroup direction="horizontal" className="flex flex-1">
        <ResizablePanel defaultSize={30}>
          <Origin />
        </ResizablePanel>
        <ResizableHandle withHandle />
        <ResizablePanel defaultSize={30}>
          <Destination />
        </ResizablePanel>
        <ResizableHandle withHandle />
        <ResizablePanel defaultSize={40}>
          <Bin disk={disk} />
        </ResizablePanel>
      </ResizablePanelGroup>
    </div>
  );
};
