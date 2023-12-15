import { Node } from '../types';

export function findNode(tree: Node[], keyToFind: string): Node | null {
  const queue: Node[] = [...tree];

  while (queue.length > 0) {
    const currentNode = queue.shift()!; // Dequeue the front node

    // Check if the current node's key matches the key to find
    console.log('ckey - ktf  ', currentNode.key, keyToFind);
    if (currentNode.key === keyToFind) {
      return currentNode;
    }

    // Enqueue the children of the current node (if any)
    if (currentNode.children && currentNode.children.length > 0) {
      queue.push(...currentNode.children);
    }
  }

  // If the key is not found, return null
  return null;
}
