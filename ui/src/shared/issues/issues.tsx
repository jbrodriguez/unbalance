import React from 'react';

import AutoSizer from 'react-virtualized-auto-sizer';
import dayjs from 'dayjs';

import { useUnraidPlan } from '~/state/unraid';
import { formatTime } from '~/helpers/units';

export const Issues: React.FunctionComponent = () => {
  const plan = useUnraidPlan();

  const owner = plan?.ownerIssue ?? 'n/a';
  const group = plan?.groupIssue ?? 'n/a';
  const folder = plan?.folderIssue ?? 'n/a';
  const file = plan?.fileIssue ?? 'n/a';

  const elapsed = plan
    ? dayjs(plan.finished).diff(dayjs(plan.started), 'second')
    : 0;

  const showWarning =
    plan &&
    plan.ownerIssue + plan.groupIssue + plan.folderIssue + plan.fileIssue > 0;

  return (
    <AutoSizer disableWidth>
      {({ height }) => (
        <div className="flex flex-1 flex-col bg-neutral-100 dark:bg-gray-950">
          <div
            className="overflow-y-auto p-2 text-base text-gray-700 dark:text-gray-500"
            style={{ height: `${height}px` }}
          >
            <div className="flex flex-row items-center justify-between">
              <h2>ISSUES</h2>
              {plan && (
                <span className="font-bold">
                  runtime: {formatTime(elapsed)}{' '}
                </span>
              )}
            </div>
            <div className="pb-2" />
            <section>
              <div className="flex flex-row items-center justify-between text-2xl text-gray-900 dark:text-gray-400">
                <span>Owner</span>
                <span className="font-bold">{owner}</span>
              </div>
              <span className="text-sm">
                file(s)/folder(s) with an owner other than 'nobody'
              </span>
            </section>
            <div className="pb-4" />
            <section>
              <div className="flex flex-row items-center justify-between text-2xl text-gray-900 dark:text-gray-400">
                <span>Group</span>
                <span className="font-bold">{group}</span>
              </div>
              <span className="text-sm">
                file(s)/folder(s) with a group other than 'users'
              </span>
            </section>
            <div className="pb-4" />
            <section>
              <div className="flex flex-row items-center justify-between text-2xl text-gray-900 dark:text-gray-400">
                <span>Folder permissions</span>
                <span className="font-bold">{folder}</span>
              </div>
              <span className="text-sm">
                folder(s) with a permission other than 'drwxrwxrwx'
              </span>
            </section>
            <div className="pb-4" />
            <section>
              <div className="flex flex-row items-center justify-between text-2xl text-gray-900 dark:text-gray-400">
                <span>File permissions</span>
                <span className="font-bold">{file}</span>
              </div>
              <span className="text-sm">
                files(s) with a permission other than '-rw-rw-rw-' or
                '-r--r--r--'
              </span>
            </section>
            {showWarning && (
              <>
                <div className="pb-4" />
                <section>
                  You can find more details about which files have issues in the
                  log file (/var/log/unbalance.log). <br />
                  At this point, you can transfer the folders/files if you want,
                  but be advised that it can cause errors in the operation.{' '}
                  <br />
                  You are suggested to install the Fix Common Problems plugin,
                  then run the Docker Safe New Permissions command.
                </section>
              </>
            )}
          </div>
        </div>
      )}
    </AutoSizer>
    // </div>
  );
};
