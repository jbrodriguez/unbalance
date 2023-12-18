import React from 'react';

import { Tree, Nodes, ITreeNode } from '~/shared/tree/tree2';
import { decorateNode } from '~/shared/tree/utils';
import { Icon } from '~/shared/icons/icon';
import { Api } from '~/api';
import { useScatterSource } from '~/state/scatter';

interface Props {
  height?: number;
}

const initialNodes = {
  root: {
    id: 'root',
    label: '/',
    leaf: false,
    parent: '',
    children: [],
    loading: false,
  },
};

export function FileSystem({ height }: Props) {
  const [nodes, setNodes] = React.useState<Nodes>(initialNodes);
  const source = useScatterSource();

  // React.useEffect(() => {
  //   if (source === '') {
  //     return;
  //   }

  //   // @ts-expect-error ts-migrate(2531) FIXME: Object is possibly 'null'.
  //   setData([{ id: '1', label: source, value: `/mnt/${source}` }]);
  // }, [source]);

  const isParent = (id: string) =>
    Object.values(nodes).some((n) => n.parent === id);

  const getAbsolutePath = (node: ITreeNode): string => {
    const parent = nodes[node.parent];
    if (!parent) {
      return node.label;
    }
    return `${getAbsolutePath(parent)}/${node.label}`;
  };

  const onLoad = async (node: ITreeNode) => {
    nodes[node.id].expanded = !nodes[node.id].expanded;

    if (isParent(node.id)) {
      setNodes({ ...nodes });
      return;
    }

    const draft = { ...nodes };
    draft.loader = {
      id: 'loader',
      label: 'loading ...',
      leaf: false,
      parent: node.id,
      children: [],
      checked: false,
      expanded: false,
      loading: true,
    };
    draft[node.id].children = ['loader'];
    setNodes(draft);

    const route = `${source}/${getAbsolutePath(node)}`;
    console.log('route ', route);
    const branch = await Api.getTree(route, node.id);
    console.log('loaded ', branch);
    await new Promise((r) => setTimeout(r, 1000));
    for (const key in branch.nodes) {
      decorateNode(branch.nodes[key]);
    }
    delete draft.loader;
    draft[node.id].children = branch.order;
    setNodes({ ...draft, ...branch.nodes });
  };

  const onCheck = (node: ITreeNode) => {
    const draft = { ...nodes };
    draft[node.id].checked = !draft[node.id].checked;
    setNodes(draft);
  };

  return (
    <div className="flex flex-1 bg-neutral-200 dark:bg-gray-950">
      <div className="overflow-y-auto p-4" style={{ height: `${height}px` }}>
        <Tree
          nodes={nodes}
          onLoad={onLoad}
          onCheck={onCheck}
          collapseIcon={
            <Icon
              name="minus"
              size={20}
              fill="fill-slate-500 dark:fill-slate-500"
            />
          }
          expandIcon={
            <Icon
              name="plus"
              size={20}
              fill="fill-slate-500 dark:fill-slate-500"
            />
          }
          checkedIcon={
            <Icon
              name="checked"
              size={20}
              fill="fill-slate-700 dark:fill-slate-200"
            />
          }
          uncheckedIcon={
            <Icon
              name="unchecked"
              size={20}
              fill="fill-slate-700 dark:fill-slate-200"
            />
          }
          leafIcon={
            <Icon
              name="file"
              size={20}
              fill="fill-blue-400 dark:fill-blue-700"
            />
          }
          parentIcon={
            <Icon
              name="folder"
              size={20}
              fill="fill-orange-400 dark:fill-yellow-300"
            />
          }
          placeholderIcon={
            <Icon
              name="placeholder"
              size={20}
              fill="fill-neutral-200 dark:fill-gray-950"
            />
          }
          loadingIcon={
            <Icon
              name="loading"
              size={20}
              fill="animate-spin fill-slate-700 dark:fill-slate-700"
            />
          }
        />
      </div>
    </div>
  );
}
