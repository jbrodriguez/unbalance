import React from 'react';

import Tree, { BasicDataNode } from 'rc-tree';
import 'rc-tree/assets/index.css';

import { Api } from '~/api';
import { useScatterSource } from '~/state/scatter';
import { findNode } from '~/helpers/bfs';

interface Props {
  height?: number;
}

export type EventDataNode<BasicDataNode> = {
  key: React.Key;
  expanded: boolean;
  selected: boolean;
  checked: boolean;
  loaded: boolean;
  loading: boolean;
  halfChecked: boolean;
  dragOver: boolean;
  dragOverGapTop: boolean;
  dragOverGapBottom: boolean;
  pos: string;
  active: boolean;
} & BasicDataNode;

// const initialData = [{ title: 'mnt', key: '/', isLeaf: false }];

type CheckedKeysType =
  | { checked: React.Key[]; halfChecked: React.Key[] }
  | React.Key[];

export function FileSystem({ height }: Props) {
  const [data, setData] = React.useState<BasicDataNode[]>([]);
  const [checked, setChecked] = React.useState<CheckedKeysType>([]);
  const source = useScatterSource();

  React.useEffect(() => {
    if (source === '') {
      return;
    }

    // @ts-expect-error ts-migrate(2531) FIXME: Object is possibly 'null'.
    setData([{ title: source, key: `/mnt/${source}`, isLeaf: false }]);
  }, [source]);

  const onSelect = (info: unknown) => {
    console.log('selected', info);
  };

  const onCheck = (checkedKeys: CheckedKeysType) => {
    console.log(checkedKeys);
    setChecked(checkedKeys);
  };

  // const getTree = async (treeNode: EventDataNode<BasicDataNode>) => {
  //   // @ts-expect-error ts-migrate(2531) FIXME: Object is possibly 'null'.
  //   treeNode.children = [{ title: 'films', key: '/mnt/films', isLeaf: false }];
  //   return treeNode;
  // };

  const onLoadData = async (treeNode: EventDataNode<BasicDataNode>) => {
    console.log('load data... ', treeNode);

    const loaded = await Api.getTree(`${treeNode.key}`);
    console.log('loaded ', loaded);
    const treeData = [...data];
    // @ts-expect-error ts-migrate(2531) FIXME: Object is possibly 'null'.
    const node = findNode(treeData, treeNode.key);
    console.log('node ', node);
    if (node === null) {
      return;
    }

    // @ts-expect-error ts-migrate(2531) FIXME: Object is possibly 'null'.
    node.children = loaded as BasicDataNode;
    setData(treeData);

    // const index = treeData.findIndex((e) => e.key === treeNode.key);
    // console.log('index ', index);
    // if (index > -1) {
    //   // @ts-expect-error ts-migrate(2531) FIXME: Object is possibly 'null'.
    //   treeData[index].children = loaded as BasicDataNode;
    //   console.log('treeData ', treeData);
    // }
  };

  return (
    <div className="flex flex-1 bg-neutral-200 dark:bg-gray-950">
      <div className={`overflow-y-auto`} style={{ height: `${height}px` }}>
        <Tree
          onSelect={onSelect}
          checkable
          onCheck={onCheck}
          checkedKeys={checked}
          loadData={onLoadData}
          treeData={data}
        />
      </div>
    </div>
  );
}
