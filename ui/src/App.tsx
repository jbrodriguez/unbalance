import React from 'react';

import { useLocation, useNavigate, Outlet } from 'react-router-dom';
import { ThemeProvider } from '@/components/theme-provider';
import ErrorBoundary from '@/components/error-boundary';
import { AuthGate } from '@/components/auth-gate';

import { Header } from './shared/header/header';
import { Footer } from './shared/footer/footer';
import { useConfigActions, useConfigVersion } from './state/config';
import {
  useAuthActions,
  useAuthenticated,
  useAuthEnabled,
  useAuthLoaded,
} from './state/auth';
import {
  useUnraidActions,
  useUnraidLoaded,
  useUnraidRoute,
} from './state/unraid';

export function App() {
  const { getConfig } = useConfigActions();
  const { load } = useAuthActions();
  const { setNavigate, getUnraid, syncRoute, connectSocket, disconnectSocket } =
    useUnraidActions();
  const isLoaded = useUnraidLoaded();
  const version = useConfigVersion();
  const route = useUnraidRoute();
  const authLoaded = useAuthLoaded();
  const authEnabled = useAuthEnabled();
  const authenticated = useAuthenticated();
  const navigate = useNavigate();
  const location = useLocation();

  React.useEffect(() => {
    load();
  }, [load]);

  React.useEffect(() => {
    console.log('setting navigation.,... ');
    setNavigate(navigate);
  }, [setNavigate, navigate]);

  React.useEffect(() => {
    // console.log('sync location >>>>>>>>>>>>> ', location);
    syncRoute(location.pathname);
  }, [location, syncRoute]);

  React.useEffect(() => {
    if (!authLoaded) {
      return;
    }

    if (authEnabled && !authenticated) {
      disconnectSocket();
      return;
    }

    connectSocket();
    getConfig();
    getUnraid();
  }, [
    authLoaded,
    authEnabled,
    authenticated,
    connectSocket,
    disconnectSocket,
    getConfig,
    getUnraid,
  ]);

  if (!authLoaded) {
    return null;
  }

  if (authEnabled && !authenticated) {
    return (
      <ThemeProvider defaultTheme="dark" storageKey="vite-ui-theme">
        <AuthGate />
      </ThemeProvider>
    );
  }

  if (!(isLoaded && version !== '')) {
    return null;
  }

  console.log('rendering App() ', route);

  return (
    <ThemeProvider defaultTheme="dark" storageKey="vite-ui-theme">
      <ErrorBoundary>
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
      </ErrorBoundary>
    </ThemeProvider>
  );
}
