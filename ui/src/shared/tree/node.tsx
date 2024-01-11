import React from 'react';

import { Node, Icons } from '~/types';

interface TreeNodeProps {
  node: Node;
  getChildNodes: (node: Node) => Node[];
  onExpandCollapse: (node: Node) => void;
  onCheckUncheck: (node: Node) => void;
  icons: Icons;
  level: number;
}

export const TreeNode: React.FunctionComponent<TreeNodeProps> = ({
  node,
  getChildNodes,
  onExpandCollapse,
  onCheckUncheck,
  icons,
  level,
}) => {
  const renderNode = (node: Node) => {
    if (node.loading) {
      return (
        <>
          <span className="ml-1">{icons.loadingIcon}</span>
          <span className="ml-1">{node.label}</span>
        </>
      );
    }

    return (
      <>
        {!node.leaf ? (
          <span onClick={() => onExpandCollapse(node)}>
            {node.expanded ? icons.collapseIcon : icons.expandIcon}
          </span>
        ) : (
          <span>{icons.hiddenIcon}</span>
        )}

        <span className="ml-1" onClick={() => onCheckUncheck(node)}>
          {node.checked ? icons.checkedIcon : icons.uncheckedIcon}
        </span>

        <span className="ml-1">
          {node.leaf ? icons.leafIcon : icons.parentIcon}
        </span>

        <span className="ml-1">{node.label}</span>
      </>
    );
  };

  return (
    <>
      <div
        style={{ paddingLeft: `${level * 24}px` }}
        className="flex flex-1 flex-row items-center whitespace-nowrap"
      >
        {renderNode(node)}
      </div>
      {node.expanded &&
        getChildNodes(node).map((childNode) => (
          <TreeNode
            node={childNode}
            getChildNodes={getChildNodes}
            onExpandCollapse={onExpandCollapse}
            onCheckUncheck={onCheckUncheck}
            level={level + 1}
            icons={icons}
          />
        ))}
    </>
  );
};
