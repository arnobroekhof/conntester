# Poor mans Connection Tester

the program is able to test connectivity to multiple connections from a config file.

## Installation
download the latest binary release from the release page.

## usage:

```shell script
./conntester -config config.yaml
```

## output:

```
✅  --> able to connect to google = TCP:google.nl:443
❌  --> protocol UDP not supported
✅  --> able to connect to microsoft = TCP:microsoft.nl:443
```

## TODO list:

- [X] Test TCP Connection
- [ ] Test UDP Connection
- [ ] Ping ip or hostname
- [ ] Test DNS A record
- [ ] Test DNS PTR Record
- [ ] Test DNS SRV Record
- [ ] Test DNS CNAME Record