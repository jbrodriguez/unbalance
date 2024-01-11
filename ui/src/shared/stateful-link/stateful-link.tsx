import React from 'react';

import { NavLink } from 'react-router-dom';

interface StatefulLinkProps {
  to: string;
  className: (props: {
    isActive?: boolean;
    isPending?: boolean;
    isTransitioning?: boolean;
  }) => string;
  disabled?: boolean;
  children: React.ReactNode;
}

export const StatefulLink: React.FunctionComponent<StatefulLinkProps> = ({
  disabled = false,
  children,
  ...rest
}) =>
  disabled ? (
    <span
      className="ml-4 cursor-not-allowed text-sky-700 dark:text-slate-400"
      aria-disabled="true"
    >
      {children}
    </span>
  ) : (
    <NavLink {...rest}>{children}</NavLink>
  );
