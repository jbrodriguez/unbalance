import React from 'react';

// import Tree, { BasicDataNode } from 'rc-tree';
// import 'rc-tree/assets/index.css';

import { Tree, Nodes, ITreeNode } from '~/shared/tree/tree2';
import { decorateNode } from '~/shared/tree/utils';
import { Icon } from '~/shared/icons/icon';
import { Api } from '~/api';
import { useScatterSource } from '~/state/scatter';
// import { findNode } from '~/helpers/bfs';

interface Props {
  height?: number;
}

// export type EventDataNode<BasicDataNode> = {
//   key: React.Key;
//   expanded: boolean;
//   selected: boolean;
//   checked: boolean;
//   loaded: boolean;
//   loading: boolean;
//   halfChecked: boolean;
//   dragOver: boolean;
//   dragOverGapTop: boolean;
//   dragOverGapBottom: boolean;
//   pos: string;
//   active: boolean;
// } & BasicDataNode;

// // const initialData = [{ title: 'mnt', key: '/', isLeaf: false }];

// type CheckedKeysType =
//   | { checked: React.Key[]; halfChecked: React.Key[] }
//   | React.Key[];

const initialNodes = {
  root: {
    id: 'root',
    label: '/',
    leaf: false,
    parent: '',
    children: [],
    // checked: false,
    // expanded: false,
    loading: false,
  },
};

