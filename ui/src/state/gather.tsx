import { create } from 'zustand';
import { immer } from 'zustand/middleware/immer';

import { Nodes, Node, Sizes } from '~/types';
import { Api } from '~/api';
import { decorateNode } from '~/shared/tree/utils';
import { isParent, getAbsolutePath, getSiblingRange } from '~/helpers/tree';

interface GatherStore {
  source: string;
  tree: Nodes;
  selected: Record<string, string>;
  location: Record<string, Array<string>>;
  sizes: Record<string, Sizes | null>;
  target: string;
  lastChecked: string;
  actions: {
    loadShares: () => Promise<void>;
    loadBranch: (node: Node) => Promise<void>;
    toggleSelected: (node: Node, shiftKey?: boolean) => Promise<void>;
    loadSize: (id: string, path: string) => Promise<void>;
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
    sizes: {},
    target: '',
    lastChecked: '',

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
      toggleSelected: async (node: Node, shiftKey = false) => {
        // shift-click extends the check/uncheck to every sibling between
        // the previously clicked node and this one
        const range = shiftKey
          ? getSiblingRange(get().lastChecked, node.id, get().tree)
          : null;
        const nodes = range ?? [node];
        const checked = !get().tree[node.id].checked;

        set((state) => {
          state.lastChecked = node.id;
        });

        for (const target of nodes) {
          // the tree may have changed while awaiting a previous locate
          const current = get().tree[target.id];
          if (!current || !!current.checked === checked) {
            continue;
          }

          set((state) => {
            state.tree[target.id].checked = checked;

            // remove parents by looping
            let parent = state.tree[target.parent];
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
            removeChildren(target);
          });

          if (!checked) {
            set((state) => {
              delete state.selected[target.id];
              delete state.location[target.id];
              state.selected = { ...state.selected };
              state.location = { ...state.location };
            });
            continue;
          }

          const fullPath = getAbsolutePath(target, get().tree);
          console.log('fullPath ', fullPath);

          const location = await Api.locate(fullPath);
          console.log('location ', location);

          set((state) => {
            state.selected[target.id] = fullPath;
            state.location[target.id] = location;
          });
        }
      },
      loadSize: async (id: string, path: string) => {
        // null marks an in-flight request, so each entry is fetched once
        if (get().sizes[id] !== undefined) {
          return;
        }

        set((state) => {
          state.sizes[id] = null;
        });

        const sizes = await Api.size(path);

        set((state) => {
          if (sizes === null) {
            // allow a retry on the next selection change
            delete state.sizes[id];
            return;
          }
          state.sizes[id] = sizes;
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
export const useGatherSizes = () => useGatherStore((state) => state.sizes);
export const useGatherTarget = () => useGatherStore((state) => state.target);
