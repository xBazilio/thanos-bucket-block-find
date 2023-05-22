# Find thanos blocks ULIDs by time range

For some tools you need block ULID to operate. This tool helps to find ULIDS using `thanos tools bucket web`.

```
usage: thanos-bucket-block-find --endpoint=ENDPOINT [<flags>]

Flags:
  --help               Show context-sensitive help (also try --help-long and --help-man).
  --min-time="0000-01-01T00:00:00Z"  
                       Start of time range to search blocks
  --max-time="9999-12-31T23:59:59Z"  
                       End of time range to search blocks
  --endpoint=ENDPOINT  Endpoint host of thanos bucket web to fetch data from
  --version            Show application version.
```
