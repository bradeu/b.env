import { http, createConfig, Config } from "wagmi";
import { base, mainnet } from "wagmi/chains";
import { injected, metaMask, safe, walletConnect } from "wagmi/connectors";

const projectId = "8277bfa3e76432ba39247d206ebb04e9";

export const config: Config = createConfig({
  chains: [mainnet, base],
  connectors: [injected(), walletConnect({ projectId }), metaMask(), safe()],
  transports: {
    [mainnet.id]: http(),
    [base.id]: http(),
  },
});
