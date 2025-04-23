package floods

import (
  "sync"
)

type Method struct {
  Type        int
  Name        string
  DisplayName       string
  Description string
  Subnet      int
  VIP         bool

  UsageCount  int
}

var (
  Methods = make(map[string]*Method)
  TotalUsage int
  mu sync.Mutex
)

func Get(name string) *Method {
  return Methods[name]
}

func Stats() (top *Method, vip, total, count int) {
  count = len(Methods)
  total = TotalUsage
  for _, m := range Methods {
      if m.VIP {
          vip++
      }
      if top == nil || m.UsageCount > top.UsageCount {
          top = m
      }
  }
  return
}