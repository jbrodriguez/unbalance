import { create } from 'zustand';
import { immer } from 'zustand/middleware/immer';

import { Targets, Nodes, Node } from '~/types';
import { Api } from '~/api';
import { decorateNode } from '~/shared/tree/utils';
import { isParent, getAbsolutePath } from '~/helpers/tree';

interface ScatterStore {
  source: string;
  selected: Array<string>;
  targets: Targets;
  tree: Nodes;
  binDisk: string;
  actions: {
    setSource: (source: string) => Promise<void>;
    loadBranch: (node: Node) => Promise<void>;
    toggleSelected: (node: Node) => void;
    toggleTarget: (name: string) => void;
    setBinDisk: (binDisk: string) => void;
  };
}

const rootNode = {
  id: 'root',
  label: '/',
  leaf: false,
  dir: false,
  parent: '',
};

const loaderNode = {
  id: 'loader',
  label: 'loading ...',
  leaf: false,
  dir: false,
  parent: 'root',
};

export const useScatterStore = create<ScatterStore>()(
  immer((set, get) => ({
    source: '',
    selected: [],
    targets: {},
    tree: { root: decorateNode(rootNode as Node) },
    logs: [],
    binDisk: '',
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
          // change reference to force re-render and show expanded/non-expanded state
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
            dir: false,
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

          // add or remove from selected
          const fullPath = getAbsolutePath(node, state.tree);
          const index = state.selected.indexOf(fullPath);
          if (index === -1) {
            state.selected.push(fullPath);
          } else {
            state.selected.splice(index, 1);
          }

          // remove parents by looping
          let parent = state.tree[node.parent];
          while (parent) {
            const parentFullPath = getAbsolutePath(parent, state.tree);
            const parentIndex = state.selected.indexOf(parentFullPath);
            if (parentIndex !== -1) {
              state.selected.splice(parentIndex, 1);
              state.tree[parent.id].checked = false;
            }
            parent = state.tree[parent.parent];
          }

          // remove children recursively
          const removeChildren = (node: Node) => {
            if (!node.children) {
              return;
            }

            node.children.forEach((childId) => {
              const child = state.tree[childId];
              const childFullPath = getAbsolutePath(child, state.tree);
              const childIndex = state.selected.indexOf(childFullPath);
              if (childIndex !== -1) {
                state.selected.splice(childIndex, 1);
                state.tree[child.id].checked = false;
              }
              removeChildren(child);
            });
          };
          removeChildren(node);
        });
      },
      toggleTarget: (name: string) => {
        set((state) => {
          if (state.targets[name] === undefined) {
            state.targets[name] = true;
            return;
          }

          delete state.targets[name];
        });
      },
      setBinDisk: (binDisk: string) => {
        set((state) => {
          state.binDisk = binDisk;
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
export const useScatterBinDisk = () =>
  useScatterStore((state) => state.binDisk);
