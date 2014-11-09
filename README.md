kagenui
============

kagenui is profiler and analyzer for golang

Good for use with web application.

Usage
============

```
go get github.com/netmarkjp/kagenui
```


```
import "github.com/netmarkjp/kagenui"

func SomeHandler(){
    mp := kagenui.Begin("someHandler")
    defer mp.End()

    ...

    mp.Step("exec XXX")

    ...
    if valid {
        mp.AddMemo("cache is valid")
        ...
    } else {
        mp.AddMemo("cache is invalid")
        ...
    }

    ...

    mp.Step("exec YYY")

}

func profileHandler(w http.ResponseWriter, r *http.Request) {
    kagenui.Analyze(w)
    // kagenui.Dump(w)
    kagenui.Flush()
}
```

Analyze Output
============

```
Sort by Count
Count   Total  Mean     Stddev     Min     Max  Description
  100   24286   242     22.909     197     350  someHandler/Last Step to End
  100  144162  1441   1669.445     631   16671  someHandler/exec XXX
   50  257665  5153  27030.228     952  194161  someHandler/exec YYY

Sort by Total
Count   Total  Mean     Stddev     Min     Max  Description
   50  257665  5153  27030.228     952  194161  someHandler/exec XXX
  100  144162  1441   1669.445     631   16671  someHandler/exec YYY
  100   24286   242     22.909     197     350  someHandler/Last Step to End

Sort by Mean
Count   Total  Mean     Stddev     Min     Max  Description
   50  257665  5153  27030.228     952  194161  someHandler/exec YYY
  100  144162  1441   1669.445     631   16671  someHandler/exec XXX
  100   24286   242     22.909     197     350  someHandler/Last Step to End

Sort by Standard Deviation
Count   Total  Mean     Stddev     Min     Max  Description
   50  257665  5153  27030.228     952  194161  someHandler/exec YYY
  100  144162  1441   1669.445     631   16671  someHandler/exec XXX
  100   24286   242     22.909     197     350  someHandler/Last Step to End

Sort by Maximum(100 Percentile)
Count   Total  Mean     Stddev     Min     Max  Description
   50  257665  5153  27030.228     952  194161  someHandler/exec YYY
  100  144162  1441   1669.445     631   16671  someHandler/exec XXX
  100   24286   242     22.909     197     350  someHandler/Last Step to End
```

Dump Output
============

ltsv format

```
log:MP<TAB><DESCRIPTION_AT_STEP>:<SPENT_TIME_IN_NANOSEC><TAB>...<TAB>description:<DESCRIPTION_AT_BEGIN><TAB>memo:<MEMO1>,<MEMO2>,...
```

Customize
============

## Enable/Disable

```
kagenui.Enable()
```

```
kagenui.Disable()
```
