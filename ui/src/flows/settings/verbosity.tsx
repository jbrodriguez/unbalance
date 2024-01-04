import React from 'react';

import { Label } from '@/components/ui/label';
import { RadioGroup, RadioGroupItem } from '@/components/ui/radio-group';

import { useConfigActions, useConfigVerbosity } from '~/state/config';

export const Verbosity: React.FunctionComponent = () => {
  const verbosity = useConfigVerbosity();
  const { setVerbosity } = useConfigActions();

  const onChange = (value: string) => setVerbosity(+value);

  return (
    <div className="text-slate-700 dark:text-gray-300 p-4">
      <h1>
        Full verbosity will affect logging in two ways: <br />- It will print
        each line generated in the transfer (rsync) phase. <br /> - It will
        print each line generated while checking for permission issues. <br />
        Normal verbosity will not, thus greatly reducing the amount of logging.
      </h1>
      <div className="pb-4" />

      <h2 className="text-lg font-bold">Verbosity</h2>
      <div className="pb-1" />
      <RadioGroup
        defaultValue={`${verbosity}`}
        value={`${verbosity}`}
        onValueChange={onChange}
      >
        <div className="flex items-center space-x-2">
          <RadioGroupItem value="0" id="r1" />
          <Label htmlFor="r1">Normal</Label>
        </div>
        <div className="flex items-center space-x-2">
          <RadioGroupItem value="1" id="r2" />
          <Label htmlFor="r2">Full</Label>
        </div>
      </RadioGroup>
    </div>
  );
};
