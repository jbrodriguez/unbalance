import React from 'react';

import { ITreeNode } from './tree2';

interface TreeNodeProps {
  node: ITreeNode;
  getChildNodes: (node: ITreeNode) => ITreeNode[];
  onExpandCollapse: (node: ITreeNode) => void;
  onCheckUncheck: (node: ITreeNode) => void;
  level: number;
  collapseIcon: React.ReactElement;
  expandIcon: React.ReactElement;
  checkedIcon: React.ReactElement;
  uncheckedIcon: React.ReactElement;
  parentIcon: React.ReactElement;
  leafIcon: React.ReactElement;
  placeholderIcon: React.ReactElement;
  loadingIcon: React.ReactElement;
}

export const TreeNode: React.FunctionComponent<TreeNodeProps> = ({
  node,
  getChildNodes,
  onExpandCollapse,
  onCheckUncheck,
  level,
  collapseIcon,
  expandIcon,
  checkedIcon,
  uncheckedIcon,
  parentIcon,
  leafIcon,
  placeholderIcon,
  loadingIcon,
}) => {
  // const [expanded, setExpanded] = React.useState<boolean>(false);

  const handleToggle = () => {
    // setExpanded(!expanded);
  };

  // console.log('rendering node ', node);

  const renderNode = (node: ITreeNode) => {
    if (node.loading) {
      return (
        <>
          <span className="ml-1">{loadingIcon}</span>
          <span className="ml-1">{node.label}</span>
        </>
      );
    }

    return (
      <>
        {!node.leaf ? (
          <span onClick={() => onExpandCollapse(node)}>
            {node.expanded ? collapseIcon : expandIcon}
          </span>
        ) : (
          <span>{placeholderIcon}</span>
        )}

        <span className="ml-1" onClick={() => onCheckUncheck(node)}>
          {node.checked ? checkedIcon : uncheckedIcon}
        </span>

        <span className="ml-1" onClick={handleToggle}>
          {node.leaf ? leafIcon : parentIcon}
        </span>

        <span className="ml-1">{node.label}</span>
      </>
    );
  };

  return (
    <>
      <div
        style={{ paddingLeft: `${level * 24}px` }}
        className="flex flex-row items-center"
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
    </>
  );
};
