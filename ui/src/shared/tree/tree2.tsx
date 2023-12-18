import React from 'react';

import { TreeNode } from './node';
// import { decorateNode } from './utils';

export interface ITreeNode {
  id: string;
  // key: string;
  label: string;
  leaf: boolean;
  parent: string;
  children: string[];
  checked?: boolean;
  expanded?: boolean;
  loading?: boolean;
  // icon?: React.ReactElement;
  // checkbox?: React.ReactElement;
}

export type Nodes = Record<string, ITreeNode>;

export interface TreeProps {
  nodes: Nodes;
  onLoad: (node: ITreeNode) => void;
  onCheck: (node: ITreeNode) => void;
  // onExpand?: (node: TreeNode) => void;
  // spinner: React.ReactElement;
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
  // const [nodes, setNodes] = React.useState<Nodes>(initialNodes);

  // React.useEffect(() => {
  //   console.log('nodes', data);
  //   setNodes(data);
  // }, [data]);

  const getRootNodes = (list: Nodes) => {
    const items = Object.values(list).filter((node) => node.parent === '');
    console.log('getRootNodes ', items);
    return items;
  };

  const getChildNodes = (node: ITreeNode) =>
    node.children ? node.children.map((id) => nodes[id]) : [];

  // const getAbsolutePath = (node: ITreeNode): string => {
  //   const parent = nodes[node.parent];
  //   if (!parent) {
  //     return node.label;
  //   }
  //   return `${getAbsolutePath(parent)}/${node.label}`;
  // };

  // const isParent = (id: string) =>
  //   Object.values(nodes).some((n) => n.parent === id);

  const onExpandCollapse = (node: ITreeNode) => {
    // // nodes[node.id].expanded = !nodes[node.id].expanded;
    // // nodes[node.id].loading = true;
    // // setNodes({ ...nodes });
    // // const loaded = await onLoad(nodes, node.id);
    // // console.log('loaded ', loaded);
    // // nodes[node.id].loading = false;
    // // setNodes(loaded);
    // nodes[node.id].expanded = !nodes[node.id].expanded;

    // // if (isParent(node.id)) {
    // //   // setNodes({ ...nodes });
    // //   return;
    // // }

    // const draft = { ...nodes };
    // draft.loader = {
    //   id: 'loader',
    //   // key: 'loader',
    //   label: 'loading ...',
    //   leaf: false,
    //   parent: node.id,
    //   children: [],
    //   checked: false,
    //   expanded: false,
    //   loading: false,
    // };
    // draft[node.id].children = ['loader'];
    // // setNodes({ ...draft });

    // const loaded = await onLoad(getAbsolutePath(node), node.id);
    // for (const key in loaded) {
    //   decorateNode(loaded[key]);
    // }
    // delete draft.loader;
    // draft[node.id].children = Object.keys(loaded);
    // // setNodes({ ...draft, ...loaded });
    onLoad(node);
  };

  const onCheckUncheck = (node: ITreeNode) => {
    // nodes[node.id].checked = !nodes[node.id].checked;
    // setNodes({ ...nodes });
    // if (onCheck) {
    //   onCheck(node);
    // }
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
