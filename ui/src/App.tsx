import React from 'react';

import { useLocation, useNavigate, Outlet } from 'react-router-dom';
import { ThemeProvider } from '@/components/theme-provider';

import { Header } from './shared/header/header';
import { Footer } from './shared/footer/footer';
import { useConfigActions, useConfigVersion } from './state/config';
import {
  useUnraidActions,
  useUnraidLoaded,
  useUnraidRoute,
} from './state/unraid';

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
    // console.log('sync location >>>>>>>>>>>>> ', location);
    syncRoute(location.pathname);
  }, [location, syncRoute]);

  React.useEffect(() => {
    getConfig();
    getUnraid();
  }, [getConfig, getUnraid]);

  if (!(isLoaded && version !== '')) {
    return null;
  }

  console.log('rendering App() ', route);

  return (
    <ThemeProvider defaultTheme="dark" storageKey="vite-ui-theme">
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
    </ThemeProvider>
  );
}
