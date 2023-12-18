import React from 'react';

import { Node, Nodes, Icons } from '~/types';
import { TreeNode } from './node';

export interface CheckboxTreeProps {
  nodes: Nodes;
  onLoad: (node: Node) => void;
  onCheck: (node: Node) => void;
  icons: Icons;
}

export const CheckboxTree: React.FunctionComponent<CheckboxTreeProps> = ({
  nodes = {},
  onLoad,
  onCheck,
  icons,
}) => {
  const getRootNodes = (list: Nodes) =>
    Object.values(list).filter((node) => node.parent === '');

  const getChildNodes = (node: Node) =>
    node.children ? node.children.map((id) => nodes[id]) : [];

  const onExpandCollapse = (node: Node) => onLoad(node);

  const onCheckUncheck = (node: Node) => onCheck(node);

  return (
    <>
      {getRootNodes(nodes).map((node) => (
        <TreeNode
          node={node}
          getChildNodes={getChildNodes}
          onExpandCollapse={onExpandCollapse}
          onCheckUncheck={onCheckUncheck}
          icons={icons}
          level={0}
        />
      ))}
    </>
  );
};
