# snap-plugin-collector-goddd

[![Build Status](https://travis-ci.org/swhsiang/snap-plugin-collector-goddd.svg?branch=master)](https://travis-ci.org/swhsiang/snap-plugin-collector-goddd)

## Dev

```{shell}
glide install

make all
```

## Usage

```{shell}
# Remember to give GODDD_URL in the configuration file of goddd.
snaptel plugin load build/rootfs/snap-plugin-collector-goddd
snaptel task create -t examples/goddd-file.json
```

## Dev

### Generate json parser.

We use [easyjson](https://github.com/mailru/easyjson) to encode and decode the json object in our program.
If you want to change the structure of json object, you must to use easyjson to generate corresponding parser for json.

```{shell}
go get -u github.com/mailru/easyjson/...
cd goddd/
easyjson -all model.go
```

## Known issue

* Configuration
Since loading a plugin will trigger GetMetricTypes() and GetConfigPolicy(), there is not configuration file uploaded. 
The plugin can not have the value of endpoint which is specified in the configuration of plugin.
Reference https://github.com/intelsdi-x/snap/blob/master/docs/PLUGIN_LIFECYCLE.md#what-happens-when-a-plugin-is-loaded

