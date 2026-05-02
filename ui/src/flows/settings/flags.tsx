import React from 'react';

import { Input } from '@/components/ui/input';
import { Button } from '@/components/ui/button';

import { useConfigActions, useConfigRsyncArgs } from '~/state/config';

export const Flags: React.FunctionComponent = () => {
  const flags = useConfigRsyncArgs();
  const [flagsValue, setFlagsValue] = React.useState(flags.join(' '));
  const [error, setError] = React.useState('');
  const { setRsyncArgs, resetRsyncArgs } = useConfigActions();

  const onFlagsChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    setFlagsValue(e.target.value);
  };

  const onApply = async () => {
    const flags = flagsValue.split(/\s+/).filter(Boolean);
    try {
      await setRsyncArgs(flags);
      setError('');
    } catch (e) {
      setError(e instanceof Error ? e.message : 'Unable to save rsync flags');
    }
  };

  const onReset = async () => {
    setFlagsValue('-X');
    try {
      await resetRsyncArgs();
      setError('');
    } catch (e) {
      setError(e instanceof Error ? e.message : 'Unable to reset rsync flags');
    }
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
          Destructive or remote-execution rsync options are blocked because
          unbalanced manages source deletion and local transfer behavior itself.
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
      {error !== '' && (
        <p className="pt-3 text-sm text-red-900 dark:text-red-700">{error}</p>
      )}
    </div>
  );
};
