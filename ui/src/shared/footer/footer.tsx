import React from 'react';

import { Icon } from '~/shared/icons/icon';
import diskmv from '~/assets/diskmv.png';
import unraid from '~/assets/unraid.png';
import jb from '~/assets/jb.png';
import { useConfigVersion } from '~/state/config';

export const Footer: React.FunctionComponent = () => {
  const version = useConfigVersion();

  return (
    <section className="grid grid-cols-2 gap-2 md:flex md:flex-row md:items-center md:justify-between bg-gray-800 dark:bg-gray-800 text-sky-700 dark:text-slate-300 p-2 my-4">
      <div className="md:order-1 flex flex-col gap-1 text-xs md:text-base">
        <div className="text-lime-600">
          <span>unbalanced</span>
          {version !== '' && <span>&nbsp;{version}</span>}
        </div>
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
      </div>
      <div className="md:order-2 flex flex-wrap gap-1 md:gap-2 items-center justify-end md:justify-end">
        <a
          href="https://jbrio.net/unbalanced"
          title="Support Fund"
          rel="noreferrer noopener"
          target="_blank"
        >
          <Icon name="gift" size={20} style="fill-lime-600" />
        </a>

        <a
          href="https://x.com/jbrodriguezio"
          title="@jbrodriguezio"
          rel="noreferrer noopener"
          target="_blank"
        >
          <Icon name="x" size={18} style="fill-neutral-300" />
        </a>

        <a
          href="https://github.com/jbrodriguez"
          title="github.com/jbrodriguez"
          rel="noreferrer noopener"
          target="_blank"
        >
          <Icon name="github" size={20} style="fill-neutral-300" />
        </a>

        <a
          href="https://forums.unraid.net/topic/34547-diskmv-a-set-of-utilities-to-move-files-between-disks/"
          title="diskmv"
          rel="noreferrer noopener"
          target="_blank"
        >
          <img src={diskmv} alt="logo" className="h-8 md:h-10" />
        </a>

        <a
          className=""
          href="https://unraid.net/"
          title="unraid.net"
          rel="noreferrer noopener"
          target="_blank"
        >
          <img src={unraid} alt="logo" className="h-6 md:h-8" />
        </a>

        <a
          className=""
          href="https://jbrio.net/"
          title="jbrio.net"
          rel="noreferrer noopener"
          target="_blank"
        >
          <img src={jb} alt="logo" className="h-6 md:h-8" />
        </a>
      </div>
    </section>
  );
};
