import React from 'react';

import { Label } from '@/components/ui/label';
import { RadioGroup, RadioGroupItem } from '@/components/ui/radio-group';

export const Notifications: React.FunctionComponent = () => {
  return (
    <div className="text-slate-700 dark:text-gray-300 p-4">
      <h1>
        Notifications rely on Unraid's notifications settings, so you need to
        set up unRAID first, in order to receive notifications from unbalanced.
      </h1>
      <div className="pb-4" />
      <h2 className="text-lg font-bold">Planning</h2>
      <div className="pb-1" />
      <RadioGroup defaultValue="planning-no-notifications">
        <div className="flex items-center space-x-2">
          <RadioGroupItem value="planning-no-notifications" id="r1" />
          <Label htmlFor="r1">No Notifications</Label>
        </div>
        <div className="flex items-center space-x-2">
          <RadioGroupItem value="planning-basic" id="r2" />
          <Label htmlFor="r2">Basic</Label>
        </div>
        <div className="flex items-center space-x-2">
          <RadioGroupItem value="planning-detailed" id="r3" />
          <Label htmlFor="r3">Detailed</Label>
        </div>
      </RadioGroup>

      <div className="pb-4" />

      <h2 className="text-lg font-bold">Transfer</h2>
      <div className="pb-1" />
      <RadioGroup defaultValue="transfer-detailed">
        <div className="flex items-center space-x-2">
          <RadioGroupItem
            value="transfer-no-notifications"
            id="r1"
            onClick={() => console.log('transferonot')}
          />
          <Label htmlFor="r1">No Notifications</Label>
        </div>
        <div className="flex items-center space-x-2">
          <RadioGroupItem value="transfer-basic" id="r2" />
          <Label htmlFor="r2">Basic</Label>
        </div>
        <div className="flex items-center space-x-2">
          <RadioGroupItem value="transfer-detailed" id="r3" />
          <Label htmlFor="r3">Detailed</Label>
        </div>
      </RadioGroup>
    </div>
  );
};
