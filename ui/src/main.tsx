import React from 'react';
import ReactDOM from 'react-dom/client';

import {
  createBrowserRouter,
  RouterProvider,
  Navigate,
} from 'react-router-dom';

import App from './App.tsx';
import './index.css';
import { Scatter } from './flows/scatter/scatter.tsx';
import { Gather } from './flows/gather/gather.tsx';

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
    ],
  },
]);

ReactDOM.createRoot(document.getElementById('root')!).render(
  <React.StrictMode>
    <RouterProvider router={router} />
  </React.StrictMode>,
);
