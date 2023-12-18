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
    setSource: (source: string) => void;
    loadBranch: (node: Node) => Promise<void>;
    toggleSelected: (id: string) => void;
    toggleTarget: (name: string) => void;
  };
}

const rootNode = {
  id: 'root',
  label: '/',
  leaf: false,
  parent: '',
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

export const useConfigStore = create<ScatterStore>()(
  immer((set, get) => ({
    source: '',
    selected: [],
    targets: {},
    tree: {},

    actions: {
      setSource: (source: string) => {
        set((state) => {
          state.source = source;
          state.targets = {};
          state.tree = { root: rootNode as Node };
          state.selected = [];
        });
      },
      loadBranch: async (node: Node) => {
        set((state) => {
          state.tree[node.id].expanded = !state.tree[node.id].expanded;
        });

        const draft = get();

        if (isParent(node.id, draft.tree)) {
          set((state) => {
            state.tree = { ...draft.tree };
          });
          return;
        }

        // draft = {
        //   ...draft,
        //   loader: {
        //     id: 'loader',
        //     label: 'loading ...',
        //     leaf: false,
        //     parent: node.id,
        //     children: [],
        //     checked: false,
        //     expanded: false,
        //     loading: true,
        //   },
        //   [node.id]: {
        //     ...node,
        //     children: ['loader'],
        //   },
        // };
        // draft[node.id].children = ['loader'];
        // draft.tree.loader = {
        //   id: 'loader',
        //   label: 'loading ...',
        //   leaf: false,
        //   parent: node.id,
        //   children: [],
        //   checked: false,
        //   expanded: false,
        //   loading: true,
        // };

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

        console.log('draft ', draft);

        const route = `${draft.source}/${getAbsolutePath(node, draft.tree)}`;
        console.log('route ', route);
        const branch = await Api.getTree(route, node.id);
        // console.log('loaded ', branch);
        // // await new Promise((r) => setTimeout(r, 1000));
        for (const key in branch.nodes) {
          decorateNode(branch.nodes[key]);
        }

        set((state) => {
          delete state.tree.loader;
          state.tree = { ...state.tree, ...branch.nodes };
          state.tree[node.id].children = branch.order;
        });

        // delete draft.tree.loader;
        // draft[node.id].children = branch.order;
        // set((state) => {
        //   state.tree = {
        //     ...draft,
        //     [node.id]: { ...node, children: branch.order },
        //     ...branch.nodes,
        //   };
        // });
      },

      // const onLoad = async (node: Node) => {

      //   const route = `${source}/${getAbsolutePath(node)}`;
      //   console.log('route ', route);
      //   const branch = await Api.getTree(route, node.id);
      //   console.log('loaded ', branch);
      //   // await new Promise((r) => setTimeout(r, 1000));
      //   for (const key in branch.nodes) {
      //     decorateNode(branch.nodes[key]);
      //   }
      //   delete draft.loader;
      //   draft[node.id].children = branch.order;
      //   setNodes({ ...draft, ...branch.nodes });
      // };

      toggleSelected: (id: string) => {
        set((state) => {
          state.tree[id].checked = !state.tree[id].checked;
          const index = state.selected.indexOf(id);
          if (index === -1) {
            state.selected.push(id);
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

export const useScatterActions = () => useConfigStore((state) => state.actions);

export const useScatterSource = () => useConfigStore((state) => state.source);
export const useScatterTree = () => useConfigStore((state) => state.tree);
export const useScatterTargets = () => useConfigStore((state) => state.targets);
