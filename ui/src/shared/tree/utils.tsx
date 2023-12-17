import { ITreeNode } from './tree2';

export function decorateNode(node: ITreeNode) {
  node.checked = false;
  node.expanded = false;
  node.loading = false;
}
