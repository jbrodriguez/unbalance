import React from 'react';

export const Description: React.FunctionComponent = () => {
  return (
    <>
      <span>
        Validate planned transfers, then choose to either{' '}
        <span className="font-bold">MOVE</span> or{' '}
        <span className="font-bold">COPY</span> the data using the buttons
        above, you can also choose to <span className="font-bold">dry-run</span>{' '}
        first (or not)
        <br />
        <span className="font-bold">*NOTE</span>: planned sizes are based on a
        move scenario (copy scenario will be the same, except the source disk
        will not change in free size)
      </span>
    </>
  );
};
