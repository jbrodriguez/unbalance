import React from 'react';

import AutoSizer from 'react-virtualized-auto-sizer';
import dayjs from 'dayjs';
import utc from 'dayjs/plugin/utc';
dayjs.extend(utc);

import { useUnraidHistory } from '~/state/unraid';
import { Operation } from '~/types';
import { Selectable } from '~/shared/selectable/selectable';
import { formatBytes } from '~/helpers/units';
import { operationKindToName } from '~/helpers/operation';

interface Props {
  current: Operation | null;
  onSelected: (operation: Operation, first: boolean) => void;
}

export const Operations: React.FunctionComponent<Props> = ({
  current,
  onSelected,
}) => {
  const operations = useUnraidHistory();
  const onClick = (op: Operation, first: boolean) => () => {
    console.log('clicked', op);
    onSelected(op, first);
  };

  console.log('operations', operations);

  return (
    <div className="h-full flex flex-col bg-neutral-100 dark:bg-gray-950">
      <div className="flex-auto">
        <AutoSizer disableWidth>
          {({ height }) => (
            <div
              className="overflow-y-auto overflow-x-auto p-2"
              style={{ height: `${height}px` }}
            >
              {operations.map((operation, index) => {
                const name = operationKindToName[operation.opKind];
                const finished = dayjs(operation.finished);
                const { value, unit } = formatBytes(operation.bytesTransferred);

                return (
                  <Selectable
                    key={operation.id}
                    onClick={onClick(operation, index === 0)}
                    selected={current?.id === operation.id}
                  >
                    <div className="flex flex-row items-start justify-between">
                      <div>
                        <div className="text-lg">{name}</div>
                        <div className="text-slate-500 dark:text-gray-500">
                          {finished.local().format('YYYY.MM.DD, HH:mm')}
                        </div>
                      </div>
                      <div className="flex flex-row items-center">
                        <div className="text-slate-900 dark:text-gray-100 text-2xl">
                          {value}
                        </div>
                        <div className="pr-1" />
                        <div className="text-slate-500 dark:text-gray-500 text-sm">
                          {unit}
                        </div>
                      </div>
                    </div>
                  </Selectable>
                );
              })}
            </div>
          )}
        </AutoSizer>
      </div>
    </div>
  );
};
