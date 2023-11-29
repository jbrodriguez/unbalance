import React from "react"

import ReactResizeDetector from "react-resize-detector"

// import { Panel, PanelGroup } from "react-resizable-panels"
// import { ResizeHandle } from "~/shared/resize/resize-handle"
// import styles from "~/shared/resize/shared.module.css"
import { Panels } from "~/shared/panels/panels"
import { Disks } from "./disks"
import { FileSystem } from "./filesystem"
import { Targets } from "./targets"

export const Scatter: React.FC = () => {
  return (
    <ReactResizeDetector handleHeight>
      {({ height }) => (
        <Panels
          type="3col"
          left={<Disks height={height} />}
          middle={<FileSystem height={height} />}
          right={<Targets height={height} />}
        />
      )}
    </ReactResizeDetector>
  )
  // return (
  //   <div className={styles.PanelGroupWrapper}>
  //     <PanelGroup className={styles.PanelGroup} direction="horizontal">
  //       <Panel
  //         className={styles.PanelRow}
  //         defaultSizePercentage={30}
  //         minSizePercentage={20}
  //       >
  //         <div className={styles.Centered}>left</div>
  //       </Panel>
  //       <ResizeHandle className={styles.ResizeHandle} />
  //       <Panel className={styles.PanelRow} minSizePercentage={30}>
  //         <div className={styles.Centered}>middle</div>
  //       </Panel>
  //       <ResizeHandle className={styles.ResizeHandle} />
  //       <Panel
  //         className={styles.PanelRow}
  //         defaultSizePercentage={30}
  //         minSizePercentage={20}
  //       >
  //         <div className={styles.Centered}>right</div>
  //       </Panel>
  //     </PanelGroup>
  //   </div>
  // )
}
