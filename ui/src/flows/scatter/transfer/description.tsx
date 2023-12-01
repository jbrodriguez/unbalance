import React from 'react';

export const Description: React.FunctionComponent = () => {
  return (
    <span>
      Choose to either <span className="font-bold">MOVE</span> or{' '}
      <span className="font-bold">COPY</span> the data using the buttons above,
      you can also choose to <span className="font-bold">dry-run</span> first
      (or not)
    </span>
  );
};
