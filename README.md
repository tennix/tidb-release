# TiDB release tools

## TiDB configurator

TiDB component configuration migration tool.

TiDB component configuration files are in TOML format, this tool helps to migrate old configuration to new configuration.

```
go build -o bin/tidb-configurator
./tidb-configurator -old-config=tikv.toml -new-config=new-tikv.toml -final-config=final-tikv.toml
```
