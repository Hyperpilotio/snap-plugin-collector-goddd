# snap-plugin-collector-goddd

## Dev

```{shell}
glide install

make all
```

## Usage

```{shell}
snaptel plugin load build/rootfs/snap-plugin-collector-goddd
snaptel task create -t examples/goddd-file.json
```

## Known issue

* Configuration
Since loading a plugin will trigger GetMetricTypes() and GetConfigPolicy(), there is not configuration file uploaded. 
The plugin can not have the value of endpoint which is specified in the configuration of plugin.
Reference https://github.com/intelsdi-x/snap/blob/master/docs/PLUGIN_LIFECYCLE.md#what-happens-when-a-plugin-is-loaded
