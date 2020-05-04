# lynis_exporter
Simple prometheus exporter for Lynis audit reports

## Install
```
go get -u github.com/MauveSoftware/lynis_exporter
```

## Configuration
This is a sample config to export 3 metrics from the recent report.

```yaml
reportfile_path: /path/to/report/file.dat
metrics:
  "lynis_hardening_index":
    "description": "Hardening index"
  "finish":
    "description": "1 if audit completed successfully"
    "name": "complete"
    "converter": "bool"
```

The output of /metrics will look like like:

```
lynis_vulnerable_packages_found 1
lynis_hardening_index 89
lynis_complete 1
```

## Usage

### Binary
```bash
./lynis_exporter
```

## License
(c) Mauve Mailorder Software GmbH & Co. KG, 2020. Licensed under [Apache 2.0](LICENSE) license.

## Prometheus
see https://prometheus.io/

## Lynis
see https://github.com/CISOfy/lynis