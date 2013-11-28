package main
import(
    "fmt"
    "net/http"
    "log"
    "time"
    "bufio"
    "os"
    "strings"
    "encoding/json"
)

func StatData(path string) (map[string]map[string]map[string]float64, map[string]bool, map[string]bool) {
    file, err := os.Open(path)
    if err != nil {
        return nil, nil, nil
    }
    defer file.Close()

    stat := make(map[string]map[string]map[string]float64)
    scanner := bufio.NewScanner(file)
    head := []string{}
    labels := make(map[string]bool)
    features := make(map[string]bool)
    n := 0
    for scanner.Scan() {
        n += 1
        line := scanner.Text()
        tks := strings.Split(line, "\t")
        if n == 1{
            head = tks
        } else {
            label := tks[0]
            if len(label) == 0{
                continue
            }
            labels[label] = true
            for i, tk := range tks {
                if i == 0 {
                    continue
                } else {
                    if head[i][0] != '#'{
                        continue
                    }
                    key := head[i][1:]
                    features[key] = true
                    value := tk
                    _, ok := stat[key]
                    if !ok {
                        stat[key] = make(map[string]map[string]float64)
                    }
                    _, ok = stat[key][value]
                    if !ok {
                        stat[key][value] = make(map[string]float64)
                    }
                    _, ok = stat[key][value][label]
                    if !ok {
                        stat[key][value][label] = 1.0
                    } else {
                        stat[key][value][label] += 1.0
                    }
                }
            }
        }
    }
    return stat, labels, features
}

var global_stat map[string]map[string]map[string]float64
var global_labels map[string]bool
var global_features map[string]bool
 
func SingleFeatureStat( w http.ResponseWriter,r *http.Request ){
    w.Header().Set("Content-Type", "application/json")
    ret := []interface{}{}
    key := r.FormValue("feature")
    value_label_dis, ok := global_stat[key]
    
    if !ok {
        b, _ := json.Marshal(ret)
        fmt.Fprint(w, string(b))
        return
    }

    for label, _ := range global_labels {
        record := make(map[string]interface{})
        record["key"] = key + ": " + label
        values := []map[string]interface{}{}
        for value, _ := range global_features{
            label_dis, _ := global_stat[value]
            count, ok := label_dis[label]
            if !ok {
                count = 0.0
            }
            point := make(map[string]interface{})
            point["x"] = value
            point["y"] = count
            values = append(values, point)
        }
        record["values"] = values
        ret = append(ret, record)
    }
    b, _ := json.Marshal(ret)
    fmt.Fprint(w, string(b))
}

func FeatureList( w http.ResponseWriter,r *http.Request ){
    w.Header().Set("Content-Type", "application/json")
    ret := []string{}

    for key, _ := range global_stat {
        ret = append(ret, key)
    }
    b, _ := json.Marshal(ret)
    fmt.Fprint(w, string(b))
}
 
func main(){
    global_stat, global_labels, global_features = StatData("./data/adult/adult.tsv")
    fmt.Println(len(global_stat))
    http.Handle("/", http.FileServer(http.Dir(".")))
    http.HandleFunc( "/single_feature",SingleFeatureStat)
    http.HandleFunc("/feature_list", FeatureList)
    s := &http.Server{  
        Addr:           ":8080",
        ReadTimeout:    30 * time.Second,
        WriteTimeout:   30 * time.Second,
        MaxHeaderBytes: 1 << 20,
    }
    log.Fatal(s.ListenAndServe())
}