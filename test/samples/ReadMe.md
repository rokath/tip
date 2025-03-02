# quick check one-liner

```bash
                         go install ../../... && ti_generate -v -i trice.bin.sample -z 8 && cp idTable.c ../../src/ && go clean -cache && go install ../../... && ti_pack -v -i trice.bin.sample

time (go clean -cache && go install ../../... && ti_generate -v -i ../../docs/TipUserManual.md -z 6 && cp idTable.c ../../src && go clean -cache && go install ../../... && ti_pack -v -i ../../docs/TipUserManual.md)
```

As long `-z 4` gets better results than `-z 12` there is no best algorithm for _idTable.c_ generation.
