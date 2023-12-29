import React from 'react';

// import AutoSizer from 'react-virtualized-auto-sizer';

import { Dashboard } from './dashboard';
import { Commands } from './commands';

export const Transfer: React.FunctionComponent = () => {
  return (
    <div className="flex flex-1 flex-col">
      <Dashboard />
      <div className="pb-2" />
      <Commands />
    </div>
  );
};

// export const Transfer: React.FunctionComponent = () => {
//   return (
//     <div className="h-full flex flex-col bg-neutral-100 dark:bg-gray-950">
//       <div className="p-2">
//         <Dashboard />
//       </div>
//       <div className="flex-auto">
//         <AutoSizer disableWidth>
//           {({ height }) => (
//             <div
//               className="overflow-y-auto px-2 pb-2"
//               style={{ height: `${height}px` }}
//             >
//               <Commands />
//             </div>
//           )}
//         </AutoSizer>
//       </div>
//     </div>
//   );
// };
