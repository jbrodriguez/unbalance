import React from 'react';

import { Label } from '@/components/ui/label';
import { RadioGroup, RadioGroupItem } from '@/components/ui/radio-group';

import {
  useConfigActions,
  useConfigNotifyPlan,
  useConfigNotifyTransfer,
} from '~/state/config';

export const Notifications: React.FunctionComponent = () => {
  const notifyPlan = useConfigNotifyPlan();
  const notifyTransfer = useConfigNotifyTransfer();
  const { setNotifyPlan, setNotifyTransfer } = useConfigActions();

  const onPlanChange = (value: string) => setNotifyPlan(+value);
  const onTransferChange = (value: string) => setNotifyTransfer(+value);

  return (
    <div className="p-4">
      <h1>
        Notifications rely on Unraid's notifications settings, so you need to
        set up unRAID first, in order to receive notifications from unbalanced.
      </h1>
      <div className="pb-4" />
      <h2 className="text-lg font-bold">Planning</h2>
      <div className="pb-1" />
      <RadioGroup
        defaultValue={`${notifyPlan}`}
        value={`${notifyPlan}`}
        onValueChange={onPlanChange}
      >
        <div className="flex items-center space-x-2">
          <RadioGroupItem value="0" id="r1" />
          <Label htmlFor="r1">No Notifications</Label>
        </div>
        <div className="flex items-center space-x-2">
          <RadioGroupItem value="1" id="r2" />
          <Label htmlFor="r2">Basic</Label>
        </div>
        <div className="flex items-center space-x-2">
          <RadioGroupItem value="2" id="r3" />
          <Label htmlFor="r3">Detailed</Label>
        </div>
      </RadioGroup>

      <div className="pb-4" />

      <h2 className="text-lg font-bold">Transfer</h2>
      <div className="pb-1" />
      <RadioGroup
        defaultValue={`${notifyTransfer}`}
        value={`${notifyTransfer}`}
        onValueChange={onTransferChange}
      >
        <div className="flex items-center space-x-2">
          <RadioGroupItem
            value="0"
            id="r1"
            onClick={() => console.log('transferonot')}
          />
          <Label htmlFor="r1">No Notifications</Label>
        </div>
        <div className="flex items-center space-x-2">
          <RadioGroupItem value="1" id="r2" />
          <Label htmlFor="r2">Basic</Label>
        </div>
        <div className="flex items-center space-x-2">
          <RadioGroupItem value="2" id="r3" />
          <Label htmlFor="r3">Detailed</Label>
        </div>
      </RadioGroup>
    </div>
  );
};
