import { Node, Nodes } from '../types';

export const isParent = (id: string, nodes: Nodes) =>
  Object.values(nodes).some((n) => n.parent === id);

export const getAbsolutePath = (node: Node, nodes: Nodes): string => {
  const parent = nodes[node.parent];
  if (!parent || parent.id === 'root') {
    return node.label;
  }
  return `${getAbsolutePath(parent, nodes)}/${node.label}`;
};
