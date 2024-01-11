import React from 'react';

import logo from '~/assets/unbalance-logo.png';
import { ModeToggle } from '@/components/mode-toggle';
import { Loading } from '~/shared/icons/loading';
import { useUnraidIsBusy } from '~/state/unraid';
import { StatefulLink } from '~/shared/stateful-link/stateful-link';

export const Header: React.FunctionComponent = () => {
  const busy = useUnraidIsBusy();

  return (
    <nav className="grid grid-cols-12 gap-2 my-4">
      <ul className="col-span-2 py-2 border bg-lime-400 dark:bg-lime-600 border-lime-300 dark:border-lime-500">
        <li className="flex items-center justify-center">
          <img src={logo} alt="logo" className="h-8 mr-2" />
          <span className="text-slate-900 dark:text-slate-900 font-medium">
            unbalanced
          </span>
        </li>
      </ul>

      <ul className="col-span-10 flex flex-row items-center justify-between py-2 bg-neutral-100 dark:bg-gray-800 text-sky-700 dark:text-slate-400">
        <li>
          <StatefulLink
            to="/scatter"
            className={({ isActive }) => {
              return isActive
                ? 'ml-4 underline underline-offset-8 font-medium dark:text-white text-sky-900'
                : 'ml-4';
            }}
            disabled={busy}
          >
            SCATTER
          </StatefulLink>

          <StatefulLink
            to="/gather"
            className={({ isActive }) => {
              return isActive
                ? 'ml-4 underline underline-offset-8 font-medium dark:text-white text-sky-900'
                : 'ml-4';
            }}
            disabled={busy}
          >
            GATHER
          </StatefulLink>

          <StatefulLink
            to="/history"
            className={({ isActive }) => {
              return isActive
                ? 'ml-4 underline underline-offset-8 font-medium dark:text-white text-sky-900'
                : 'ml-4';
            }}
            disabled={busy}
          >
            HISTORY
          </StatefulLink>

          <StatefulLink
            to="/settings"
            className={({ isActive }) => {
              return isActive
                ? 'ml-4 underline underline-offset-8 font-medium dark:text-white text-sky-900'
                : 'ml-4';
            }}
            disabled={busy}
          >
            SETTINGS
          </StatefulLink>

          <StatefulLink
            to="/logs"
            className={({ isActive }) => {
              return isActive
                ? 'ml-4 underline underline-offset-8 font-medium dark:text-white text-sky-900'
                : 'ml-4';
            }}
            disabled={busy}
          >
            LOG
          </StatefulLink>
        </li>
        <li className="flex flex-row items-center">
          {busy && <Loading />}
          <ModeToggle />
          <span className="pl-2" />
        </li>
      </ul>
    </nav>
  );
};
