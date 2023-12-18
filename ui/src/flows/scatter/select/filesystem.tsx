import React from 'react';

import AutoSizer from 'react-virtualized-auto-sizer';

import { CheckboxTree } from '~/shared/tree/checkbox-tree';
import { Nodes, Node } from '~/types';
// import { decorateNode } from '~/shared/tree/utils';
import { Icon } from '~/shared/icons/icon';
// import { Api } from '~/api';
import { useScatterTree, useScatterActions } from '~/state/scatter';

interface Props {
  height?: number;
  width?: number;
}

// const rootNode = {
//   id: 'root',
//   label: '/',
//   leaf: false,
//   parent: '',
// };

export function FileSystem({ height }: Props) {
  // const source = useScatterSource();
  const tree = useScatterTree();
  const { loadBranch } = useScatterActions();
  const [nodes, setNodes] = React.useState<Nodes>({});

  // React.useEffect(() => {
  //   if (source === '') {
  //     return;
  //   }

  //   const root = decorateNode(rootNode as Node);
  //   setNodes({ root } as Nodes);
  // }, [source]);

  // const isParent = (id: string) =>
  //   Object.values(nodes).some((n) => n.parent === id);

  // const getAbsolutePath = (node: Node): string => {
  //   const parent = nodes[node.parent];
  //   if (!parent) {
  //     return node.label;
  //   }
  //   return `${getAbsolutePath(parent)}/${node.label}`;
  // };

  const onLoad = async (node: Node) => {
    await loadBranch(node);
    // nodes[node.id].expanded = !nodes[node.id].expanded;

    // if (isParent(node.id)) {
    //   setNodes({ ...nodes });
    //   return;
    // }

    // const draft = { ...nodes };
    // draft.loader = {
    //   id: 'loader',
    //   label: 'loading ...',
    //   leaf: false,
    //   parent: node.id,
    //   children: [],
    //   checked: false,
    //   expanded: false,
    //   loading: true,
    // };
    // draft[node.id].children = ['loader'];
    // setNodes(draft);

    // const route = `${source}/${getAbsolutePath(node)}`;
    // console.log('route ', route);
    // const branch = await Api.getTree(route, node.id);
    // console.log('loaded ', branch);
    // // await new Promise((r) => setTimeout(r, 1000));
    // for (const key in branch.nodes) {
    //   decorateNode(branch.nodes[key]);
    // }
    // delete draft.loader;
    // draft[node.id].children = branch.order;
    // setNodes({ ...draft, ...branch.nodes });
  };

  const onCheck = (node: Node) => {
    const draft = { ...nodes };
    draft[node.id].checked = !draft[node.id].checked;
    setNodes(draft);
  };

  return (
    <div className="flex flex-1 bg-neutral-200 dark:bg-gray-950">
      <AutoSizer disableHeight>
        {({ width }) => (
          <div
            className="overflow-y-auto overflow-x-auto p-4"
            style={{ height: `${height}px`, width: `${width}px` }}
          >
            <CheckboxTree
              nodes={tree}
              onLoad={onLoad}
              onCheck={onCheck}
              icons={{
                collapseIcon: (
                  <Icon
                    name="minus"
                    size={20}
                    fill="fill-slate-500 dark:fill-gray-700"
                  />
                ),
                expandIcon: (
                  <Icon
                    name="plus"
                    size={20}
                    fill="fill-slate-500 dark:fill-gray-700"
                  />
                ),
                checkedIcon: (
                  <Icon
                    name="checked"
                    size={20}
                    fill="fill-slate-700 dark:fill-slate-200"
                  />
                ),
                uncheckedIcon: (
                  <Icon
                    name="unchecked"
                    size={20}
                    fill="fill-slate-700 dark:fill-slate-200"
                  />
                ),
                leafIcon: (
                  <Icon
                    name="file"
                    size={20}
                    fill="fill-blue-400 dark:fill-blue-700"
                  />
                ),
                parentIcon: (
                  <Icon
                    name="folder"
                    size={20}
                    fill="fill-orange-400 dark:fill-yellow-300"
                  />
                ),
                hiddenIcon: (
                  <Icon
                    name="square"
                    size={20}
                    fill="fill-neutral-200 dark:fill-gray-950"
                  />
                ),
                loadingIcon: (
                  <Icon
                    name="loading"
                    size={20}
                    fill="animate-spin fill-slate-700 dark:fill-slate-700"
                  />
                ),
              }}
            />
          </div>
        )}
      </AutoSizer>
    </div>
  );
}
