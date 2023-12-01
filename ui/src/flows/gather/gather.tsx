import React from 'react';

import { Navbar } from './navbar';
import { Pane } from '~/shared/pane/pane';
import { Ticker } from './ticker';
import { Content } from './content';

export const Gather: React.FunctionComponent = () => {
  return (
    <div className="flex flex-col h-full">
      <Navbar />
      <Pane>
        <Ticker />
      </Pane>
      <Content />
    </div>
  );
};
