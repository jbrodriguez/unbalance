import React from 'react';

import { Operations } from './operations';
import { Operation } from './operation';
import { Operation as IOperation } from '~/types';

export const History: React.FunctionComponent = () => {
  const [selected, setSelected] = React.useState<IOperation | null>(null);
  const [first, setFirst] = React.useState<boolean>(false);

  const onSelected = (operation: IOperation, first: boolean) => {
    setSelected(operation);
    setFirst(first);
  };

  return (
    <div className="grid grid-cols-12 gap-1 h-full">
      <div className="col-span-3 flex flex-col flex-1">
        <Operations current={selected} onSelected={onSelected} />
      </div>
      <div className="col-span-9 flex flex-col flex-1">
        <Operation current={selected} first={first} />
      </div>
    </div>
  );
};
