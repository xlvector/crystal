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
    "strconv"
)

type WeightLabel struct {
    weight float64
    label string
}

type WeightLabelList []WeightLabel

func (self *WeightLabelList) MaxWeight() float64 {
    ret := (*self)[0].weight
    for _, wl := range *self {
        if ret < wl.weight {
            ret = wl.weight
        }
    }
    return ret
}

type StatResult struct {
    discrete_stat map[string]map[string]map[string]float64
    labels map[string]bool
    features map[string]bool
    continuous_stat map[string]WeightLabelList
}

func (self *StatResult) Init() {
    self.discrete_stat = make(map[string]map[string]map[string]float64)
    self.labels = make(map[string]bool)
    self.features = make(map[string]bool)
    self.continuous_stat = make(map[string]WeightLabelList)
}

func StatData(path string) *StatResult {
    file, err := os.Open(path)
    if err != nil {
        return nil
    }
    defer file.Close()

    stat := StatResult{}
    stat.Init()
    scanner := bufio.NewScanner(file)
    head := []string{}
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
            stat.labels[label] = true
            for i, tk := range tks {
                if i == 0 {
                    continue
                } else {
                    key := head[i][1:]
                    if head[i][0] == '#'{
                        stat.features[key] = true
                        value := tk
                        _, ok := stat.discrete_stat[key]
                        if !ok {
                            stat.discrete_stat[key] = make(map[string]map[string]float64)
                        }
                        _, ok = stat.discrete_stat[key][value]
                        if !ok {
                            stat.discrete_stat[key][value] = make(map[string]float64)
                        }
                        _, ok = stat.discrete_stat[key][value][label]
                        if !ok {
                            stat.discrete_stat[key][value][label] = 1.0
                        } else {
                            stat.discrete_stat[key][value][label] += 1.0
                        }
                    } else if head[i][0] == '*' {
                        stat.features[key] = true
                        _, ok := stat.continuous_stat[key]
                        if !ok {
                            stat.continuous_stat[key] = WeightLabelList{}
                        }
                        value, _ := strconv.ParseFloat(tk, 64)
                        stat.continuous_stat[key] = append(stat.continuous_stat[key], WeightLabel{weight: value, label: label})
                    }
                }
            }
        }
    }
    return &stat
}

var global_stat *StatResult
 
func SingleDiscreteFeatureStat( w http.ResponseWriter,r *http.Request ){
    w.Header().Set("Content-Type", "application/json")
    ret := []interface{}{}
    key := r.FormValue("feature")
    is_100percent := r.FormValue("is_100percent")
    value_label_dis, ok := global_stat.discrete_stat[key]
    
    valuestr := []string{}
    for value, _ := range value_label_dis {
        valuestr = append(valuestr, value)
    }

    if !ok {
        b, _ := json.Marshal(ret)
        fmt.Fprint(w, string(b))
        return
    }

    value_sum := make(map[string]float64)

    for value, label_dis := range value_label_dis {
        value_sum[value] = 0.0
        for _, count := range label_dis {
            value_sum[value] += count
        }
    }

    for label, _ := range global_stat.labels {
        record := make(map[string]interface{})
        record["key"] = key + ": " + label
        values := []map[string]interface{}{}
        for _, value := range valuestr{
            label_dis, _ := value_label_dis[value]
            count, ok := label_dis[label]
            if !ok {
                count = 0.0
            }
            point := make(map[string]interface{})
            point["x"] = value
            if is_100percent == "true" {
                point["y"] = count / value_sum[value]
            } else {
                point["y"] = count
            }
            values = append(values, point)
        }
        record["values"] = values
        ret = append(ret, record)
    }
    b, _ := json.Marshal(ret)
    fmt.Fprint(w, string(b))
}

func SingleContinuousFeatureStat(w http.ResponseWriter, r *http.Request){
    w.Header().Set("Content-Type", "application/json")
    ret := []interface{}{}
    feature := r.FormValue("feature")
    weight_labels, ok := global_stat.continuous_stat[feature]

    if !ok {
        b, _ := json.Marshal(ret)
        fmt.Fprint(w, string(b))
        return
    }
    max_weight := weight_labels.MaxWeight()
}

func FeatureList( w http.ResponseWriter,r *http.Request ){
    w.Header().Set("Content-Type", "application/json")
    ret := []string{}

    for key, _ := range global_stat.discrete_stat {
        ret = append(ret, key)
    }
    b, _ := json.Marshal(ret)
    fmt.Fprint(w, string(b))
}
 
func main(){
    global_stat = StatData("./data/adult/adult.tsv")
    http.Handle("/", http.FileServer(http.Dir(".")))
    http.HandleFunc( "/single_feature",SingleDiscreteFeatureStat)
    http.HandleFunc("/feature_list", FeatureList)
    s := &http.Server{  
        Addr:           ":8080",
        ReadTimeout:    30 * time.Second,
        WriteTimeout:   30 * time.Second,
        MaxHeaderBytes: 1 << 20,
    }
    log.Fatal(s.ListenAndServe())
}