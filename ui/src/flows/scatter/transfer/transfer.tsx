import React from 'react';

import { Outlet } from 'react-router-dom';
// import AutoSizer from 'react-virtualized-auto-sizer';

// import { Panels } from '~/shared/panels/panels';

export const Transfer: React.FunctionComponent = () => {
  return (
    <Outlet />
    // <div style={{ flex: '1 1 auto' }}>
    //   <AutoSizer disableWidth>
    //     {({ height }) => (
    //       <>
    //         <Panels
    //           type="3col"
    //           left={<Disks height={height} />}
    //           middle={<FileSystem height={height} />}
    //           right={<Targets height={height} />}
    //         />
    //       </>
    //     )}
    //   </AutoSizer>
    // </div>
  );
};
