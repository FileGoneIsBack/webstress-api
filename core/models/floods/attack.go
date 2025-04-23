package floods

import "time"

type Attack struct {
	ID     int
	Target string
	*Method
	Port, Threads, PPS, Subnet int
	Parent, Stopped, Conns     int
	Duration                   int
	Created                    int64
}


func New(method string) *Attack {
    m := Get(method)
    if m == nil {
        return nil
    }

    // bump usage counters
    mu.Lock()
    m.UsageCount++
    TotalUsage++
    mu.Unlock()

    return &Attack{
        Target:   "",
        Duration: 0,
        Port:     0,
        Threads:  3,
        PPS:      250000,
        Method:   m,
        Created:  time.Now().Unix(),
        Parent:   0,
    }
}