import { create } from 'zustand';
import { immer } from 'zustand/middleware/immer';

import { Nodes, Node } from '~/types';
import { Api } from '~/api';
import { decorateNode } from '~/shared/tree/utils';
import { isParent, getAbsolutePath } from '~/helpers/tree';

interface GatherStore {
  source: string;
  selected: Array<string>;
  // targets: Targets;
  tree: Nodes;
  actions: {
    loadShares: () => Promise<void>;
    loadBranch: (node: Node) => Promise<void>;
    toggleSelected: (node: Node) => void;
    // setSource: (source: string) => Promise<void>;
    // loadBranch: (node: Node) => Promise<void>;
    // toggleSelected: (node: Node) => void;
    // toggleTarget: (name: string) => void;
    // addLine: (line: string) => void;
  };
}

const rootNode = {
  id: 'root',
  label: '/',
  leaf: false,
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
    source: '/user',
    selected: [],
    // targets: {},
    tree: { root: decorateNode(rootNode as Node) },
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
    },
  })),
);

export const useGatherActions = () => useGatherStore((state) => state.actions);
export const useGatherTree = () => useGatherStore((state) => state.tree);
