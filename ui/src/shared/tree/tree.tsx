import React, { useState, useEffect } from 'react';

export interface TreeNode {
  id: string;
  label: string;
  value: string;
  icon: React.ReactElement;
  checkbox: React.ReactElement;
  children?: TreeNode[];
  isChecked: boolean;
  isExpanded: boolean;
  isLoading?: boolean;
}

export interface TreeProps {
  data: TreeNode[];
  onCheck: (node: TreeNode) => void;
  onLoad: (node: TreeNode) => Promise<void>;
  onExpand?: (node: TreeNode) => void;
  onCollapse?: (node: TreeNode) => void;
  spinner: React.ReactElement;
  collapseIcon: React.ReactElement;
  expandIcon: React.ReactElement;
}

const TreeNodeComponent: React.FC<TreeProps> = ({
  data,
  onCheck,
  onLoad,
  onExpand,
  onCollapse,
  spinner,
  collapseIcon,
  expandIcon,
}) => {
  const [expandedNodes, setExpandedNodes] = useState<string[]>([]);

  const toggleNode = async (node: TreeNode) => {
    if (node.isExpanded) {
      setExpandedNodes((prevExpanded) =>
        prevExpanded.filter((id) => id !== node.id),
      );
      onCollapse?.(node);
    } else {
      onExpand?.(node);
      setExpandedNodes((prevExpanded) => [...prevExpanded, node.id]);

      if (node.children && !node.isLoading) {
        // Only load if there are children and not already loading
        try {
          node.isLoading = true;
          await onLoad(node);
        } finally {
          node.isLoading = false;
        }
      }
    }
  };

  const renderTreeNodes = (nodes: TreeNode[], depth: number) => {
    return nodes.map((node) => (
      <div key={node.id} style={{ marginLeft: `${depth * 20}px` }}>
        <button onClick={() => toggleNode(node)}>
          {expandedNodes.includes(node.id) ? collapseIcon : expandIcon}
        </button>
        <span onClick={() => onCheck(node)}>{node.checkbox}</span>
        {node.isLoading ? spinner : node.icon}
        {node.label}
        {expandedNodes.includes(node.id) &&
          node.children &&
          renderTreeNodes(node.children, depth + 1)}
      </div>
    ));
  };

  useEffect(() => {
    setExpandedNodes([]);
  }, [data]); // Reset expanded nodes when data changes

  return <div>{renderTreeNodes(data, 0)}</div>;
};

export const Tree: React.FC<TreeProps> = ({
  data,
  onCheck,
  onLoad,
  onExpand,
  onCollapse,
  spinner,
  collapseIcon,
  expandIcon,
}) => {
  const [cachedNodes, setCachedNodes] = useState<TreeNode[]>(data);

  const handleLoad = async (node: TreeNode) => {
    await onLoad(node);
    setCachedNodes((prevNodes) => {
      const updatedNodes = [...prevNodes];
      const index = updatedNodes.findIndex((n) => n.id === node.id);
      if (index !== -1) {
        updatedNodes[index] = {
          ...updatedNodes[index],
          children: node.children,
          isLoading: false,
        };
      }
      return updatedNodes;
    });
  };

  return (
    <TreeNodeComponent
      data={cachedNodes}
      onCheck={onCheck}
      onLoad={handleLoad}
      onExpand={onExpand}
      onCollapse={onCollapse}
      spinner={spinner}
      collapseIcon={collapseIcon}
      expandIcon={expandIcon}
    />
  );
};

// please improve the onload logic so that isloading is set to true when the node is expanded and set to false when the node is collapsed
