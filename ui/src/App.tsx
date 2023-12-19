import React from 'react';

import { useLocation, useNavigate, Outlet } from 'react-router-dom';
// import { useLocation, Outlet } from 'react-router-dom';

import { Header } from './shared/header/header';
import { Footer } from './shared/footer/footer';
import { useConfigActions, useConfigVersion } from './state/config';
import {
  useUnraidActions,
  useUnraidLoaded,
  useUnraidStatus,
} from './state/unraid';
import { getRouteFromOp } from './helpers/steps';
// import { useUnraidStatus } from './state/unraid';
// import { getCurrentStep, getRoute } from './helpers/steps';

export function App() {
  const { getConfig } = useConfigActions();
  const { getUnraid, syncRouteAndStep } = useUnraidActions();
  const isLoaded = useUnraidLoaded();
  const version = useConfigVersion();
  const unraidStatus = useUnraidStatus();
  const navigate = useNavigate();
  const location = useLocation();

  React.useEffect(() => {
    // Google Analytics
    // ga('send', 'pageview');
    console.log('App.useEffect().synclocation ', location);
    syncRouteAndStep(location.pathname);
  }, [location, syncRouteAndStep]);

  React.useEffect(() => {
    getConfig();
    getUnraid();
  }, [getConfig, getUnraid]);

  React.useEffect(() => {
    if (!isLoaded) {
      return;
    }
    const route = getRouteFromOp(unraidStatus);
    console.log('routing ', unraidStatus, route);
    navigate(route);
  }, [unraidStatus, isLoaded, navigate]);

  // React.useEffect(() => {
  //   if (array === 'stopped') {
  //     navigate('/scatter');
  //   }
  // }, [navigate, unraidStatus]);

  // if (unraidStatus === OpStopped) {
  //   return null;
  // }
  // console.log('App() ', isLoaded, version);

  if (!(isLoaded && version !== '')) {
    return null;
  }

  console.log('rendering App() ', unraidStatus);

  return (
    <div className="container mx-auto h-screen flex flex-col">
      <header>
        <Header />
      </header>
      <main className="flex-1">
        <Outlet />
      </main>
      <footer>
        <Footer />
      </footer>
    </div>
  );
}
