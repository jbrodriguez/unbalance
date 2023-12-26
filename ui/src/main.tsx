import React from 'react';
import ReactDOM from 'react-dom/client';

import {
  createBrowserRouter,
  RouterProvider,
  Navigate,
} from 'react-router-dom';

import { App } from './App.tsx';
import './index.css';
import { Scatter } from '~/flows/scatter/scatter';
import { Select as ScatterSelect } from '~/flows/scatter/select/select';
import { Plan as ScatterPlan } from '~/flows/scatter/plan/plan';
import { Transfer } from '~/flows/scatter/transfer/transfer';
import { Validation as ScatterValidation } from '~/flows/scatter/transfer/validation';
// import { Log as ScatterLog } from '~/shared/log/log';
import { Gather } from '~/flows/gather/gather';
import { Select as GatherSelect } from '~/flows/gather/select/select';
import { Plan as GatherPlan } from '~/flows/gather/plan/plan';
import { History } from '~/flows/history/history';
import { Settings } from '~/flows/settings/settings';
import { Notifications } from '~/flows/settings/notifications';
import { Logs } from '~/flows/logs/logs';
import { Transfer as SharedTransfer } from '~/shared/transfer/transfer';

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
            element: <ScatterPlan />,
          },
          {
            path: 'transfer',
            element: <Transfer />,
            children: [
              {
                path: 'validation',
                element: <ScatterValidation />,
              },
              {
                path: 'operation',
                element: <SharedTransfer />,
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
            element: <GatherPlan />,
          },
          {
            path: 'transfer',
            element: <Transfer />,
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
