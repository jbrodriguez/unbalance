import React from "react"

import { NavLink } from "react-router-dom"
// import debounce from "lodash.debounce"

// import { useOptionsStore, useOptionsActions } from "~/state/options"
// import Chevron from "~/shared/components/chevron"
import logo from "~/assets/unbalance-logo.png"

export const Header: React.FC = () => {
  // const { filterBy, filterByOptions, sortBy, sortByOptions } = useOptionsStore(
  //   (state) => ({
  //     filterBy: state.filterBy,
  //     sortBy: state.sortBy,
  //     filterByOptions: state.filterByOptions,
  //     sortByOptions: state.sortByOptions,
  //   })
  // )
  // const { setFilterBy, setSortBy, setQuery } = useOptionsActions()

  // const onFilterByChange = (e: React.ChangeEvent<HTMLSelectElement>) =>
  //   setFilterBy(e.target.value)

  // const onSortByChange = (e: React.ChangeEvent<HTMLSelectElement>) =>
  //   setSortBy(e.target.value)

  // const updateQuery = debounce((e: React.ChangeEvent<HTMLInputElement>) => {
  //   // console.log("search", e.target.value);
  //   setQuery(e.target.value)
  // }, 750)

  //   <NavLink
  //   to="/scatter"
  //   className={({ isActive }) => {
  //     return isActive
  //       ? "ml-4 bg-lime-600 dark:bg-lime-600 text-neutral-50"
  //       : "ml-4"
  //   }}
  // >
  //   SCATTER
  // </NavLink>

  return (
    <nav className="grid grid-cols-12 gap-2 my-4">
      <ul className="col-span-2 py-2 bg-lime-600 dark:bg-lime-600 text-neutral-50">
        <li className="flex items-center justify-center">
          <img src={logo} alt="logo" className="h-8 mr-2" />
          <span className="text-slate-950 font-medium">unbalance-ng</span>
        </li>
      </ul>

      <ul className="col-span-10 items-center justify-center py-2 bg-neutral-100 dark:bg-gray-800 text-sky-700 dark:text-slate-400">
        <li>
          <NavLink
            to="/scatter"
            className={({ isActive }) => {
              return isActive
                ? "ml-4 underline underline-offset-8 font-medium dark:text-white text-sky-900"
                : "ml-4"
            }}
          >
            SCATTER
          </NavLink>

          <NavLink
            to="/gather"
            className={({ isActive }) => {
              return isActive
                ? "ml-4 underline underline-offset-8 font-medium dark:text-white text-sky-900"
                : "ml-4"
            }}
          >
            GATHER
          </NavLink>

          <NavLink
            to="/history"
            className={({ isActive }) => {
              return isActive
                ? "ml-4 underline underline-offset-8 font-medium dark:text-white text-sky-900"
                : "ml-4"
            }}
          >
            HISTORY
          </NavLink>

          <NavLink
            to="/settings"
            className={({ isActive }) => {
              return isActive
                ? "ml-4 underline underline-offset-8 font-medium dark:text-white text-sky-900"
                : "ml-4"
            }}
          >
            SETTINGS
          </NavLink>

          <NavLink
            to="/log"
            className={({ isActive }) => {
              return isActive
                ? "ml-4 underline underline-offset-8 font-medium dark:text-white text-sky-900"
                : "ml-4"
            }}
          >
            LOG
          </NavLink>
        </li>
      </ul>
    </nav>
  )
}
