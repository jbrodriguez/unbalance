// import React from 'react';

import AutoSizer from 'react-virtualized-auto-sizer';

import { CheckboxTree } from '~/shared/tree/checkbox-tree';
import { Node } from '~/types';
import { Icon } from '~/shared/icons/icon';
import { useScatterTree, useScatterActions } from '~/state/scatter';

export const FileSystem: React.FunctionComponent = () => {
  const tree = useScatterTree();
  const { loadBranch, toggleSelected } = useScatterActions();

  const onLoad = async (node: Node) => await loadBranch(node);
  const onCheck = (node: Node) => toggleSelected(node);

  return (
    <div className="h-full flex flex-col bg-neutral-100 dark:bg-gray-950">
      <div className="flex flex-col pt-2 px-2">
        <h1 className="text-lg text-slate-500 dark:text-gray-500 pb-2">
          Folders/Files
        </h1>
        <hr className="border-slate-300 dark:border-gray-700" />
      </div>
      <div className="flex-auto">
        <AutoSizer disableWidth>
          {({ height }) => (
            <div
              className="overflow-y-auto overflow-x-auto p-2 text-slate-700 dark:text-gray-300"
              style={{ height: `${height}px` }}
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
                      style="fill-slate-500 dark:fill-gray-700"
                    />
                  ),
                  expandIcon: (
                    <Icon
                      name="plus"
                      size={20}
                      style="fill-slate-500 dark:fill-gray-700"
                    />
                  ),
                  checkedIcon: (
                    <Icon
                      name="checked"
                      size={20}
                      style="fill-green-700 dark:fill-lime-600"
                    />
                  ),
                  uncheckedIcon: (
                    <Icon
                      name="unchecked"
                      size={20}
                      style="fill-slate-700 dark:fill-slate-200"
                    />
                  ),
                  leafIcon: (
                    <Icon
                      name="file"
                      size={20}
                      style="fill-blue-400 dark:fill-gray-700"
                    />
                  ),
                  parentIcon: (
                    <Icon
                      name="folder"
                      size={20}
                      style="fill-orange-400 dark:fill-gray-700"
                    />
                  ),
                  hiddenIcon: (
                    <Icon
                      name="square"
                      size={20}
                      style="fill-neutral-100 dark:fill-gray-950"
                    />
                  ),
                  loadingIcon: (
                    <Icon
                      name="loading"
                      size={20}
                      style="animate-spin fill-slate-700 dark:fill-slate-700"
                    />
                  ),
                }}
              />
            </div>
          )}
        </AutoSizer>
      </div>
    </div>
  );
};
