import { Node } from '~/types';

export const decorateNode = (node: Node): Node => {
  node.checked = false;
  node.expanded = false;
  node.loading = false;
  node.children = [];

  return node;
};
