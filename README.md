# Docker machine driver for [Vscale Cloud](https://vscale.io)

## Installation

### Easy way
Download the release bundle from the releases section and place the binary that corresponds to your platform it somewhere in your PATH

### Hard way
Use go get github.com/eduardnikolenko/docker-machine-driver-vscale and make sure that docker-machine-driver-vscale is located somwhere in your PATH

## Usage

    $ docker-machine create \
      --driver vscale \
      --vscale-access-token=<YOU_ACCESS_TOKEN> \
      my-machine

## Options

| Parameter                   | Env                   | Default |
| --------------------------- | --------------------- | ------- |
| **`--vscale-access-token`** | `VSCALE_ACCESS_TOKEN` | -       |
| **`--vscale-location`**     | `VSCALE_LOCATION`     | `spb0`  |
| **`--vscale-made-from`**    | `VSCALE_MADE_FROM`    |         |
| **`--vscale-rplan`**        | `VSCALE_RPLAN`        | `small` |
| **`--vscale-swap-file`**    | `VSCALE_SWAP_FILE`    | `0`     |

## License

MIT © [Eduard Nikolenko](https://github.com/eduardnikolenko)