export function FileSystem({ height }: Props) {
  const [nodes, setNodes] = React.useState<Nodes>(initialNodes);
  // const [checked, setChecked] = React.useState<CheckedKeysType>([]);
  const source = useScatterSource();

  // React.useEffect(() => {
  //   if (source === '') {
  //     return;
  //   }

  //   // @ts-expect-error ts-migrate(2531) FIXME: Object is possibly 'null'.
  //   setData([{ id: '1', label: source, value: `/mnt/${source}` }]);
  // }, [source]);

  // const onSelect = (info: unknown) => {
  //   console.log('selected', info);
  // };

  // const onCheck = (checkedKeys) => {
  //   console.log(checkedKeys);
  //   setChecked(checkedKeys);
  // };

  // const getTree = async (treeNode: EventDataNode<BasicDataNode>) => {
  //   // @ts-expect-error ts-migrate(2531) FIXME: Object is possibly 'null'.
  //   treeNode.children = [{ title: 'films', key: '/mnt/films', isLeaf: false }];
  //   return treeNode;
  // };

  // const onLoadData = async (treeNode: EventDataNode<BasicDataNode>) => {
  //   // console.log('load data... ', treeNode);
  //   // const loaded = await Api.getTree(`${treeNode.key}`);
  //   // console.log('loaded ', loaded);
  //   // const treeData = [...data];
  //   // // @ts-expect-error ts-migrate(2531) FIXME: Object is possibly 'null'.
  //   // const node = findNode(treeData, treeNode.key);
  //   // console.log('node ', node);
  //   // if (node === null) {
  //   //   return;
  //   // }
  //   // // @ts-expect-error ts-migrate(2531) FIXME: Object is possibly 'null'.
  //   // node.children = loaded as BasicDataNode;
  //   // setData(treeData);
  //   // // const index = treeData.findIndex((e) => e.key === treeNode.key);
  //   // // console.log('index ', index);
  //   // // if (index > -1) {
  //   // //   // @ts-expect-error ts-migrate(2531) FIXME: Object is possibly 'null'.
  //   // //   treeData[index].children = loaded as BasicDataNode;
  //   // //   console.log('treeData ', treeData);
  //   // // }
  // };

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
    // console.log('onLoad ', path, id);
    // const fullPath = `${source}/${path}`;
    // console.log('fullPath ', fullPath);
    // const tmp = await Api.getTree(fullPath, id);
    // console.log('loaded ', tmp);
    // await new Promise((r) => setTimeout(r, 3000));

    // const loaded = {
    //   '2': {
    //     id: '2',
    //     key: '/mnt/films',
    //     label: 'films',
    //     leaf: false,
    //     parent: '1',
    //     children: [],
    //     checked: false,
    //     expanded: false,
    //     loading: false,
    //   },
    //   '3': {
    //     id: '3',
    //     key: '/mnt/tvshows',
    //     label: 'tvshows',
    //     leaf: false,
    //     parent: '1',
    //     children: [],
    //     checked: false,
    //     expanded: false,
    //     loading: false,
    //   },
    // };
    // // nodes[id].children = Object.keys(loaded);
    // // return { ...nodes, ...loaded };
    // return loaded;

    // nodes[node.id].expanded = !nodes[node.id].expanded;
    // nodes[node.id].loading = true;
    // setNodes({ ...nodes });
    // const loaded = await onLoad(nodes, node.id);
    // console.log('loaded ', loaded);
    // nodes[node.id].loading = false;
    // setNodes(loaded);
    nodes[node.id].expanded = !nodes[node.id].expanded;

    if (isParent(node.id)) {
      setNodes({ ...nodes });
      return;
    }

    const draft = { ...nodes };
    draft.loader = {
      id: 'loader',
      // key: 'loader',
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
    // const branch = await onLoad(getAbsolutePath(node), node.id);
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
          // spinner={<Icon name="gift" />}
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

// import React from 'react';

// import Tree, { BasicDataNode } from 'rc-tree';
// import 'rc-tree/assets/index.css';

// import { Api } from '~/api';
// import { useScatterSource } from '~/state/scatter';
// import { findNode } from '~/helpers/bfs';

// interface Props {
//   height?: number;
// }

// export type EventDataNode<BasicDataNode> = {
//   key: React.Key;
//   expanded: boolean;
//   selected: boolean;
//   checked: boolean;
//   loaded: boolean;
//   loading: boolean;
//   halfChecked: boolean;
//   dragOver: boolean;
//   dragOverGapTop: boolean;
//   dragOverGapBottom: boolean;
//   pos: string;
//   active: boolean;
// } & BasicDataNode;

// // const initialData = [{ title: 'mnt', key: '/', isLeaf: false }];

// type CheckedKeysType =
//   | { checked: React.Key[]; halfChecked: React.Key[] }
//   | React.Key[];

// export function FileSystem({ height }: Props) {
//   const [data, setData] = React.useState<BasicDataNode[]>([]);
//   const [checked, setChecked] = React.useState<CheckedKeysType>([]);
//   const source = useScatterSource();

//   React.useEffect(() => {
//     if (source === '') {
//       return;
//     }

//     // @ts-expect-error ts-migrate(2531) FIXME: Object is possibly 'null'.
//     setData([{ title: source, key: `/mnt/${source}`, isLeaf: false }]);
//   }, [source]);

//   const onSelect = (info: unknown) => {
//     console.log('selected', info);
//   };

//   const onCheck = (checkedKeys: CheckedKeysType) => {
//     console.log(checkedKeys);
//     setChecked(checkedKeys);
//   };

//   // const getTree = async (treeNode: EventDataNode<BasicDataNode>) => {
//   //   // @ts-expect-error ts-migrate(2531) FIXME: Object is possibly 'null'.
//   //   treeNode.children = [{ title: 'films', key: '/mnt/films', isLeaf: false }];
//   //   return treeNode;
//   // };

//   const onLoadData = async (treeNode: EventDataNode<BasicDataNode>) => {
//     console.log('load data... ', treeNode);

//     const loaded = await Api.getTree(`${treeNode.key}`);
//     console.log('loaded ', loaded);
//     const treeData = [...data];
//     // @ts-expect-error ts-migrate(2531) FIXME: Object is possibly 'null'.
//     const node = findNode(treeData, treeNode.key);
//     console.log('node ', node);
//     if (node === null) {
//       return;
//     }

//     // @ts-expect-error ts-migrate(2531) FIXME: Object is possibly 'null'.
//     node.children = loaded as BasicDataNode;
//     setData(treeData);

//     // const index = treeData.findIndex((e) => e.key === treeNode.key);
//     // console.log('index ', index);
//     // if (index > -1) {
//     //   // @ts-expect-error ts-migrate(2531) FIXME: Object is possibly 'null'.
//     //   treeData[index].children = loaded as BasicDataNode;
//     //   console.log('treeData ', treeData);
//     // }
//   };

//   return (
//     <div className="flex flex-1 bg-neutral-200 dark:bg-gray-950">
//       <div className={`overflow-y-auto`} style={{ height: `${height}px` }}>
//         <Tree
//           onSelect={onSelect}
//           checkable
//           onCheck={onCheck}
//           checkedKeys={checked}
//           loadData={onLoadData}
//           treeData={data}
//         />
//       </div>
//     </div>
//   );
// }
