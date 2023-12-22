import React from 'react';

import { Panel, PanelGroup } from 'react-resizable-panels';
import { ResizeHandle } from '~/shared/resize/resize-handle';

interface PanelsProps {
  type?: '2col' | '3col';
  left?: React.ReactNode;
  middle?: React.ReactNode;
  right?: React.ReactNode;
  dimensions?: { left: number; middle?: number; right?: number };
}

export const Panels: React.FunctionComponent<PanelsProps> = ({
  type = '3col',
  left,
  middle,
  right,
  dimensions = { left: 30, right: 30 },
}) => {
  return (
    <div>
      <PanelGroup direction="horizontal">
        <Panel
          className="flex flex-row"
          defaultSizePercentage={dimensions.left}
          minSizePercentage={20}
        >
          {left}
        </Panel>
        <ResizeHandle />
        <Panel
          className="flex flex-row"
          defaultSizePercentage={dimensions.middle}
          minSizePercentage={20}
        >
          {middle}
        </Panel>
        {type === '3col' && (
          <>
            <ResizeHandle />
            <Panel
              className="flex flex-row"
              defaultSizePercentage={dimensions.right}
              minSizePercentage={20}
            >
              {right}
            </Panel>
          </>
        )}
      </PanelGroup>
    </div>
  );
};
