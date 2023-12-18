import React from 'react';

import { TreeNode } from './node';

export interface ITreeNode {
  id: string;
  label: string;
  leaf: boolean;
  parent: string;
  children: string[];
  checked?: boolean;
  expanded?: boolean;
  loading?: boolean;
}

export type Nodes = Record<string, ITreeNode>;

export interface TreeProps {
  nodes: Nodes;
  onLoad: (node: ITreeNode) => void;
  onCheck: (node: ITreeNode) => void;
  collapseIcon: React.ReactElement;
  expandIcon: React.ReactElement;
  checkedIcon: React.ReactElement;
  uncheckedIcon: React.ReactElement;
  parentIcon: React.ReactElement;
  leafIcon: React.ReactElement;
  placeholderIcon: React.ReactElement;
  loadingIcon: React.ReactElement;
}

export const Tree: React.FunctionComponent<TreeProps> = ({
  nodes = {},
  onLoad,
  onCheck,
  collapseIcon,
  expandIcon,
  checkedIcon,
  uncheckedIcon,
  parentIcon,
  leafIcon,
  placeholderIcon,
  loadingIcon,
}) => {
  const getRootNodes = (list: Nodes) => {
    const items = Object.values(list).filter((node) => node.parent === '');
    console.log('getRootNodes ', items);
    return items;
  };

  const getChildNodes = (node: ITreeNode) =>
    node.children ? node.children.map((id) => nodes[id]) : [];

  const onExpandCollapse = (node: ITreeNode) => {
    onLoad(node);
  };

  const onCheckUncheck = (node: ITreeNode) => {
    onCheck(node);
  };

  return (
    <div>
      {getRootNodes(nodes).map((node) => (
        <TreeNode
          node={node}
          getChildNodes={getChildNodes}
          level={0}
          onExpandCollapse={onExpandCollapse}
          onCheckUncheck={onCheckUncheck}
          collapseIcon={collapseIcon}
          expandIcon={expandIcon}
          checkedIcon={checkedIcon}
          uncheckedIcon={uncheckedIcon}
          parentIcon={parentIcon}
          leafIcon={leafIcon}
          placeholderIcon={placeholderIcon}
          loadingIcon={loadingIcon}
        />
      ))}
    </div>
  );
};
