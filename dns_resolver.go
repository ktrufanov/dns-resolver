package main

import (
    "os"
    "time"
    "strings"
    "regexp"
    "encoding/json"
    "github.com/rs/dnscache"
    "context"
    "net/http"
    "log"
)


/////////DNS Resolver

type Resolver struct {
    resolver *dnscache.Resolver
}

func AppendIfMissingStr(slice []string, i string) []string {
    for _, ele := range slice {
            if ele == i {
                    return slice
            }
    }
    return append(slice, i)
}

func IsIpv4Regex(ipAddress string) bool {
   ipRegex, _ := regexp.Compile(`^(([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])\.){3}([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])$`)
   ipAddress = strings.Trim(ipAddress, " ")
   return ipRegex.MatchString(ipAddress)
}

func (u *Resolver) Init() {

        if u.resolver == nil {
           u.resolver = &dnscache.Resolver{}
           go func() {
              t := time.NewTicker(2 * time.Minute)
              defer t.Stop()
              for range t.C {
                   u.resolver.Refresh(true)
              }
           }()
        }

        return
}


func (u *Resolver) Get(w http.ResponseWriter, r *http.Request) {

        r.ParseForm()
        w.Header().Set("Content-Type", "application/json; charset=utf-8")


       out_s := make(map[string][]string)
       out_s["Nets"] = []string{}
       if item, ok := r.Form["item"]; ok && item[0]!="" {
          addrs, err_ad := u.resolver.LookupHost(context.Background(), item[0])
          if err_ad==nil {
             for _,addr := range addrs {
                if IsIpv4Regex(addr) {
                   out_s["Nets"] = AppendIfMissingStr(out_s["Nets"], addr+"/32")
                }
             }
          }

       }
       js, _ := json.Marshal(out_s)
       w.Header().Set("Content-Type", "application/json; charset=utf-8")
       w.Write( []byte(js) )
}

func main() {
   resolver := new(Resolver)
   resolver.Init()
   http.HandleFunc( "/dns",func(w http.ResponseWriter, r *http.Request) { resolver.Get(w,r) } )
   port:=os.Getenv("SERVICE_PORT")
   if port == "" {
       port = "8080"
   }
   log.Fatal(http.ListenAndServe(":"+port, nil))
}

