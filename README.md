# CNC Drilling

Command to generate a G-Code to drill. It takes, as input, a dxf file.

Only points are taken into account.

```bash
go run ./cmd -d 10 -s 30 -z 20 ./testdata/points.dxf
```