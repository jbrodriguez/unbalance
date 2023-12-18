import { create } from 'zustand';
import { immer } from 'zustand/middleware/immer';

import { Targets, Nodes, Node } from '~/types';
import { Api } from '~/api';
import { decorateNode } from '~/shared/tree/utils';

interface ScatterStore {
  source: string;
  selected: Array<string>;
  targets: Targets;
  tree: Nodes;
  actions: {
    setSource: (source: string) => Promise<void>;
    loadBranch: (node: Node) => Promise<void>;
    toggleSelected: (node: Node) => void;
    toggleTarget: (name: string) => void;
  };
}

const rootNode = {
  id: 'root',
  label: '/',
  leaf: false,
  parent: '',
};

const loaderNode = {
  id: 'loader',
  label: 'loading ...',
  leaf: false,
  parent: 'root',
};

const isParent = (id: string, nodes: Nodes) =>
  Object.values(nodes).some((n) => n.parent === id);

const getAbsolutePath = (node: Node, nodes: Nodes): string => {
  const parent = nodes[node.parent];
  if (!parent) {
    return node.label;
  }
  return `${getAbsolutePath(parent, nodes)}/${node.label}`;
};

export const useScatterStore = create<ScatterStore>()(
  immer((set, get) => ({
    source: '',
    selected: [],
    targets: {},
    tree: { root: decorateNode(rootNode as Node) },

    actions: {
      setSource: async (source: string) => {
        const loader = decorateNode({ ...loaderNode } as Node);
        loader.loading = true;

        set((state) => {
          state.source = source;
          state.targets = {};
          state.selected = [];
          state.tree.root.children = ['loader'];
          state.tree = { ...state.tree, loader };
        });

        const route = `${get().source}/`;
        console.log('route ', route);
        const branch = await Api.getTree(route, 'root');
        // await new Promise((r) => setTimeout(r, 5000));
        for (const key in branch.nodes) {
          decorateNode(branch.nodes[key]);
        }

        console.log('decorated ', branch.nodes);

        set((state) => {
          delete state.tree.loader;
          state.tree = { ...state.tree, ...branch.nodes };
          state.tree.root.children = branch.order;
        });
      },
      loadBranch: async (node: Node) => {
        set((state) => {
          state.tree[node.id].expanded = !state.tree[node.id].expanded;
        });

        if (isParent(node.id, get().tree)) {
          set((state) => {
            state.tree = { ...state.tree };
          });
          return;
        }

        set((state) => {
          state.tree.loader = {
            id: 'loader',
            label: 'loading ...',
            leaf: false,
            parent: node.id,
            children: [],
            checked: false,
            expanded: false,
            loading: true,
          };
          state.tree[node.id].children = ['loader'];
        });

        const route = `${get().source}/${getAbsolutePath(node, get().tree)}`;
        console.log('route ', route);
        const branch = await Api.getTree(route, node.id);
        // await new Promise((r) => setTimeout(r, 5000));
        for (const key in branch.nodes) {
          decorateNode(branch.nodes[key]);
        }

        set((state) => {
          delete state.tree.loader;
          state.tree = { ...state.tree, ...branch.nodes };
          state.tree[node.id].children = branch.order;
        });
      },
      toggleSelected: (node: Node) => {
        set((state) => {
          console.log('toggleSelected ', node);
          state.tree[node.id].checked = !state.tree[node.id].checked;
          console.log('node.id ', state.tree[node.id]);
          const fullPath = getAbsolutePath(node, state.tree);
          const index = state.selected.indexOf(fullPath);
          if (index === -1) {
            state.selected.push(fullPath);
          } else {
            state.selected.splice(index, 1);
          }
        });
      },
      toggleTarget: (name: string) => {
        set((state) => {
          state.targets[name] = !state.targets[name];
        });
      },
    },
  })),
);

export const useScatterActions = () =>
  useScatterStore((state) => state.actions);

export const useScatterSource = () => useScatterStore((state) => state.source);
// export const useScatterRoots = () =>
//   useScatterStore((state) => state.roots.map((root) => state.tree[root]));
export const useScatterTree = () => useScatterStore((state) => state.tree);
export const useScatterSelected = () =>
  useScatterStore((state) => state.selected);
export const useScatterTargets = () =>
  useScatterStore((state) => state.targets);
