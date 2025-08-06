# ftdc-reader

Requirements:

* Go 1.24 or later

Run as follows:

```shell
make build
./dist/ftdc-reader /path/to/diagnostic.data/metrics.2025-08-05T11-12-59Z-00000 > ftdc.json
```

You can also pipe it to `mongoimport` and load FTDC data directly into a MongoDB instance, e.g.:

```shell
./dist/ftdc-reader /path/to/diagnostic.data/metrics.2025-08-05T11-12-59Z-00000 \
  | mongoimport --db ftdc --collection ftdc20250805
```
