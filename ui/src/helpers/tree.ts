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

// returns the nodes between anchorId and nodeId (inclusive, in display order),
// provided both are siblings under the same parent; null otherwise
export const getSiblingRange = (
  anchorId: string,
  nodeId: string,
  nodes: Nodes,
): Node[] | null => {
  const anchor = nodes[anchorId];
  const node = nodes[nodeId];
  if (!anchor || !node || anchor.parent !== node.parent) {
    return null;
  }

  const siblings = nodes[node.parent]?.children ?? [];
  const from = siblings.indexOf(anchorId);
  const to = siblings.indexOf(nodeId);
  if (from === -1 || to === -1) {
    return null;
  }

  const [start, end] = from < to ? [from, to] : [to, from];
  return siblings.slice(start, end + 1).map((id) => nodes[id]);
};
