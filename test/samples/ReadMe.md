# quick check one-liner

```bash
go install ../../... && ti_generate.exe -v -i trice.bin.sample -z 8 && cp idTable.c ../../src/ && go clean -cache && go install ../../... && ti_pack -v -i trice.bin.sample
```

As long `-z 4` gets better results than `-z 12` there is no best algorithm for _idTable.c_ generation.
