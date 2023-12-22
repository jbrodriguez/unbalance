import React from 'react';

interface Props {
  height?: number;
}

export const Operation: React.FunctionComponent<Props> = ({ height = 0 }) => {
  return (
    <div className="flex flex-1 flex-col bg-neutral-200 dark:bg-gray-950">
      <div className={`overflow-y-auto`} style={{ height: `${height}px` }}>
        <h1>operation</h1>
        <h1>operation</h1>
        <h1>operation</h1>
        <h1>operation</h1>
        <h1>operation</h1>
        <h1>operation</h1>
        <h1>operation</h1>
        <h1>operation</h1>
        <h1>operation</h1>
        <h1>operation</h1>
        <h1>operation</h1>
        <h1>operation</h1>
        <h1>operation</h1>
        <h1>operation</h1>
        <h1>operation</h1>
      </div>
    </div>
  );
};
