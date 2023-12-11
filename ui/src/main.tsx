import React from 'react';
import ReactDOM from 'react-dom/client';

import {
  createBrowserRouter,
  RouterProvider,
  Navigate,
} from 'react-router-dom';

import App from './App.tsx';
import './index.css';
import { Scatter } from './flows/scatter/scatter';
import { Gather } from './flows/gather/gather';
import { History } from './flows/history/history';
import { Settings } from './flows/settings/settings';
import { Notifications } from './flows/settings/notifications';

const router = createBrowserRouter([
  {
    path: '/',
    element: <App />,
    children: [
      { index: true, element: <Navigate to="/scatter" replace /> },
      {
        path: '/scatter',
        element: <Scatter />,
      },
      {
        path: '/gather',
        element: <Gather />,
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
    ],
  },
]);

ReactDOM.createRoot(document.getElementById('root')!).render(
  <React.StrictMode>
    <RouterProvider router={router} />
  </React.StrictMode>,
);
