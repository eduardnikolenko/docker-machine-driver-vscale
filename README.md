# Docker machine driver for [Vscale Cloud](https://vscale.io)

[![Go Report Card](https://goreportcard.com/badge/github.com/eduardnikolenko/docker-machine-driver-vscale)](https://goreportcard.com/report/github.com/eduardnikolenko/docker-machine-driver-vscale)
[![License](https://img.shields.io/badge/License-MIT-blue.svg)](/LICENSE)

## Installation

Use go get github.com/eduardnikolenko/docker-machine-driver-vscale and make sure that docker-machine-driver-vscale is located somwhere in your PATH

## Usage

    $ docker-machine create \
      --driver vscale \
      --vscale-access-token=<YOU_ACCESS_TOKEN> \
      my-machine

## Options

| Parameter                   | Env                   | Default                      |
| --------------------------- | --------------------- | ---------------------------- |
| **`--vscale-access-token`** | `VSCALE_ACCESS_TOKEN` | -                            |
| **`--vscale-location`**     | `VSCALE_LOCATION`     | `spb0`                       |
| **`--vscale-made-from`**    | `VSCALE_MADE_FROM`    | `ubuntu_16.04_64_001_docker` |
| **`--vscale-rplan`**        | `VSCALE_RPLAN`        | `small`                      |
| **`--vscale-swap-file`**    | `VSCALE_SWAP_FILE`    | `0`                          |

## License

MIT © [Eduard Nikolenko](https://github.com/eduardnikolenko)
