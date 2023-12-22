import React from 'react';

interface Props {
  children?: React.ReactNode;
}

export const Pane: React.FunctionComponent<Props> = ({ children }) => {
  return (
    <div className="border border-solid dark:border-slate-700 rounded p-2 mb-4">
      {children}
    </div>
  );
};
