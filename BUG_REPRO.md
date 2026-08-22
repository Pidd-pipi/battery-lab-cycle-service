# BUG_REPRO

## Bug 是什么
`OpsStore`/`OpsCache` 的 `Get`/`List`/`Update`/`Store` 返回或存下内部可变引用（`Labels` map 未深拷贝），调用方改写会污染存储；`Count`/`Size` 读路径无锁，与写操作并发时产生 data race。

## 如何触发
- 拿到返回的记录或列表后改写其标签，再次读取可见被污染；
- 并发执行 `Count()`（或 `Size()`）与 `Update()`/`Store()`。

## 真实错误信息
```
WARNING: DATA RACE
Write at 0x00c000123020 by goroutine 15:
  battery-lab-cycle-service.(*OpsStore).Update()
Previous read at 0x00c000123020 by goroutine 10:
  battery-lab-cycle-service.(*OpsStore).Count()
```
