import { create } from 'zustand';
import { immer } from 'zustand/middleware/immer';

import { Nodes, Node } from '~/types';
import { Api } from '~/api';
import { decorateNode } from '~/shared/tree/utils';
import { isParent, getAbsolutePath } from '~/helpers/tree';

interface GatherStore {
  source: string;
  tree: Nodes;
  selected: Record<string, string>;
  location: Record<string, Array<string>>;
  target: string;
  actions: {
    loadShares: () => Promise<void>;
    loadBranch: (node: Node) => Promise<void>;
    toggleSelected: (node: Node) => Promise<void>;
    setTarget: (target: string) => void;
    // loadBranch: (node: Node) => Promise<void>;
    // toggleSelected: (node: Node) => void;
    // toggleTarget: (name: string) => void;
  };
}

const rootNode = {
  id: 'root',
  label: '/',
  leaf: false,
  dir: false,
  parent: '',
};

// const loaderNode = {
//   id: 'loader',
//   label: 'loading ...',
//   leaf: false,
//   parent: 'root',
// };

export const useGatherStore = create<GatherStore>()(
  immer((set, get) => ({
    source: 'user',
    // selected: [],
    // targets: {},
    tree: { root: decorateNode(rootNode as Node) },
    selected: {},
    location: {},
    target: '',

    actions: {
      loadShares: async () => {
        const actions = get().actions;
        actions.loadBranch(decorateNode(rootNode as Node));
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
      toggleSelected: async (node: Node) => {
        set((state) => {
          state.tree[node.id].checked = !state.tree[node.id].checked;
        });

        set((state) => {
          // remove parents by looping
          let parent = state.tree[node.parent];
          while (parent) {
            if (state.selected[parent.id]) {
              delete state.selected[parent.id];
              delete state.location[parent.id];
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
              if (state.selected[child.id]) {
                delete state.selected[child.id];
                delete state.location[child.id];
                state.tree[child.id].checked = false;
              }
              removeChildren(child);
            });
          };
          removeChildren(node);
        });

        // get().location[fullPath] = location;

        // const fullPath = getAbsolutePath(node, get().tree);
        // const branch = await Api.locate(fullPath, node.id);
        // for (const key in branch.nodes) {
        //   decorateNode(branch.nodes[key]);
        // }

        // const selected = get().selected;
        if (get().selected[node.id]) {
          set((state) => {
            delete state.selected[node.id];
            delete state.location[node.id];
            state.selected = { ...state.selected };
            state.location = { ...state.location };
          });
          return;
        }

        const fullPath = getAbsolutePath(node, get().tree);
        console.log('fullPath ', fullPath);

        const location = await Api.locate(fullPath);
        console.log('location ', location);

        set((state) => {
          state.selected[node.id] = fullPath;
          state.location[node.id] = location;
        });
      },
      setTarget: (target: string) => {
        set((state) => {
          state.target = target;
        });
      },
    },
  })),
);

export const useGatherActions = () => useGatherStore((state) => state.actions);

export const useGatherTree = () => useGatherStore((state) => state.tree);
export const useGatherSelected = () =>
  useGatherStore((state) => state.selected);
export const useGatherLocation = () =>
  useGatherStore((state) => state.location);
export const useGatherTarget = () => useGatherStore((state) => state.target);
