import React from 'react';

import { NavLink } from 'react-router-dom';

interface Props {
  height?: number;
}

export const Sidebar: React.FunctionComponent<Props> = ({ height }) => {
  return (
    <div
      className="bg-neutral-200 dark:bg-gray-950"
      style={{ height: `${height}px` }}
    >
      <div className="overflow-y-auto" style={{ height: `${height}px` }}>
        <aside
          className="overflow-y-auto bg-gray-50 dark:bg-gray-900"
          style={{ height: `${height}px` }}
        >
          <ul className="space-y-2 font-medium">
            <li>
              <NavLink
                to="validation"
                className={({ isActive }) =>
                  isActive
                    ? 'flex items-center p-2 text-gray-900 rounded-l-lg dark:text-white bg-neutral-200 dark:bg-gray-950 border-r-2 border-gray-50 dark:border-gray-900'
                    : 'flex items-center p-2 text-gray-900 rounded-l-lg dark:text-white hover:bg-gray-100 dark:hover:bg-gray-700'
                }
              >
                <svg
                  className="w-5 h-5 text-gray-500 transition duration-75 dark:text-gray-400 group-hover:text-gray-900 dark:group-hover:text-white"
                  aria-hidden="true"
                  xmlns="http://www.w3.org/2000/svg"
                  fill="currentColor"
                  viewBox="0 0 22 21"
                >
                  <path d="M16.975 11H10V4.025a1 1 0 0 0-1.066-.998 8.5 8.5 0 1 0 9.039 9.039.999.999 0 0 0-1-1.066h.002Z" />
                  <path d="M12.5 0c-.157 0-.311.01-.565.027A1 1 0 0 0 11 1.02V10h8.975a1 1 0 0 0 1-.935c.013-.188.028-.374.028-.565A8.51 8.51 0 0 0 12.5 0Z" />
                </svg>
                <span className="ms-3">Validation</span>
              </NavLink>
            </li>
            <li>
              <NavLink
                to="log"
                className={({ isActive }) =>
                  isActive
                    ? 'flex items-center p-2 text-gray-900 rounded-l-lg dark:text-white bg-neutral-200 dark:bg-gray-950 border-r-2 border-gray-50 dark:border-gray-900'
                    : 'flex items-center p-2 text-gray-900 rounded-l-lg dark:text-white hover:bg-gray-100 dark:hover:bg-gray-700'
                }
              >
                <svg
                  className="flex-shrink-0 w-5 h-5 text-gray-500 transition duration-75 dark:text-gray-400 group-hover:text-gray-900 dark:group-hover:text-white"
                  aria-hidden="true"
                  xmlns="http://www.w3.org/2000/svg"
                  fill="currentColor"
                  viewBox="0 0 18 18"
                >
                  <path d="M6.143 0H1.857A1.857 1.857 0 0 0 0 1.857v4.286C0 7.169.831 8 1.857 8h4.286A1.857 1.857 0 0 0 8 6.143V1.857A1.857 1.857 0 0 0 6.143 0Zm10 0h-4.286A1.857 1.857 0 0 0 10 1.857v4.286C10 7.169 10.831 8 11.857 8h4.286A1.857 1.857 0 0 0 18 6.143V1.857A1.857 1.857 0 0 0 16.143 0Zm-10 10H1.857A1.857 1.857 0 0 0 0 11.857v4.286C0 17.169.831 18 1.857 18h4.286A1.857 1.857 0 0 0 8 16.143v-4.286A1.857 1.857 0 0 0 6.143 10Zm10 0h-4.286A1.857 1.857 0 0 0 10 11.857v4.286c0 1.026.831 1.857 1.857 1.857h4.286A1.857 1.857 0 0 0 18 16.143v-4.286A1.857 1.857 0 0 0 16.143 10Z" />
                </svg>
                <span className="flex-1 ms-3 whitespace-nowrap">Logs</span>
              </NavLink>
            </li>
          </ul>
        </aside>
      </div>
    </div>
  );
};
