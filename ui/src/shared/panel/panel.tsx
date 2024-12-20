import React from 'react';

import AutoSizer from 'react-virtualized-auto-sizer';

interface PanelProps {
  title?: string;
  children: React.ReactNode;
  scrollToTop?: boolean;
  subtitle?: React.ReactNode;
}

export const Panel: React.FunctionComponent<PanelProps> = ({
  title = '',
  children,
  scrollToTop = false,
  subtitle = null,
}) => {
  const ref = React.useRef<HTMLDivElement>(null);

  React.useEffect(() => {
    if (ref.current) {
      ref.current.scrollTop = 0;
    }
  }, [scrollToTop]);

  return (
    <div className="h-full flex flex-col bg-neutral-100 dark:bg-gray-950">
      {title.length > 0 ? (
        <div className="flex flex-col pt-2 px-2">
          <div className="flex flex-row gap-6 items-center justify-start">
            <h1 className="text-lg text-slate-500 dark:text-gray-500 pb-2">
              {title}
            </h1>
            {subtitle}
          </div>
          <hr className="border-slate-300 dark:border-gray-700" />
        </div>
      ) : null}
      <div className="flex-auto">
        <AutoSizer disableWidth>
          {({ height }) => (
            <div
              className="overflow-y-auto overflow-x-auto p-2"
              style={{ height: `${height}px` }}
              ref={ref}
            >
              {children}
            </div>
          )}
        </AutoSizer>
      </div>
    </div>
  );
};
