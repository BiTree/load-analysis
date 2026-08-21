# load-analysis

这是一个用 Go 实现的 Linux 内核风格负载均值计算器。它**考虑了内核的指数平滑计算方法**，而不是对最近样本做普通算术平均。

## 计算方法

Linux 负载均值使用指数加权移动平均（EWMA）：

```text
load(t) = load(t-1) * exp + active(t) * (1 - exp)
```

实现与内核算法保持一致：

- 每 5 秒输入一次活跃任务数；
- 使用 11 位小数的定点数（`FIXED_1 = 2048`），避免浮点误差；
- 使用内核预计算的 1、5、15 分钟衰减系数 `1884`、`2014`、`2037`；
- 当负载上升时采用与内核 `calc_load` 相同的向上取整方式。

这里的 `active` 应是该采样时刻处于可运行状态或不可中断睡眠状态的任务数。调用者必须每 5 秒调用一次 `Update`；空闲采样也要传入 `0`，否则时间衰减会不正确。

## 作为 Go 包使用

```go
package main

import (
	"fmt"

	"load-analysis/loadavg"
)

func main() {
	var avg loadavg.Averages
	_ = avg.Update(2)
	one, five, fifteen := avg.Values()
	fmt.Printf("%.2f %.2f %.2f\n", one, five, fifteen)
}
```

如果需要逐位核对内核结果，可以使用 `FixedValues` 获取原始 Q11 定点数。

## 命令行使用

每行（或每个空白分隔字段）提供一个 5 秒采样的活跃任务数：

```sh
printf '2\n2\n2\n' | go run ./cmd/load-analysis
```

输出的第一列是样本序号，之后依次是 1、5、15 分钟负载均值：

```text
1       0.16 0.03 0.01
2       0.31 0.07 0.02
3       0.44 0.10 0.03
```
