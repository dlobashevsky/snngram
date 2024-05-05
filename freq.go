package main

import (
  "os"
  "log"
  "bufio"
  "bytes"
  "strconv"
  "math"
  "sort"
)

type Score struct {
  Parts []string `json:"parts"`
  Val float64 `json:"val"`
}

// log10(frequency)
var ngrams = map[string] float64 {
}

// log10(frequency)
var freqs = map[int] float64 {
}

var total uint64

func FreqInit() {

  var ngrams2 = map[string] uint64 {
  }

  var freqs2 = map[int] uint64 {
  }

  total=0
  f, err := os.Open(CFG.Source)
  if err != nil {
    log.Panic("Cannot open file with frequencies (%s): %v",CFG.Source,err)
  }
  defer f.Close()

  freqs2=make(map[int] uint64)
  ngrams2=make(map[string] uint64)

  scanner:=bufio.NewScanner(f)
  scanner.Split(bufio.ScanLines)
  for scanner.Scan() {
    s:=scanner.Text()
    if len(s)==0 {
      continue
    }
    if bytes.IndexByte([]byte(" \t#"),s[0])>=0 {
      continue
    }

    before,after,found:=bytes.Cut([]byte(s),[]byte("\t"))
    if !found || len(before)==0 {
      continue
    }

    var c uint64
    c,err=strconv.ParseUint(string(after),10,64)
    if c==0 || err!=nil {
      continue
    }

    ngrams2[string(before)]+=c
    freqs2[len(before)]+=c
    total+=c

//    log.Printf("%s\t%d\n",before,c)
  }

  freqs=make(map[int] float64)
  ngrams=make(map[string] float64)

  for k,v:=range freqs2 {
    freqs[k]=math.Log10(float64(v))
  }

  for k,v:=range ngrams2 {
    ngrams[k]=math.Log10(float64(v))
  }

  log.Printf("importing done, total power is %d, records %d\n",total,len(ngrams))
}


func Split2(s string) []Score {

  idx:=0
  var ret []Score=make([]Score,CFG.Sys.Maxout)

  v,ok:=ngrams[s]
  f,fok:=freqs[len(s)]

  if ok && fok {
    ret[idx].Parts=make([]string,1)
    ret[idx].Parts[0]=s
    ret[idx].Val=f-v
    idx++
  }

  for i:=1;i<len(s);i++ {
    v1,ok1:=ngrams[s[:i]]
    v2,ok2:=ngrams[s[i:]]
    f1,fok1:=freqs[len(s[:i])]
    f2,fok2:=freqs[len(s[i:])]

    if ok1 && ok2 && fok1 && fok2 {
      ret[idx].Parts=make([]string,2)
      ret[idx].Parts[0]=s[:i]
      ret[idx].Parts[1]=s[i:]
      ret[idx].Val=f1+f2-v1-v2
      idx++
    }
  }
  rets:=ret[:idx]
  sort.SliceStable(rets, func(i,j int) bool {
    return rets[i].Val < rets[j].Val
  })
  return rets
}

func SplitAll(s string) (float64,bool) {

  v,ok:=ngrams[s]
  f,fok:=freqs[len(s)]

  if ok && fok {
    return f-v,true
  }

  if len(s)<=1 {
    return 1e100,false
  }

  minidx:=-1
  var minscore float64=1e100

  for i:=1;i<len(s);i++ {
    v1,ok1:=SplitAll(s[:i])
    v2,ok2:=SplitAll(s[i:])
    if ok1 && ok2 && v1+v2<minscore {
      minscore=v1+v2
      minidx=i
    }
  }

  return minscore,minidx>=0
}

