import React from 'react';

import { NavLink, Outlet } from 'react-router-dom';
import AutoSizer from 'react-virtualized-auto-sizer';

import { Icon } from '~/shared/icons/icon';

export const Settings: React.FunctionComponent = () => {
  return (
    <div className="flex h-full">
      <div style={{ flex: '1 1 auto' }}>
        <AutoSizer disableWidth>
          {({ height }) => (
            <div
              className="flex flex-1 bg-neutral-200 dark:bg-gray-950"
              style={{ height: `${height}px` }}
            >
              <aside className="w-64">
                <div
                  className="px-3 py-4 overflow-y-auto bg-gray-50 dark:bg-gray-800"
                  style={{ height: `${height}px` }}
                >
                  <ul className="space-y-2 font-medium">
                    <li>
                      <NavLink
                        to="notifications"
                        className="flex items-center p-2 text-gray-900 rounded-lg dark:text-white hover:bg-gray-100 dark:hover:bg-gray-700 group"
                      >
                        <Icon name="bell" size={24} style="fill-gray-500" />
                        <span className="ms-3">Notifications</span>
                      </NavLink>
                    </li>
                    <li>
                      <a
                        href="#"
                        className="flex items-center p-2 text-gray-900 rounded-lg dark:text-white hover:bg-gray-100 dark:hover:bg-gray-700 group"
                      >
                        <Icon name="lifebuoy" size={24} style="fill-gray-500" />
                        <span className="flex-1 ms-3 whitespace-nowrap">
                          Reserved Space
                        </span>
                        <span className="inline-flex items-center justify-center px-2 ms-3 text-sm font-medium text-gray-800 bg-gray-100 rounded-full dark:bg-gray-700 dark:text-gray-300">
                          Pro
                        </span>
                      </a>
                    </li>
                    <li>
                      <a
                        href="#"
                        className="flex items-center p-2 text-gray-900 rounded-lg dark:text-white hover:bg-gray-100 dark:hover:bg-gray-700 group"
                      >
                        <Icon name="flag" size={24} style="fill-gray-500" />
                        <span className="flex-1 ms-3 whitespace-nowrap">
                          Rsync Flags
                        </span>
                        <span className="inline-flex items-center justify-center w-3 h-3 p-3 ms-3 text-sm font-medium text-blue-800 bg-blue-100 rounded-full dark:bg-blue-900 dark:text-blue-300">
                          1
                        </span>
                      </a>
                    </li>
                    <li>
                      <a
                        href="#"
                        className="flex items-center p-2 text-gray-900 rounded-lg dark:text-white hover:bg-gray-100 dark:hover:bg-gray-700 group"
                      >
                        <Icon name="log" size={24} style="fill-gray-500" />
                        <span className="flex-1 ms-3 whitespace-nowrap">
                          Log Verbosity
                        </span>
                      </a>
                    </li>
                    <li>
                      <a
                        href="#"
                        className="flex items-center p-2 text-gray-900 rounded-lg dark:text-white hover:bg-gray-100 dark:hover:bg-gray-700 group"
                      >
                        <Icon name="update" size={24} style="fill-gray-500" />
                        <span className="flex-1 ms-3 whitespace-nowrap">
                          Updates
                        </span>
                      </a>
                    </li>
                  </ul>
                  <div
                    id="dropdown-cta"
                    className="p-4 mt-6 rounded-lg bg-blue-200 dark:bg-blue-900"
                    role="alert"
                  >
                    <div className="flex items-center mb-3">
                      <span className="bg-green-500 text-green-100 text-sm font-semibold me-2 px-2.5 py-0.5 rounded dark:bg-green-800 dark:text-green-200">
                        Promo
                      </span>
                      <button
                        type="button"
                        className="ms-auto -mx-1.5 -my-1.5 bg-blue-50 inline-flex justify-center items-center w-6 h-6 text-blue-900 rounded-lg focus:ring-2 focus:ring-blue-400 p-1 hover:bg-blue-200 dark:bg-blue-900 dark:text-blue-400 dark:hover:bg-blue-800"
                        data-dismiss-target="#dropdown-cta"
                        aria-label="Close"
                      >
                        <span className="sr-only">Close</span>
                        <svg
                          className="w-2.5 h-2.5"
                          aria-hidden="true"
                          xmlns="http://www.w3.org/2000/svg"
                          fill="none"
                          viewBox="0 0 14 14"
                        >
                          <path
                            stroke="currentColor"
                            stroke-linecap="round"
                            stroke-linejoin="round"
                            stroke-width="2"
                            d="m1 1 6 6m0 0 6 6M7 7l6-6M7 7l-6 6"
                          />
                        </svg>
                      </button>
                    </div>
                    <p className="mb-3 text-sm text-blue-800 dark:text-blue-400">
                      For a limited time only, donate to unbalance support fund
                      to encourage continuous development of the app ! <br />
                      😀 🙌
                    </p>
                    <a
                      className="text-sm text-blue-800 underline font-medium hover:text-blue-900 dark:text-blue-400 dark:hover:text-blue-300"
                      href="#"
                    >
                      Contribute
                    </a>
                  </div>
                </div>
              </aside>

              <div
                className={`overflow-y-auto`}
                style={{ height: `${height}px` }}
              >
                <Outlet />
              </div>
            </div>
          )}
        </AutoSizer>
      </div>
    </div>
  );
};
