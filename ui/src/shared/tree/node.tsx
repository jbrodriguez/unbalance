import React from 'react';

import { ITreeNode } from './tree2';

interface TreeNodeProps {
  node: ITreeNode;
  getChildNodes: (node: ITreeNode) => ITreeNode[];
  onToggle: (node: ITreeNode) => void;
  level: number;
}

export const TreeNode: React.FunctionComponent<TreeNodeProps> = ({
  node,
  getChildNodes,
  onToggle,
  level,
}) => {
  // const [expanded, setExpanded] = React.useState<boolean>(false);

  const handleToggle = () => {
    // setExpanded(!expanded);
  };

  console.log('rendering node ', node);

  return (
    <>
      <div style={{ paddingLeft: `${level * 1}rem` }}>
        <span onClick={() => onToggle(node)}>{node.expanded ? '-' : '+'}</span>

        <span className="ml-2" onClick={handleToggle}>
          checkbox
        </span>

        {node.loading ? (
          <span className="ml-2" onClick={handleToggle}>
            loading
          </span>
        ) : (
          <span className="ml-2" onClick={handleToggle}>
            icon
          </span>
        )}

        <span className="ml-2">{node.label}</span>
      </div>
      {node.expanded &&
        getChildNodes(node).map((childNode) => (
          <TreeNode
            node={childNode}
            getChildNodes={getChildNodes}
            onToggle={onToggle}
            level={level + 1}
          />
        ))}
    </>
  );
};
