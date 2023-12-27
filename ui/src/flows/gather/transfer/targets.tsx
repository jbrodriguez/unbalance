import React from 'react';

import AutoSizer from 'react-virtualized-auto-sizer';

export const Targets: React.FunctionComponent = () => {
  return (
    <div style={{ flex: '1 1 auto' }}>
      <AutoSizer disableWidth>
        {({ height }) => (
          <div className="flex flex-1 flex-col bg-neutral-200 dark:bg-gray-950">
            <div
              className={`overflow-y-auto`}
              style={{ height: `${height}px` }}
            >
              <div className="relative overflow-x-auto">
                <table className="w-full text-sm text-left rtl:text-right text-gray-500 dark:text-gray-400">
                  <thead className="text-xs text-gray-700 uppercase bg-gray-50 dark:bg-gray-700 dark:text-gray-400">
                    <tr>
                      <th scope="col" className="p-4">
                        <div className="flex items-center">
                          <input
                            id="checkbox-all-search"
                            type="checkbox"
                            className="w-4 h-4 text-blue-600 bg-gray-100 border-gray-300 rounded focus:ring-blue-500 dark:focus:ring-blue-600 dark:ring-offset-gray-800 dark:focus:ring-offset-gray-800 focus:ring-2 dark:bg-gray-700 dark:border-gray-600"
                          />
                          <label
                            htmlFor="checkbox-all-search"
                            className="sr-only"
                          >
                            checkbox
                          </label>
                        </div>
                      </th>
                      <th scope="col" className="px-6 py-3">
                        DISK
                      </th>
                      <th scope="col" className="px-6 py-3">
                        TYPE
                      </th>
                      <th scope="col" className="px-6 py-3">
                        SERIAL
                      </th>
                      <th scope="col" className="px-6 py-3">
                        TRANSFER
                      </th>
                      <th scope="col" className="px-6 py-3">
                        SIZE
                      </th>
                      <th scope="col" className="px-6 py-3">
                        CURRENT
                      </th>
                      <th scope="col" className="px-6 py-3">
                        PLANNED
                      </th>
                    </tr>
                  </thead>
                  <tbody>
                    <tr className="bg-white border-b dark:bg-gray-800 dark:border-gray-700 hover:bg-gray-50 dark:hover:bg-gray-600">
                      <td className="w-4 p-4">
                        <div className="flex items-center">
                          <input
                            id="checkbox-table-search-1"
                            type="checkbox"
                            className="w-4 h-4 text-blue-600 bg-gray-100 border-gray-300 rounded focus:ring-blue-500 dark:focus:ring-blue-600 dark:ring-offset-gray-800 dark:focus:ring-offset-gray-800 focus:ring-2 dark:bg-gray-700 dark:border-gray-600"
                          />
                          <label
                            htmlFor="checkbox-table-search-1"
                            className="sr-only"
                          >
                            checkbox
                          </label>
                        </div>
                      </td>
                      <th
                        scope="row"
                        className="px-6 py-4 font-medium text-gray-900 whitespace-nowrap dark:text-white"
                      >
                        disk1
                      </th>
                      <td className="px-6 py-4">xfs</td>
                      <td className="px-6 py-4">
                        TOSHIBA_DT01ACA200_7322ZM4KS (sdb)
                      </td>
                      <td className="px-6 py-4">90.4 GB</td>
                      <td className="px-6 py-4">2 TB</td>
                      <td className="px-6 py-4">107 GB</td>
                      <td className="px-6 py-4">16.3 GB</td>
                    </tr>
                  </tbody>
                </table>
              </div>
            </div>
          </div>
        )}
      </AutoSizer>
    </div>
  );
};
