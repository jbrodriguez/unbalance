import React from "react"

interface Props {
  height?: number
}

export const FileSystem: React.FC<Props> = ({ height = 0 }) => {
  return (
    <div className="flex flex-1 bg-neutral-200 dark:bg-gray-950">
      <div className={`overflow-y-auto`} style={{ height: `${height}px` }}>
        <h1>files/folders</h1>
      </div>
    </div>
  )
}
