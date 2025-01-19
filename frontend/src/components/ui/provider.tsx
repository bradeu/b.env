"use client";

import { config } from "@/web3";
import { ChakraProvider, defaultSystem } from "@chakra-ui/react";
import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import { WagmiProvider } from "wagmi";
import { ColorModeProvider, type ColorModeProviderProps } from "./color-mode";

const queryClient = new QueryClient();

export function Provider(props: ColorModeProviderProps) {
  return (
    <WagmiProvider config={config}>
      <QueryClientProvider client={queryClient}>
        <ChakraProvider value={defaultSystem}>
          <ColorModeProvider {...props} />
        </ChakraProvider>
      </QueryClientProvider>
    </WagmiProvider>
  );
}
