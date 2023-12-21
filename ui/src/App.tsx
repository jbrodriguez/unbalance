import React from 'react';

import { useLocation, useNavigate, Outlet } from 'react-router-dom';
// import { useLocation, Outlet } from 'react-router-dom';

import { Header } from './shared/header/header';
import { Footer } from './shared/footer/footer';
import { useConfigActions, useConfigVersion } from './state/config';
import {
  useUnraidActions,
  useUnraidLoaded,
  useUnraidRoute,
} from './state/unraid';
// import { getRouteFromStatus } from './helpers/routes';
// import { useUnraidStatus } from './state/unraid';
// import { getCurrentStep, getRoute } from './helpers/steps';

export function App() {
  const { getConfig } = useConfigActions();
  const { setNavigate, getUnraid, syncRoute } = useUnraidActions();
  const isLoaded = useUnraidLoaded();
  const version = useConfigVersion();
  const route = useUnraidRoute();
  const navigate = useNavigate();
  const location = useLocation();

  React.useEffect(() => {
    console.log('setting navigation.,... ');
    setNavigate(navigate);
  }, [setNavigate, navigate]);

  React.useEffect(() => {
    // Google Analytics
    // ga('send', 'pageview');
    console.log('sync location >>>>>>>>>>>>> ', location);
    syncRoute(location.pathname);
  }, [location, syncRoute]);

  React.useEffect(() => {
    getConfig();
    getUnraid();
  }, [getConfig, getUnraid]);

  // React.useEffect(() => {
  //   if (!isLoaded) {
  //     return;
  //   }
  //   if (location.pathname === route) {
  //     return;
  //   }

  //   // const route = getRouteFromStatus(unraidStatus);
  //   console.log('routing ', route);
  //   navigate(route);
  // }, [route, location, isLoaded, navigate]);
  // React.useEffect(() => {
  //   if (!isLoaded) {
  //     return;
  //   }

  //   // const route = getRouteFromStatus(unraidStatus);
  //   console.log('routing ', route);
  //   navigate(route);
  // }, [route, isLoaded, navigate]);

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

  console.log('rendering App() ', route);

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
