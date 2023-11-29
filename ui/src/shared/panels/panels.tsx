import React from "react"

import { Panel, PanelGroup } from "react-resizable-panels"
import { ResizeHandle } from "~/shared/resize/resize-handle"
import styles from "~/shared/resize/shared.module.css"

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
    <div className={styles.PanelGroupWrapper}>
      <PanelGroup className={styles.PanelGroup} direction="horizontal">
        <Panel
          className={styles.PanelRow}
          defaultSizePercentage={30}
          minSizePercentage={20}
        >
          {left}
        </Panel>
        <ResizeHandle className={styles.ResizeHandle} />
        <Panel className={styles.PanelRow} minSizePercentage={30}>
          {middle}
        </Panel>
        {type === "3col" && (
          <>
            <ResizeHandle className={styles.ResizeHandle} />
            <Panel
              className={styles.PanelRow}
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
