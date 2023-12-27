import React from 'react';

export const Description: React.FunctionComponent = () => {
  return (
    <span>
      Select <span className="font-bold">ONE</span> target disk where the
      shares/folders on the left will be gathered, then click on{' '}
      <span className="font-bold">MOVE</span>
      <br />
      <span className="font-bold">*NOTE</span>: Drives are ordered by the least
      amount of data transfer. A star next to a disk means that one or more
      source folders are present. Some folders may not be picked up due to free
      space constraints.
    </span>
  );
};
