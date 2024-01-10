import React from 'react';

import { NavLink, Outlet } from 'react-router-dom';
// import AutoSizer from 'react-virtualized-auto-sizer';

import { Icon } from '~/shared/icons/icon';

export const Settings: React.FunctionComponent = () => {
  return (
    <div className="grid grid-cols-12 gap-2 h-full ">
      <aside className="col-span-2 bg-gray-50 dark:bg-gray-800">
        <div className="px-3 py-4">
          <ul className="space-y-2 font-medium">
            <li>
              <NavLink
                to="notifications"
                className={({ isActive }) =>
                  isActive
                    ? 'flex items-center p-2 text-gray-900 rounded-lg dark:text-white bg-blue-700 dark:bg-blue-700 group'
                    : 'flex items-center p-2 text-gray-900 rounded-lg dark:text-white hover:bg-gray-100 dark:hover:bg-gray-700 group'
                }
                // className="flex items-center p-2 text-gray-900 rounded-lg dark:text-white hover:bg-gray-100 dark:hover:bg-gray-700 group"
              >
                <Icon name="bell" size={24} style="fill-gray-500" />
                <span className="pr-3" />
                Notifications
              </NavLink>
            </li>
            <li>
              <NavLink
                to="reserved"
                className={({ isActive }) =>
                  isActive
                    ? 'flex items-center p-2 text-gray-900 rounded-lg dark:text-white bg-blue-700 dark:bg-blue-700 group'
                    : 'flex items-center p-2 text-gray-900 rounded-lg dark:text-white hover:bg-gray-100 dark:hover:bg-gray-700 group'
                }
                // className="flex items-center p-2 text-gray-900 rounded-lg dark:text-white hover:bg-gray-100 dark:hover:bg-gray-700 group"
              >
                <Icon name="lifebuoy" size={24} style="fill-gray-500" />
                <span className="pr-3" />
                Reserved Space
              </NavLink>
            </li>
            <li>
              <NavLink
                to="flags"
                className={({ isActive }) =>
                  isActive
                    ? 'flex items-center p-2 text-gray-900 rounded-lg dark:text-white bg-blue-700 dark:bg-blue-700 group'
                    : 'flex items-center p-2 text-gray-900 rounded-lg dark:text-white hover:bg-gray-100 dark:hover:bg-gray-700 group'
                }
                // className="flex items-center p-2 text-gray-900 rounded-lg dark:text-white hover:bg-gray-100 dark:hover:bg-gray-700 group"
              >
                <Icon name="flag" size={24} style="fill-gray-500" />
                <span className="pr-3" />
                Rsync Flags
              </NavLink>
            </li>
            <li>
              <NavLink
                to="verbosity"
                className={({ isActive }) =>
                  isActive
                    ? 'flex items-center p-2 text-gray-900 rounded-lg dark:text-white bg-blue-700 dark:bg-blue-700 group'
                    : 'flex items-center p-2 text-gray-900 rounded-lg dark:text-white hover:bg-gray-100 dark:hover:bg-gray-700 group'
                }
                // className="flex items-center p-2 text-gray-900 rounded-lg dark:text-white hover:bg-gray-100 dark:hover:bg-gray-700 group"
              >
                <Icon name="log" size={24} style="fill-gray-500" />
                <span className="pr-3" />
                Log Verbosity
              </NavLink>
            </li>
          </ul>
          <div className="pt-6" />
          <div className="p-4 rounded-lg border border-gray-400 dark:border-gray-700">
            <div className="flex items-center mb-3">
              <span className="bg-green-500 text-green-100 text-sm font-semibold px-2.5 py-0.5 rounded dark:bg-green-800 dark:text-green-200">
                Promo
              </span>
            </div>
            <p className="mb-3 text-sm">
              Sponsor continuous development of the plugin and receive heartfelt
              thanks from the developer ! 😀 🙌
            </p>
            <a
              className="text-sm text-blue-800 underline font-medium hover:text-blue-900 dark:text-blue-400 dark:hover:text-blue-300"
              rel="noreferrer noopener"
              target="_blank"
              href="https://jbrio.net/unbalanced"
            >
              Contribute
            </a>
          </div>
        </div>
      </aside>
      <div className="col-span-10 bg-neutral-100 dark:bg-gray-950">
        <Outlet />
      </div>
    </div>
  );
};
