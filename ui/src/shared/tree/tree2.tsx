import React from 'react';

import { TreeNode } from './node';

export interface ITreeNode {
  id: string;
  key: string;
  label: string;
  leaf: boolean;
  parent: string;
  children: string[];
  checked?: boolean;
  expanded?: boolean;
  loading?: boolean;
  icon?: React.ReactElement;
  checkbox?: React.ReactElement;
}

export type Nodes = Record<string, ITreeNode>;

export interface TreeProps {
  initialNodes?: Nodes;
  onLoad: (nodes: Nodes, id: string) => Promise<Nodes>;
  onCheck?: (node: ITreeNode) => void;
  // onExpand?: (node: TreeNode) => void;
  // spinner: React.ReactElement;
  // collapseIcon: React.ReactElement;
  // expandIcon: React.ReactElement;
}

export const Tree: React.FunctionComponent<TreeProps> = ({
  initialNodes = {},
  onLoad,
  // onCheck,
}) => {
  const [nodes, setNodes] = React.useState<Nodes>(initialNodes);

  // React.useEffect(() => {
  //   console.log('nodes', data);
  //   setNodes(data);
  // }, [data]);

  const getRootNodes = () => {
    const items = Object.values(nodes).filter((node) => node.parent === '');
    console.log('getRootNodes ', items);
    return items;
  };

  const getChildNodes = (node: ITreeNode) =>
    node.children ? node.children.map((id) => nodes[id]) : [];

  const isParent = (id: string) =>
    Object.values(nodes).some((n) => n.parent === id);

  const onToggle = async (node: ITreeNode) => {
    // nodes[node.id].expanded = !nodes[node.id].expanded;
    // nodes[node.id].loading = true;
    // setNodes({ ...nodes });
    // const loaded = await onLoad(nodes, node.id);
    // console.log('loaded ', loaded);
    // nodes[node.id].loading = false;
    // setNodes(loaded);
    nodes[node.id].expanded = !nodes[node.id].expanded;

    if (isParent(node.id)) {
      setNodes({ ...nodes });
      return;
    }

    const items = { ...nodes };
    items.loader = {
      id: 'loader',
      key: 'loader',
      label: 'loading ...',
      leaf: false,
      parent: node.id,
      children: [],
      checked: false,
      expanded: false,
      loading: false,
    };
    items[node.id].children = ['loader'];
    setNodes({ ...items });

    const loaded = await onLoad(nodes, node.id);
    delete items.loader;
    items[node.id].children = Object.keys(loaded);
    setNodes({ ...items, ...loaded });
  };

  return (
    <div>
      {getRootNodes().map((node) => (
        <TreeNode
          node={node}
          getChildNodes={getChildNodes}
          level={0}
          onToggle={onToggle}
        />
      ))}
    </div>
  );
};
