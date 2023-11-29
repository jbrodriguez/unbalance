import React from "react"

import AutoSizer from "react-virtualized-auto-sizer"

import { Panels } from "~/shared/panels/panels"
import { Disks } from "./disks"
import { FileSystem } from "./filesystem"
import { Targets } from "./targets"

export const Scatter: React.FC = () => {
  return (
    <div className="flex flex-col h-full">
      <div className="bg-red-200">
        <span>scatter</span>
      </div>
      <div className="bg-red-200">
        <span>scatter</span>
      </div>
      <div style={{ flex: "1 1 auto" }}>
        <AutoSizer disableWidth>
          {({ height }) => (
            <Panels
              type="3col"
              left={<Disks height={height} />}
              middle={<FileSystem height={height} />}
              right={<Targets height={height} />}
            />
          )}
        </AutoSizer>
      </div>
    </div>
  )
}
