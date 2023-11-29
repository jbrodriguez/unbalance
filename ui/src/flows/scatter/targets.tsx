import React from "react"

interface Props {
  height?: number
}

export const Targets: React.FC<Props> = ({ height = 0 }) => {
  return (
    <div className="flex flex-1 bg-neutral-200 dark:bg-gray-950">
      <div className={`overflow-y-auto`} style={{ height: `${height}px` }}>
        {/* Your component code here */}
        <h1>target disks</h1>
      </div>
    </div>
  )
}
