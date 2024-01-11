import React from 'react';

import dayjs from 'dayjs';

import { useUnraidOperation } from '~/state/unraid';
import { formatBytes, formatTime } from '~/helpers/units';
import { Gauge } from './gauge';

export const Dashboard: React.FunctionComponent = () => {
  const operation = useUnraidOperation();

  if (!operation) {
    return null;
  }

  const completion = Math.round(operation.completed * 100) / 100;
  const completed = isNaN(completion) ? '0' : completion.toFixed(2);
  const velocity = Math.round(operation.speed * 100) / 100;
  const speed = isNaN(velocity) ? '0' : velocity.toFixed(2);

  let bytes = formatBytes(operation.bytesTransferred + operation.deltaTransfer);
  const transferredValue = bytes.value;
  const transferredUnit = bytes.unit;

  bytes = formatBytes(operation.bytesToTransfer);
  const totalValue = bytes.value;
  const totalUnit = bytes.unit;

  const started = dayjs(operation.started);
  const diff = dayjs().diff(started, 'second');
  const timeElapsed = formatTime(diff);

  const elapsed = timeElapsed === '' ? 'n/a' : timeElapsed;
  const remaining =
    !operation.remaining || operation.remaining === ''
      ? 'n/a'
      : operation.remaining;

  return (
    <div className="grid grid-cols-6 gap-4 ">
      <Gauge value={completed} unit="%" label="Completed" />

      <Gauge value={speed} unit="MB/s" label="Speed" />

      <Gauge
        value={transferredValue}
        unit={transferredUnit}
        label="Transferred"
      />

      <Gauge value={totalValue} unit={totalUnit} label="Total" />

      <Gauge value={elapsed} label="Elapsed" />

      <Gauge value={remaining} label="Remaining" />
    </div>
  );
};
