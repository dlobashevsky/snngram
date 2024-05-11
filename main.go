package main


import (
  "log"
  "os"
  "os/signal"
  "runtime"
  "sync"
  "encoding/json"
  "net/http"
  "io"
  "io/ioutil"
)

const help="Usage:\t./split1 config_file.yaml"

type Jsplit2 struct {
  Src string `json:"src"`
  Cand []Score `json:"candidates"`
}

type Jscore struct {
  Src string `json:"src"`
  Val float64 `json:"score"`
}

func procSplit2(s string) string {
  log.Println(s)
  var req []string
  err:=json.Unmarshal([]byte(s),&req)
  if err != nil {
     return `{"error":"json parsing","data":[]}`
  }

  var resp []Jsplit2=make([]Jsplit2,len(req))

  for k,v := range req {
    resp[k].Src=v
    resp[k].Cand=Split2(v)
  }

  r,_:=json.Marshal(resp)
  return string(r) + "\n"
}


func getRoot(w http.ResponseWriter, r *http.Request) {
  log.Println("root requested")
    desc:=`<html><body><h1>short description</h1><p>
All methods take a json as input and produce json as output. Score is the inverse logarithm of estimation of the probability, i.e. score 7.8 mean probability 10^-7.8. So less score has a better probability. 
Estimation base is the google n-gram index.<hr><li>/split2 (POST) <u>Calculate scores for spliting word on two parts</u></li><br>Limitation: Words for the scoring should contain only lowercase ASCII whithout numbers 
and punctuations (regex: /^[a-z]+$/)<br>Request example:<br><b>["cardozaantonio","asdfasdf","andhereand"]</b><br>Response example:<br><b>[{"src":"cardozaantonio","candidates":[{"parts":["cardoza","antonio"],"val":7.472}]},...]
</b><br>Less score case is more frequent, results are sorted in ascending order.<hr>
<li>/score (POST)<u>Calculate scores for word</u></li><b>EXPERIMENTAL</b><br>Limitation: Words for the scoring should contain only lowercase ASCII whithout numbers and punctuations (regex: /^[a-z]+$/)<br>
Request example:<br><b>["cardozaantonio","asdfasdf","andhereand"]</b><br>Response example:<br><b>[{"src":"cardozaantonio","score"::7.472},...]</b><br>Less score is more frequent.<hr></body></html>`
    io.WriteString(w,desc)
}

func getSplit2(w http.ResponseWriter, r *http.Request) {
  if r.Method!="POST" {
    http.Error(w,"Invalid request method", http.StatusMethodNotAllowed)
    return
  }

  body, err:=ioutil.ReadAll(r.Body)
  if err != nil {
    http.Error(w,"Error reading request body",http.StatusInternalServerError)
  }
  w.Header().Set("Content-Type", "application/json")
  w.Header().Set("accept", "application/json")

  io.WriteString(w,procSplit2(string(body)))
}

func getScore(w http.ResponseWriter, r *http.Request) {
  if r.Method!="POST" {
    http.Error(w,"Invalid request method", http.StatusMethodNotAllowed)
    return
  }

  decoder := json.NewDecoder(r.Body)
  var req []string
  err := decoder.Decode(&req)
  if err != nil {
    http.Error(w,"Error parsing request body",http.StatusInternalServerError)
  }
  log.Println(req)

  var resp []Jscore=make([]Jscore,len(req))
  for k,v:=range req {
    resp[k].Src=v
    resp[k].Val=-1

    score,f:=SplitAll(v)
    if f {
      resp[k].Val=score
    }
  }
  json.NewEncoder(w).Encode(resp)
}


func main() {

//    flag.Parse()
    log.SetFlags(0)
    if len(os.Args)!=2 {
      log.Printf(help)
      os.Exit(1)
    }

    err:=Config_init(os.Args[1])
    if err!=nil {
      log.Println("Cant open config file",err)
      os.Exit(1)
    }

    {
      file, err := os.OpenFile(CFG.Log.File, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
      if err != nil {
        log.Fatal(err)
      }

      log.SetOutput(file)
    }

    FreqInit()

/*
//test
    log.Printf("expect 9.73241 %#v\n",Split2("cardozaantonio"))
    log.Printf("expect no candidates %#v\n",Split2("122134"))
    log.Printf("expect 4 cases (7.45, 13.75, 13.68, 8.25 %#v\n",Split2("andhereand"))

    log.Printf("result01 %s",procSplit2(`["andhereand","cardozaantonio","123412#$"]`))
    log.Printf("result02 %s",procSplit2(`["asdfqwerty","johndoe","dmitrylobashevsky","dimitrilobashevsky","dmytrolobashevskyi"]`))
    {
      f,t:=SplitAll("andhereand")
      log.Printf("result1 %v %v",f,t)
      f,t=SplitAll("cardozaantonio")
      log.Printf("result2 %v %v",f,t)
    }
os.Exit(0)
*/


    runtime.GOMAXPROCS(CFG.Sys.Threads)

    var wg sync.WaitGroup
    wg.Add(CFG.Sys.Threads)
//    MetricsInit()

// run http

    mux := http.NewServeMux()
    mux.HandleFunc("/", getRoot)
    mux.HandleFunc("/split2", getSplit2)
    mux.HandleFunc("/score", getScore)

    err=http.ListenAndServe(CFG.Listen.Service,mux)
    if err != nil {
      log.Fatal(err)
    }

    interrupt := make(chan os.Signal, 1)
    signal.Notify(interrupt, os.Interrupt)

    for {
        select {
        case <-interrupt:
            log.Println("interrupt")
            return
        }
    }
}
