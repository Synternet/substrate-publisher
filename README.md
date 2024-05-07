# Substrate publisher
[![Latest release](https://img.shields.io/github/v/release/Synternet/substrate-publisher)](https://github.com/Synternet/substrate-publisher/releases/latest)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![GitHub Workflow Status (with event)](https://img.shields.io/github/actions/workflow/status/Synternet/substrate-publisher/github-ci.yml?label=github-ci)](https://github.com/Synternet/substrate-publisher/actions/workflows/github-ci.yml)

Establishes connection with Substrate blockchain based node and publishes blockchain data to Synternet via NATS connection.

# Supported blockchains

This repository uses [GSRPC](https://github.com/centrifuge/go-substrate-rpc-client) ensuring compatibility with all substrate based chains. Tested blockchain: Polkadot, peaq.

## Polkadot

Polkadot streams available in [testnet Synternet portal](https://datalayer.synternet.com/subscribe/amber1x64mphk6fx8xrcnxn3ynepsqhv446uhp0k77z4/AAWG2YVSOTUW5RKT2JCOHWDHBV3UF4DUBZOOBMOPHH5VGSECAGROWBVI/)

## peaq

peaq specific custom `PeaqStorage` events `ItemAdded`, `ItemRead`, `ItemUpdated` are decoded. Available in [testnet Synternet portal](https://datalayer.synternet.com/subscribe/amber1x64mphk6fx8xrcnxn3ynepsqhv446uhp0k77z4/AADZCLQXAARU4JYV4ZEQ3ZZUBNCSTPZSJVSMP6AU5UJNJ2HUOIEONW2R/).

# Usage

Building from source
```bash
make build
```

Running executable
```bash
./dist/substrate-publisher start --rpc-url wss://rpc.polkadot.io --nats-nkey SA..BC
```

### Environment variables and flags

Environment variables can be passed to docker container. Flags can be passed as executable arguments.

| Environment variable   | Flag                   | Description                                                                                                               |
| ---------------------- | ---------------------- | ------------------------------------------------------------------------------------------------------------------------- |
| RPC_URL                | rpc-url                | Substrate node URI to establish RPC connection, e.g.: `wss://rpc.polkadot.io`                                             |
| NATS_URLS              | nats-urls              | NATS connection URLs to Synternet brokers, e.g.: `nats://e.f.g.h`. URL to [broker](https://docs.synternet.com/docs/actors/broker). Default: testnet. |
| NATS_NKEY              | nats-nkey              | NATS account NKEY a.k.a access token, e.g.: `SA..SI` (58 chars). See [here](https://docs.synternet.com/build/data-layer/developer-portal/publish-streams#7-get-the-access-token). |
| STREAM_PREFIX          | stream-prefix          | Stream prefix, e.g.: `foo` prefix results in `foo.substrate.<tx,log-even,header,...>` stream subjects. Stream prefix should be same as registered wallet [alias](https://docs.synternet.com/build/data-layer/developer-portal/publish-streams#2-register-a-wallet---get-your-alias). |
| STREAM_PUBLISHER_INFIX | stream-publisher-infix | (optional) Stream publisher infix, e.g.: `foo` infix results in `prefix.foo.<tx,log-even,header,...>` stream subjects. Stream publisher infix should be same as registered publisher [alias](https://docs.synternet.com/build/data-layer/developer-portal/publish-streams#3-register-a-publisher). Default: `substrate`. |
| STREAM_NETWORK_INFIX   | stream-network-infix   | (optional) Specify stream network infix, e.g.: `mainnet` prefix results in `<prefix>.substrate.mainnet.<tx,...>` subjects. Default: empty (`prefix.substrate.<tx,...>`). |

See [Data Layer Quick Start](https://docs.synternet.com/build/data-layer/data-layer-quick-start) to learn more.

## Docker

### Build from source

1. Build image.
```
docker build -f ./build/Dockerfile -t substrate-publisher .
```

2. Run container with passed environment variables.
```
docker run -it --rm --env-file=.env substrate-publisher
```

### Prebuilt image

Run container with passed environment variables.
```
docker run -it --rm --env-file=.env ghcr.io/synternet/substrate-publisher:latest
```

### Docker Compose

`docker-compose.yml` file.
```
version: '3.8'

services:
  substrate-publisher:
    image: ghcr.io/synternet/substrate-publisher:latest
    environment:
      - RPC_URL=wss://rpc.polkadot.io
      - NATS_NKEY=secret-access-token
      - STREAM_PREFIX=your-org
      - STREAM_PUBLISHER_INFIX=substrate-based-chain
```

## Contributing

We welcome contributions from the community. Whether it's a bug report, a new feature, or a code fix, your input is valued and appreciated.

## Synternet

If you have any questions, ideas, or simply want to connect with us, we encourage you to reach out through any of the following channels:

- **Discord**: Join our vibrant community on Discord at [https://discord.com/invite/jqZur5S3KZ](https://discord.com/invite/jqZur5S3KZ). Engage in discussions, seek assistance, and collaborate with like-minded individuals.
- **Telegram**: Connect with us on Telegram at [https://t.me/Synternet](https://t.me/Synternet). Stay updated with the latest news, announcements, and interact with our team members and community.
- **Email**: If you prefer email communication, feel free to reach out to us at devrel@synternet.com. We're here to address your inquiries, provide support, and explore collaboration opportunities.
