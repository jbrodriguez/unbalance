import React from 'react';

export const Dashboard: React.FC = () => {
  return (
    <div className="grid grid-cols-5 gap-6 text-blue-600 ">
      <div className="border border-stroke dark:border-gray-800 border-slate-300 py-3 shadow-default dark:border-strokedark dark:bg-boxdark">
        <div className="flex items-end justify-between px-4">
          <div>
            <h1 className="text-3xl">52.33</h1>
            <span className="text-sm font-medium dark:text-slate-600 text-slate-600">
              Completed
            </span>
          </div>
          <span className="text-sm font-medium">%</span>
        </div>
      </div>

      <div className="border border-stroke dark:border-gray-800 border-slate-300 py-3 shadow-default dark:border-strokedark dark:bg-boxdark">
        <div className="flex items-end justify-between px-4">
          <div>
            <h1 className="text-3xl">42.38</h1>
            <span className="text-sm font-medium dark:text-slate-600 text-slate-600">
              Speed
            </span>
          </div>
          <span className="text-sm font-medium">MB/s</span>
        </div>
      </div>

      <div className="border border-stroke dark:border-gray-800 border-slate-300 py-3 shadow-default dark:border-strokedark dark:bg-boxdark">
        <div className="flex items-end justify-between px-4">
          <div>
            <h1 className="text-3xl">1.39</h1>
            <span className="text-sm font-medium dark:text-slate-600 text-slate-600">
              Transferred
            </span>
          </div>
          <span className="text-sm font-medium">GB</span>
        </div>
      </div>

      <div className="border border-stroke dark:border-gray-800 border-slate-300 py-3 shadow-default dark:border-strokedark dark:bg-boxdark">
        <div className="flex items-end justify-between px-4">
          <div>
            <h1 className="text-3xl">2.66</h1>
            <span className="text-sm font-medium dark:text-slate-600 text-slate-600">
              Total
            </span>
          </div>
          <span className="text-sm font-medium">GB</span>
        </div>
      </div>

      <div className="border border-stroke dark:border-gray-800 border-slate-300 py-3 shadow-default dark:border-strokedark dark:bg-boxdark">
        <div className="flex items-end justify-between px-4">
          <div>
            <h1 className="text-3xl">28</h1>
            <span className="text-sm font-medium dark:text-slate-600 text-slate-600">
              Remaining
            </span>
          </div>
          <span className="text-sm font-medium">s</span>
        </div>
      </div>
    </div>
  );
};
