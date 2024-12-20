import React from 'react';
import { Checkbox } from '~/shared/checkbox/checkbox';

interface ToggleProps {
  allChecked: boolean;
  onCheck: () => void;
}

export const Toggle: React.FunctionComponent<ToggleProps> = ({
  allChecked,
  onCheck,
}) => {
  return (
    <div className="flex flex-row gap-6">
      <Checkbox checked={allChecked} onCheck={onCheck} />
      <span className="text-slate-500 dark:text-gray-500">Select All</span>
    </div>
  );
};
