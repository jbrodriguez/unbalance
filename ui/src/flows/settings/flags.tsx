import React from 'react';

import { Input } from '@/components/ui/input';
import { Button } from '@/components/ui/button';

import { useConfigActions, useConfigRsyncArgs } from '~/state/config';

export const Flags: React.FunctionComponent = () => {
  const flags = useConfigRsyncArgs();
  const [flagsValue, setFlagsValue] = React.useState(flags.join(' '));
  const { setRsyncArgs, resetRsyncArgs } = useConfigActions();

  const onFlagsChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    setFlagsValue(e.target.value);
  };

  const onApply = () => {
    const flags = flagsValue.split(' ');
    setRsyncArgs(flags);
  };

  const onReset = () => {
    setFlagsValue('-X');
    resetRsyncArgs();
  };

  return (
    <div className="p-4">
      <h1>
        Internally,{' '}
        <span className="text-lime-600 dark:text-lime-700">unbalanced</span>{' '}
        uses rsync to transfer files across disks. <br />
        By default, rsync is invoked with{' '}
        <span className="font-bold text-blue-900 dark:text-blue-600">
          -avPRX
        </span>{' '}
        flags. Note that the{' '}
        <span className="font-bold text-blue-900 dark:text-blue-600">X</span>{' '}
        flag is customizable, so you can remove it if needed. <br />
        You can add custom flags, except for the dry run flag which will be
        automatically added, if needed. <br />
        <span className="text-red-900 dark:text-red-700 font-bold">
          NOTE: These settings are meant to be changed by advanced users only.
          No validation is currently performed on the flags you manually enter
          here, so you need to know what you're doing before changing them.
        </span>
      </h1>
      <div className="pb-4" />

      <div className="flex flex-row items-center">
        <Input
          type="text"
          placeholder="rsync flags"
          className="w-40"
          defaultValue={flags.join(' ')}
          value={flagsValue}
          onChange={onFlagsChange}
        />
        <div className="pr-4" />
        <Button onClick={onApply} variant="secondary">
          Apply
        </Button>
        <div className="pr-4" />
        <Button onClick={onReset} variant="secondary">
          Reset to Defaults
        </Button>
      </div>
    </div>
  );
};
