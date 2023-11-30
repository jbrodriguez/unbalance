import React, { PropsWithChildren } from 'react';

export const Pane: React.FunctionComponent<PropsWithChildren> = ({
  children,
}) => {
  return (
    <div className="border border-solid dark:border-slate-700 rounded p-2 mb-4">
      {children}
    </div>
  );
};
