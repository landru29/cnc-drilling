# CNC Drilling

Command to generate a G-Code to drill or engrave. It takes, as input, a dxf file.

Only points are taken into account with command `drill`.

Only arcs and lines are taken into account with command `engrave`.

```bash
go run ./cmd drill -d 10 -f 30 -z 20 ./testdata/point01.dxf
go run ./cmd engrave -d 10 -f 30 -z 20 ./testdata/rectangle.dxf
```

## Configuration

Some parameters can be set in a config file. The config file is looked for in the following order:
* `./drill.yaml`
* `$HOME/.cnc-drilling/drill.yaml`
* `/etc/cnc-drilling/drill.yaml`

Use the following command to generate `./drill.yaml` as a model:

```bash
go run ./cmd save-config
```

## Environment variables

The following environment variables can be set:

