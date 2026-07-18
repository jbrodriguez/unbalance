import React from 'react';

import { Label } from '@/components/ui/label';
import { RadioGroup, RadioGroupItem } from '@/components/ui/radio-group';

import { useConfigActions, useConfigLogLines } from '~/state/config';

const choices = [100, 500, 1000, 5000];

export const Logging: React.FunctionComponent = () => {
  const logLines = useConfigLogLines();
  const { setLogLines } = useConfigActions();

  const onChange = (value: string) => setLogLines(+value);

  return (
    <div className="p-4">
      <h1>
        Controls how many lines of the log file are loaded into the Log page.
        <br />
        Higher values make it easier to troubleshoot longer operations, at the
        cost of a heavier page.
        <br />
        The log file itself is rotated on disk regardless of this setting.
      </h1>
      <div className="pb-4" />

      <h2 className="text-lg font-bold">Log lines</h2>
      <div className="pb-1" />
      <p className="text-sm text-gray-500 dark:text-gray-500 pb-2">
        currently loading the last {logLines} lines
      </p>
      <RadioGroup
        defaultValue={`${logLines}`}
        value={`${logLines}`}
        onValueChange={onChange}
      >
        {choices.map((choice) => (
          <div key={choice} className="flex items-center space-x-2">
            <RadioGroupItem value={`${choice}`} id={`log-lines-${choice}`} />
            <Label htmlFor={`log-lines-${choice}`}>{choice} lines</Label>
          </div>
        ))}
      </RadioGroup>
    </div>
  );
};
