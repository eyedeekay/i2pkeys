i2pkeys
=======

Generates and displays the contents of files that are storing i2p keys in the
incompatible format used for sam3

### Verbosity
Logging can be enabled and configured using the DEBUG_I2P environment variable. By default, logging is disabled.

There are three available log levels:

- Debug
```shell
export DEBUG_I2P=debug
```
- Warn
```shell
export DEBUG_I2P=warn
```
- Error
```shell
export DEBUG_I2P=error
```

If in case I2P_DEBUG is set to an unrecognized variable, it will fall back to "debug".