import { StrictMode } from "react";
import { createRoot } from "react-dom/client";
import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import App from "./App";
import { Provider } from "@/components/ui/provider";
import { Layout } from "./components/ui/Layout";
import { BrowserRouter } from "react-router-dom";
import { Toaster } from "./components/ui/toaster";

const queryClient = new QueryClient({
  defaultOptions: {
    queries: {
      retry: 1,
      refetchOnWindowFocus: false,
    },
  },
});

createRoot(document.getElementById("root")!).render(
  <StrictMode>
    <QueryClientProvider client={queryClient}>
      <BrowserRouter>
        <Provider>
          <Layout>
            <App />
            <Toaster />
          </Layout>
        </Provider>
      </BrowserRouter>
    </QueryClientProvider>
  </StrictMode>
);
