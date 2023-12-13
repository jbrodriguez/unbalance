import React from 'react';
import CheckboxTree from 'react-checkbox-tree';
// import IconAccountBox from 'virtual:vite-icons/ic/baseline-check-box/';
// import { Icon } from '~/shared/icons/icon';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import {
  faCheckSquare,
  faSquareMinus,
  faChevronRight,
  faChevronDown,
  faPlusSquare,
  faMinusSquare,
  faFolder,
  faFolderOpen,
  faFile,
} from '@fortawesome/free-solid-svg-icons';
import { faSquare } from '@fortawesome/free-regular-svg-icons';

import 'react-checkbox-tree/lib/react-checkbox-tree.css';

const nodes = [
  {
    value: 'Documents',
    label: 'Documents',
    children: [
      {
        value: 'Employee Evaluations.zip',
        label: 'Employee Evaluations.zip',
        // icon: <i className="far fa-file-archive" />,
      },
      {
        value: 'Expense Report.pdf',
        label: 'Expense Report.pdf',
        // icon: <i className="far fa-file-pdf" />,
      },
      {
        value: 'notes.txt',
        label: 'notes.txt',
        // icon: <i className="far fa-file-alt" />,
      },
    ],
  },
  {
    value: 'Photos',
    label: 'Photos',
    children: [
      {
        value: 'nyan-cat.gif',
        label: 'nyan-cat.gif',
        // icon: <i className="far fa-file-image" />,
      },
      {
        value: 'SpaceX Falcon9 liftoff.jpg',
        label: 'SpaceX Falcon9 liftoff.jpg',
        // icon: <i className="far fa-file-image" />,
      },
    ],
  },
];

interface Props {
  height?: number;
}

export const FileSystem: React.FC<Props> = ({ height = 0 }) => {
  const [checked, setChecked] = React.useState<string[]>([]);
  const [expanded, setExpanded] = React.useState<string[]>(['Documents']);

  const onCheck = (value: string[]) => {
    console.log('onCheck ', value);
    setChecked(value);
  };

  const onExpand = (value: string[]) => {
    console.log('onExpand ', value);
    setExpanded(value);
  };

  return (
    <div className="flex flex-1 bg-neutral-200 dark:bg-gray-950">
      <div className={`overflow-y-auto`} style={{ height: `${height}px` }}>
        <FontAwesomeIcon
          className="rct-icon rct-icon-check"
          icon="check-square"
        />

        <CheckboxTree
          checked={checked}
          expanded={expanded}
          nodes={nodes}
          onCheck={onCheck}
          onExpand={onExpand}
          icons={{
            check: (
              <FontAwesomeIcon
                className="fill-red-400 text-yellow-600"
                icon={faCheckSquare}
              />
            ),
            uncheck: (
              <FontAwesomeIcon className="text-slate-700" icon={faSquare} />
            ),
            halfCheck: (
              <FontAwesomeIcon
                className="rct-icon rct-icon-half-check"
                icon={faSquareMinus}
              />
            ),
            expandClose: (
              <FontAwesomeIcon
                className="rct-icon rct-icon-expand-close"
                icon={faChevronRight}
              />
            ),
            expandOpen: (
              <FontAwesomeIcon
                className="rct-icon rct-icon-expand-open"
                icon={faChevronDown}
              />
            ),
            expandAll: (
              <FontAwesomeIcon className="fill-red-600" icon={faPlusSquare} />
            ),
            collapseAll: (
              <FontAwesomeIcon
                className="rct-icon rct-icon-collapse-all"
                icon={faMinusSquare}
              />
            ),
            parentClose: (
              <FontAwesomeIcon
                className="rct-icon rct-icon-parent-close"
                icon={faFolder}
              />
            ),
            parentOpen: (
              <FontAwesomeIcon
                className="rct-icon rct-icon-parent-open"
                icon={faFolderOpen}
              />
            ),
            leaf: (
              <FontAwesomeIcon
                className="rct-icon rct-icon-leaf-close"
                icon={faFile}
              />
            ),
          }}
          // icons={{
          //   check: (
          //     <Icon name="checkbox" fill="dark:fill-slate-400 fill-slate-400" />
          //   ),
          //   uncheck: <Icon name="checkbox-blank" />,
          //   halfCheck: <span className="rct-icon rct-icon-half-check" />,
          //   expandClose: <span className="rct-icon rct-icon-expand-close" />,
          //   expandOpen: <span className="rct-icon rct-icon-expand-open" />,
          //   expandAll: <span className="rct-icon rct-icon-expand-all" />,
          //   collapseAll: <span className="rct-icon rct-icon-collapse-all" />,
          //   parentClose: <span className="rct-icon rct-icon-parent-close" />,
          //   parentOpen: <span className="rct-icon rct-icon-parent-open" />,
          //   leaf: <span className="rct-icon rct-icon-leaf" />,
          // }}
        />
      </div>
    </div>
  );
};
