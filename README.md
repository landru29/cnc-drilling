# CNC Drilling

Command to generate a G-Code to drill or engrave. It takes, as input, a dxf file.

Only points are taken into account with command `drill`.

Only arcs and lines are taken into account with command `engrave`.

```bash
go run ./cmd drill -d 10 -f 30 -z 20 ./testdata/point01.dxf
go run ./cmd engrave -d 10 -f 30 -z 20 ./testdata/rectangle.dxf
```