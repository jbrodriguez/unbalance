import React from "react"

import { Panel, PanelGroup } from "react-resizable-panels"
import { ResizeHandle } from "~/shared/resize/resize-handle"

interface PanelsProps {
  type?: "2col" | "3col"
  left?: React.ReactNode
  middle?: React.ReactNode
  right?: React.ReactNode
}

export const Panels: React.FC<PanelsProps> = ({
  type = "3col",
  left,
  middle,
  right,
}) => {
  return (
    <div>
      <PanelGroup direction="horizontal">
        <Panel
          className="flex flex-row"
          defaultSizePercentage={30}
          minSizePercentage={20}
        >
          {left}
        </Panel>
        <ResizeHandle />
        <Panel className="flex flex-row" minSizePercentage={30}>
          {middle}
        </Panel>
        {type === "3col" && (
          <>
            <ResizeHandle />
            <Panel
              className="flex flex-row"
              defaultSizePercentage={30}
              minSizePercentage={20}
            >
              {right}
            </Panel>
          </>
        )}
      </PanelGroup>
    </div>
  )
}
