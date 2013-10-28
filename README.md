envcfg
======

Lightweight Go configuration library that can read configuration variables from the environment.

Each non-blank line should take the form of one of the following directives:

```
# Settings["KEY_FLAG"] == envcfg.TRUE
KEY_FLAG
# Settings["KEY"] == LITERAL
KEY			LITERAL
# Settings["KEY"] == os.GetEnv(ENV_KEY). Throws error if env variable isn't set.
KEY			"ENV:"ENV_KEY
# Settings["KEY"] == os.GetEnv(ENV_KEY). Defaults to DEFAULT if env variable isn't set.
KEY			"ENV:"ENV_KEY	DEFAULT
# Lines beginning with octothorpe are comments
"#" COMMENT
```

You can then retrieve a Key-Value map in your code like so:

```go

cfgFile := os.Open(configFilePath)
settings, err := envcfg.ReadSettings(cfgFile)


```

note that ReadSettings takes io.Reader interface as input.
