import { Node } from '~/types';

export const decorateNode = (node: Node): Node => ({
  ...node,
  checked: false,
  expanded: false,
  loading: false,
  children: [],
});
