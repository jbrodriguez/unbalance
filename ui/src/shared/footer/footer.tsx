import React from "react"

import { Icon } from "~/shared/icons/icon"
import diskmv from "~/assets/diskmv.png"
import unraid from "~/assets/unraid.png"
import jb from "~/assets/jb.png"

// import useSWR from "swr"

// import { getConfig } from "~/api"
// import { Spinner } from "~/shared/components/spinner"

export const Footer: React.FC = () => {
  // const { data, isLoading } = useSWR("/config", getConfig)

  return (
    <section className="flex flex-row items-center justify-between bg-gray-800 dark:bg-gray-800 text-sky-700 dark:text-slate-300 p-2 my-4">
      <div>
        <span className="text-slate-500 dark:text-slate-600 mr-1">
          Copyright &copy;
        </span>
        <a
          href="https://jbrio.net/"
          target="_blank"
          title="jbrio.net"
          className="text-lime-600"
        >
          Juan B. Rodriguez
        </a>
      </div>
      <div className="text-lime-600">
        <>
          <span>unbalance-ng &nbsp;</span>
          <span>v2023.11.23</span>
        </>

        {/* {isLoading ? (
          <Spinner />
        ) : (
          <>
            <span>mediaGUI &nbsp;</span>
            <span>v{data?.version}</span>
          </>
        )} */}
      </div>
      <div className="flex flex-row items-center">
        {/* <a
          className="flex items-center"
          href="https://www.themoviedb.org/"
          title="themoviedb.org"
          target="_blank"
        >
          <img src="/img/tmdb.png" alt="Logo for tmdb" className="w-10 mr-4" />
        </a>

        <a
          className="flex items-center"
          href="https://jbrio.net/"
          title="jbrio.net"
          target="_blank"
        >
          <img src="/img/logo.png" alt="Logo for jbrio.net" className="w-10" />
        </a> */}
        <a
          href="https://github.com/jbrodriguez/unbalance/blob/main/DONATIONS.md"
          title="Support Fund"
          rel="noreferrer noopener"
          target="_blank"
        >
          <Icon name="gift" size={24} fill="fill-lime-600" />
        </a>

        <a
          href="https://x.com/jbrodriguezio"
          title="@jbrodriguezio"
          rel="noreferrer noopener"
          target="_blank"
          className="ml-2"
        >
          <Icon name="x" size={20} fill="fill-neutral-300" />
        </a>

        <a
          href="https://github.com/jbrodriguez"
          title="github.com/jbrodriguez"
          rel="noreferrer noopener"
          target="_blank"
          className="ml-2"
        >
          <Icon name="github" size={24} fill="fill-neutral-300" />
        </a>

        <a
          href="https://forums.unraid.net/topic/34547-diskmv-a-set-of-utilities-to-move-files-between-disks/"
          title="diskmv"
          rel="noreferrer noopener"
          target="_blank"
        >
          <img src={diskmv} alt="logo" className="h-10" />
        </a>

        <a
          className="ml-2"
          href="https://unraid.net/"
          title="unraid.net"
          rel="noreferrer noopener"
          target="_blank"
        >
          <img src={unraid} alt="logo" className="h-8" />
        </a>

        <a
          className="ml-3"
          href="https://jbrio.net/"
          title="jbrio.net"
          rel="noreferrer noopener"
          target="_blank"
        >
          <img src={jb} alt="logo" className="h-8" />
        </a>
      </div>
    </section>
  )
}
