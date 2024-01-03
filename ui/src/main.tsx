import React from 'react';
import ReactDOM from 'react-dom/client';

import {
  createBrowserRouter,
  RouterProvider,
  Navigate,
  Outlet,
} from 'react-router-dom';

import { App } from './App.tsx';
import './index.css';
import { Scatter } from '~/flows/scatter/scatter';
import { Select as ScatterSelect } from '~/flows/scatter/select/select';
import { Validation as ScatterValidation } from '~/flows/scatter/transfer/validation';
import { Gather } from '~/flows/gather/gather';
import { Select as GatherSelect } from '~/flows/gather/select/select';
import { Targets } from '~/flows/gather/transfer/targets';
import { History } from '~/flows/history/history';
import { Settings } from '~/flows/settings/settings';
import { Notifications } from '~/flows/settings/notifications';
import { Reserved } from '~/flows/settings/reserved';
import { Logs } from '~/flows/logs/logs';
import { Transfer } from '~/shared/transfer/transfer';
import { Feedback } from '~/shared/feedback/feedback';

const router = createBrowserRouter([
  {
    path: '/',
    element: <App />,
    children: [
      { index: true, element: <Navigate to="/scatter" replace /> },
      {
        path: '/scatter',
        element: <Scatter />,
        children: [
          { index: true, element: <Navigate to="/scatter/select" replace /> },
          {
            path: 'select',
            element: <ScatterSelect />,
          },
          {
            path: 'plan',
            element: <Feedback />,
          },
          {
            path: 'transfer',
            element: <Outlet />,
            children: [
              {
                path: 'validation',
                element: <ScatterValidation />,
              },
              {
                path: 'operation',
                element: <Transfer />,
              },
            ],
          },
        ],
      },
      {
        path: '/gather',
        element: <Gather />,
        children: [
          { index: true, element: <Navigate to="/gather/select" replace /> },
          {
            path: 'select',
            element: <GatherSelect />,
          },
          {
            path: 'plan',
            element: <Feedback />,
          },
          {
            path: 'transfer',
            element: <Outlet />,
            children: [
              {
                path: 'targets',
                element: <Targets />,
              },
              {
                path: 'operation',
                element: <Transfer />,
              },
            ],
          },
        ],
      },
      {
        path: '/history',
        element: <History />,
      },
      {
        path: '/settings',
        element: <Settings />,
        children: [
          {
            path: 'notifications',
            element: <Notifications />,
          },
          {
            path: 'reserved',
            element: <Reserved />,
          },
        ],
      },
      {
        path: '/logs',
        element: <Logs />,
      },
    ],
  },
]);

ReactDOM.createRoot(document.getElementById('root')!).render(
  <React.StrictMode>
    <RouterProvider router={router} />
  </React.StrictMode>,
);
