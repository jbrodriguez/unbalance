import React from 'react';

import { Icon } from '~/shared/icons/icon';

interface Props {
  checked: boolean;
  onCheck: () => void;
}

export const Checkbox: React.FunctionComponent<Props> = ({
  checked,
  onCheck,
}) => {
  return (
    <span onClick={onCheck}>
      {checked ? (
        <Icon
          name="checked"
          size={20}
          fill="fill-slate-700 dark:fill-lime-600"
        />
      ) : (
        <Icon
          name="unchecked"
          size={20}
          fill="fill-slate-700 dark:fill-slate-200"
        />
      )}
    </span>
  );
};
