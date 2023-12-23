import React from 'react';

import { NavLink } from 'react-router-dom';

import logo from '~/assets/unbalance-logo.png';

export const Header: React.FunctionComponent = () => {
  return (
    <nav className="grid grid-cols-12 gap-2 my-4">
      <ul className="col-span-2 py-2 border border-lime-600 dark:border-lime-600 text-neutral-50">
        <li className="flex items-center justify-center">
          <img src={logo} alt="logo" className="h-8 mr-2" />
          <span className="dark:text-lime-600 font-medium">unbalanced</span>
        </li>
      </ul>

      <ul className="col-span-10 items-center justify-center py-2 bg-neutral-100 dark:bg-gray-800 text-sky-700 dark:text-slate-400">
        <li>
          <NavLink
            to="/scatter"
            className={({ isActive }) => {
              return isActive
                ? 'ml-4 underline underline-offset-8 font-medium dark:text-white text-sky-900'
                : 'ml-4';
            }}
          >
            SCATTER
          </NavLink>

          <NavLink
            to="/gather"
            className={({ isActive }) => {
              return isActive
                ? 'ml-4 underline underline-offset-8 font-medium dark:text-white text-sky-900'
                : 'ml-4';
            }}
          >
            GATHER
          </NavLink>

          <NavLink
            to="/history"
            className={({ isActive }) => {
              return isActive
                ? 'ml-4 underline underline-offset-8 font-medium dark:text-white text-sky-900'
                : 'ml-4';
            }}
          >
            HISTORY
          </NavLink>

          <NavLink
            to="/settings"
            className={({ isActive }) => {
              return isActive
                ? 'ml-4 underline underline-offset-8 font-medium dark:text-white text-sky-900'
                : 'ml-4';
            }}
          >
            SETTINGS
          </NavLink>

          <NavLink
            to="/logs"
            className={({ isActive }) => {
              return isActive
                ? 'ml-4 underline underline-offset-8 font-medium dark:text-white text-sky-900'
                : 'ml-4';
            }}
          >
            LOG
          </NavLink>
        </li>
      </ul>
    </nav>
  );
};
