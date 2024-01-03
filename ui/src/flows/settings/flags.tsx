import React from 'react';

import { Input } from '@/components/ui/input';
import { Button } from '@/components/ui/button';
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuLabel,
  DropdownMenuRadioGroup,
  DropdownMenuRadioItem,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from '@/components/ui/dropdown-menu';

export const Flags: React.FunctionComponent = () => {
  const [position, setPosition] = React.useState('bottom');

  return (
    <div className="text-slate-700 dark:text-gray-300 p-4">
      <h1>
        unbalanced uses the threshold defined here as the minimum free space
        that should be kept available in a target disk, when planning how much
        the disk can be filled. <br />
        This threshold cannot be less than 512Mb (hard limit set by this app).
      </h1>
      <div className="pb-4" />

      <div className="flex flex-row items-center">
        <Input type="number" placeholder="Size" className="w-40" />
        <div className="pr-4" />
        <DropdownMenu>
          <DropdownMenuTrigger asChild>
            <Button>{position}</Button>
          </DropdownMenuTrigger>
          <DropdownMenuContent className="">
            <DropdownMenuLabel>Unit</DropdownMenuLabel>
            <DropdownMenuSeparator />
            <DropdownMenuRadioGroup
              value={position}
              onValueChange={setPosition}
            >
              <DropdownMenuRadioItem value="top">Top</DropdownMenuRadioItem>
              <DropdownMenuRadioItem value="bottom">
                Bottom
              </DropdownMenuRadioItem>
              <DropdownMenuRadioItem value="right">Right</DropdownMenuRadioItem>
            </DropdownMenuRadioGroup>
          </DropdownMenuContent>
        </DropdownMenu>
      </div>
    </div>
  );
};
