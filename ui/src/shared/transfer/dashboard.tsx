import React from 'react';

export const Dashboard: React.FC = () => {
  return (
    <div className="grid grid-cols-5 gap-6 text-blue-600 ">
      <div className="border border-stroke dark:border-gray-800 border-slate-300 px-4 py-3 shadow-default dark:border-strokedark dark:bg-boxdark">
        <div className="flex justify-between items-center">
          <span className="text-3xl">52.33</span>
          <span>%</span>
        </div>
        <span className="text-sm font-medium dark:text-slate-600 text-slate-600">
          Completed
        </span>
      </div>

      <div className="border border-stroke dark:border-gray-800 border-slate-300 px-4 py-3 shadow-default dark:border-strokedark dark:bg-boxdark">
        <div className="flex justify-between items-center">
          <span className="text-3xl">42.38</span>
          <span>MB/s</span>
        </div>
        <span className="text-sm font-medium dark:text-slate-600 text-slate-600">
          Speed
        </span>
      </div>

      <div className="border border-stroke dark:border-gray-800 border-slate-300 px-4 py-3 shadow-default dark:border-strokedark dark:bg-boxdark">
        <div className="flex justify-between items-center">
          <span className="text-3xl">1.39</span>
          <span>GB</span>
        </div>
        <span className="text-sm font-medium dark:text-slate-600 text-slate-600">
          Transferred
        </span>
      </div>

      <div className="border border-stroke dark:border-gray-800 border-slate-300 px-4 py-3 shadow-default dark:border-strokedark dark:bg-boxdark">
        <div className="flex justify-between items-center">
          <span className="text-3xl">2.66</span>
          <span>GB</span>
        </div>
        <span className="text-sm font-medium dark:text-slate-600 text-slate-600">
          Total
        </span>
      </div>

      <div className="border border-stroke dark:border-gray-800 border-slate-300 px-4 py-3 shadow-default dark:border-strokedark dark:bg-boxdark">
        <div className="flex justify-between items-center">
          <span className="text-3xl">28</span>
          <span>sec</span>
        </div>
        <span className="text-sm font-medium dark:text-slate-600 text-slate-600">
          Remaining
        </span>
      </div>
    </div>
  );
};
