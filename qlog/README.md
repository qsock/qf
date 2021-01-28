# qlog

## 创建日志
```golang
config := qlog.Config{
    Name: "default.log",
    Dir: "/tmp",
    Level: "info",
}
logger := config.Build()
logger.SetLevel(qlog.DebugLevel)
logger.Debug("debug", qlog.String("a", "b"))
logger.Debugf("debug %s", "a")
logger.Debugw("debug", "a", "b")
```

