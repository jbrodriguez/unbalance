import React from 'react';

import { Checkbox } from '~/shared/checkbox/checkbox';

import { useConfigActions, useConfigPreserveHardlinks } from '~/state/config';

export const Hardlinks: React.FunctionComponent = () => {
  const preserveHardlinks = useConfigPreserveHardlinks();
  const { setPreserveHardlinks } = useConfigActions();

  const onToggle = () => {
    setPreserveHardlinks(!preserveHardlinks);
  };

  return (
    <div className="p-4">
      <h1 className="text-xl font-semibold mb-4">Hardlink Preservation</h1>
      
      <div className="mb-6">
        <p className="mb-3">
          Hardlinks are multiple directory entries that point to the same file data on disk.
          This setting controls how <span className="text-lime-600 dark:text-lime-700">unbalanced</span> handles hardlinks during transfers.
        </p>
        
        <div className="bg-blue-50 dark:bg-blue-900 border border-blue-200 dark:border-blue-700 rounded-lg p-4 mb-4">
          <h3 className="font-semibold text-blue-800 dark:text-blue-200 mb-2">
            How it works:
          </h3>
          <ul className="text-blue-700 dark:text-blue-300 space-y-1">
            <li>• <strong>Enabled:</strong> Adds -H flag to rsync to preserve hardlinks (saves space)</li>
            <li>• <strong>Disabled:</strong> Hardlinks become separate files (uses more space but safer)</li>
          </ul>
        </div>

        <div className="bg-amber-50 dark:bg-amber-900 border border-amber-200 dark:border-amber-700 rounded-lg p-4 mb-4">
          <h3 className="font-semibold text-amber-800 dark:text-amber-200 mb-2">
            Planning Impact:
          </h3>
          <p className="text-amber-700 dark:text-amber-300">
            When disabled, unbalanced uses conservative space calculations assuming hardlinks will become separate files.
            This prevents "insufficient space" errors during transfers.
          </p>
        </div>
      </div>

      <div className="flex flex-row items-center gap-3 p-4 bg-gray-50 dark:bg-gray-800 rounded-lg">
        <Checkbox checked={preserveHardlinks} onCheck={onToggle} />
        <label className="text-sm font-medium">
          Preserve hardlinks during file transfers
          {preserveHardlinks ? (
            <span className="block text-xs text-green-600 dark:text-green-400">
              Hardlinks will be preserved (rsync -H flag enabled)
            </span>
          ) : (
            <span className="block text-xs text-amber-600 dark:text-amber-400">
              Hardlinks will become separate files (conservative space planning)
            </span>
          )}
        </label>
      </div>

      <div className="mt-6 text-sm text-gray-600 dark:text-gray-400">
        <p>
          <strong>Note:</strong> This setting affects space calculations during the planning phase and rsync behavior during transfers.
          Changes take effect immediately for new operations.
        </p>
      </div>
    </div>
  );
};