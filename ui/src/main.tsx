import React from "react"
import ReactDOM from "react-dom/client"

import { createBrowserRouter, RouterProvider, Navigate } from "react-router-dom"

import App from "./App.tsx"
import "./index.css"
import { Scatter } from "./flows/scatter/scatter.tsx"

const router = createBrowserRouter([
  {
    path: "/",
    element: <App />,
    children: [
      { index: true, element: <Navigate to="/scatter" replace /> },
      {
        path: "/scatter",
        element: <Scatter />,
      },
      //   {
      //     path: "/import",
      //     element: <Import />,
      //   },
      //   {
      //     path: "/prune",
      //     element: <Prune />,
      //   },
      //   {
      //     path: "/duplicates",
      //     element: <Duplicates />,
      //   },
      //   {
      //     path: "/covers",
      //     element: <CoversScreen />,
      //   },
    ],
  },
])

ReactDOM.createRoot(document.getElementById("root")!).render(
  <React.StrictMode>
    <RouterProvider router={router} />
  </React.StrictMode>
)
