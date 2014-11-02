miniprofiler
============

miniprofiler for golang

Good for use with web application.

Usage
============

```
go get github.com/netmarkjp/miniprofiler
```


```
import "github.com/netmarkjp/miniprofiler"

func SomeHandler(){
    mp := miniprofiler.Begin("someHandler")
    defer mp.End()

    ...

    mp.Step("exec XXX")

    ...

    mp.Step("exec YYY")

}

func ProfileHandler(){
    miniprofiler.Flush()
}
```

Outout
============

ltsv format

```
log:MP<TAB><DESCRIPTION_AT_STEP>:<SPENT_TIME_IN_NANOSEC><TAB>...<TAB>description:<DESCRIPTION_AT_BEGIN>
```

Customize
============

## Enable/Disable

```
miniprofiler.Enable()
```

```
miniprofiler.Disable()
```

## Change Writer

default: ``os.Stdout``

```
miniprofiler.SetWriter(writer)
```

