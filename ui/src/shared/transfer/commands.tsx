import React from 'react';

import AutoSizer from 'react-virtualized-auto-sizer';

import { useUnraidOperation } from '~/state/unraid';
import { Icon } from '~/shared/icons/icon';
import { Command } from '~/shared/command/command';

export const Commands: React.FunctionComponent = () => {
  const operation = useUnraidOperation();

  if (!operation) {
    return null;
  }

  return (
    <div className="h-full flex flex-col bg-neutral-100 dark:bg-gray-950">
      <div className="flex flex-col px-2 pt-2">
        <div className="grid grid-cols-12 gap-1 items-center text-lg text-slate-500 dark:text-gray-500 pb-2">
          <div className="col-span-2 flex items-center">
            <Icon
              name="loading"
              size={14}
              style="fill-neutral-100 dark:fill-gray-950"
            />
            <span className="px-2" />
            Source
          </div>
          <div className="col-span-8">Command</div>
          <div className="col-span-2">Progress</div>
        </div>
        <hr className="border-slate-300 dark:border-gray-700" />
      </div>
      <div className="flex-auto">
        <AutoSizer disableWidth>
          {({ height }) => (
            <div
              className="px-2 pb-2 overflow-y-auto"
              style={{ height: `${height}px` }}
            >
              {operation.commands.map((command) => (
                <Command
                  key={command.id}
                  command={command}
                  rsyncStrArgs={operation.rsyncStrArgs}
                />
              ))}
            </div>
          )}
        </AutoSizer>
      </div>
    </div>
  );
};

// import React from 'react';

// import AutoSizer from 'react-virtualized-auto-sizer';

// import { useUnraidOperation } from '~/state/unraid';
// import { CommandStatus } from '~/types';
// import { Icon } from '~/shared/icons/icon';

// const getCommandStatus = (status: CommandStatus): React.ReactNode => {
//   switch (status) {
//     case CommandStatus.Complete:
//       return (
//         <Icon
//           name="check-circle"
//           size={20}
//           fill="fill-green-600 dark:fill-green-600"
//         />
//       );
//     case CommandStatus.Pending:
//       return (
//         <Icon
//           name="minus-circle"
//           size={20}
//           fill="fill-blue-600 dark:fill-blue-600"
//         />
//       );
//     case CommandStatus.Flagged:
//       return (
//         <Icon
//           name="check-circle"
//           size={20}
//           fill="fill-yellow-600 dark:fill-yellow-600"
//         />
//       );
//     case CommandStatus.Stopped:
//       return (
//         <Icon
//           name="minus-circle"
//           size={20}
//           fill="fill-red-600 dark:fill-red-600"
//         />
//       );
//     case CommandStatus.SourceRemoval:
//       return (
//         <Icon
//           name="loading"
//           size={20}
//           fill="fill-yellow-600 dark:fill-yellow-600 animate-spin"
//         />
//       );
//     default:
//       return (
//         <Icon
//           name="loading"
//           size={20}
//           fill="fill-slate-600 dark:fill-slate-600 animate-spin"
//         />
//       );
//   }
// };

// export const Commands: React.FunctionComponent = () => {
//   const operation = useUnraidOperation();

//   if (!operation) {
//     return null;
//   }

//   return (
//     // <div className="flex flex-1">
//     <table className="h-full flex flex-col text-sm text-left text-gray-500 dark:text-gray-400">
//       <thead className="text-xs text-slate-500 dark:text-gray-500 uppercase bg-neutral-100 dark:bg-gray-950">
//         <tr>
//           <th scope="col" className="p-4"></th>
//           <th scope="col" className="px-6 py-3">
//             Source
//           </th>
//           <th scope="col" className="px-6 py-3">
//             Command
//           </th>
//           <th scope="col" className="px-6 py-3">
//             Progress
//           </th>
//         </tr>
//       </thead>
//       {/* <span>test</span> */}
//       <div className="flex-auto">
//         <AutoSizer disableWidth>
//           {({ height }) => (
//             // <div
//             //   className="flex flex-1 flex-col overflow-y-auto px-2 pb-2"
//             //   style={{ height: `${height}px` }}
//             // >
//             <tbody
//               className="flex flex-1 flex-col overflow-y-auto "
//               style={{ height: `${height}px` }}
//             >
//               {operation.commands.map((command) => {
//                 const progress = (
//                   (command.transferred / command.size) *
//                   100
//                 ).toFixed(0);
//                 return (
//                   <tr className="bg-gray-300 border-b dark:bg-gray-800 dark:border-gray-700 hover:bg-gray-50 dark:hover:bg-gray-600">
//                     <td className="w-4 p-4">
//                       {getCommandStatus(command.status)}
//                     </td>
//                     <td className="px-6 py-4">{command.src}</td>
//                     <td className="px-6 py-4  font-medium text-gray-900 whitespace-nowrap dark:text-white">
//                       rsync {operation.rsyncStrArgs} &quot;{command.entry}
//                       &quot; &quot;{command.dst}&quot;
//                     </td>
//                     <td className="flex-auto px-6 py-4">
//                       <div className="w-full rounded bg-gray-400 dark:bg-gray-800">
//                         <div
//                           className="p-0.5 leading-none rounded bg-blue-900 "
//                           style={{ width: `${progress}%` }}
//                         ></div>
//                       </div>
//                     </td>
//                   </tr>
//                 );
//               })}
//             </tbody>
//             // </div>
//           )}
//         </AutoSizer>
//       </div>
//     </table>
//     // </div>
//   );
// };
