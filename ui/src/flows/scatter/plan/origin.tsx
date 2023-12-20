import React from 'react';

interface Props {
  height?: number;
}

export const Origin: React.FunctionComponent<Props> = ({ height = 0 }) => {
  return (
    <div className="flex flex-1 flex-col bg-neutral-200 dark:bg-gray-950">
      <div className="p-2 overflow-y-auto" style={{ height: `${height}px` }}>
        <h1>origin</h1>
      </div>
    </div>
  );
};
