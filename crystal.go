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
    "io/ioutil"
    "sort"
    "crystal/service"
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

func (self *WeightLabelList) MinWeight() float64 {
    ret := (*self)[0].weight
    for _, wl := range *self {
        if ret > wl.weight {
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
                    key := head[i]
                    if key[0] == '['{
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
                    } else {
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

func StatAllData() map[string]*StatResult {
    ret := make(map[string]*StatResult)
    files, _ := ioutil.ReadDir("./data/")
    for _, f := range files {
        if f.IsDir(){
            stat := StatData("./data/" + f.Name() + "/data.tsv")
            ret[f.Name()] = stat
        }
    }
    return ret
}

var global_stat map[string]*StatResult

func SingleContinuousFeatureStat(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    params := service.GetParameters(w, r)
    is_100percent := r.FormValue("is_100percent")
    weight_labels, _ := global_stat[params.Dataset].continuous_stat[params.Feature]

    max_weight := weight_labels.MaxWeight()
    min_weight := weight_labels.MinWeight()

    bin_stat := make(map[int]map[string]float64)
    bin_sum := make(map[int]float64)
    for _, wl := range weight_labels {
        bin := int(20.0 * (wl.weight - min_weight) / (max_weight - min_weight))
        _, ok := bin_stat[bin]
        if !ok {
            bin_stat[bin] = make(map[string]float64)
            bin_sum[bin] = 0.0
        }
        bin_sum[bin] += 1.0
        _, ok = bin_stat[bin][wl.label]
        if !ok {
            bin_stat[bin][wl.label] = 1.0
        } else {
            bin_stat[bin][wl.label] += 1.0
        }
    }
    ret := []interface{}{}

    bins := []int{}
    for bin, _ := range bin_sum {
        bins = append(bins, bin)
    }
    sort.Ints(bins)
    for label, _ := range global_stat[params.Dataset].labels {
        record := make(map[string]interface{})
        record["key"] = params.Feature + ": " + label
        values := []map[string]interface{}{}
        for _, bin := range bins {
            label_dis, _ := bin_stat[bin]
            count, ok := label_dis[label]
            if !ok {
                count = 0.0
            }
            point := make(map[string]interface{})
            point["x"] = min_weight + float64(bin) * (max_weight - min_weight) / 20.0
            if is_100percent == "true" {
                point["y"] = count / bin_sum[bin]
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

func SingleFeatureStat(w http.ResponseWriter, r *http.Request){
    feature := r.FormValue("feature")
    if feature[0] == '[' {
        SingleDiscreteFeatureStat(w, r)
    } else {
        SingleContinuousFeatureStat(w, r)
    }
}

type StringDoublePair struct {
    Key string
    Value float64
}

type StringDoublePairList struct {
    data []*StringDoublePair
}

func (self *StringDoublePairList) Less(i, j int) bool{
    return self.data[i].Value > self.data[j].Value
}

func (self *StringDoublePairList) Len() int{
    return len(self.data)
}

func (self *StringDoublePairList) Swap(i, j int) {
    self.data[i], self.data[j] = self.data[j], self.data[i]
}
 
func SingleDiscreteFeatureStat(w http.ResponseWriter,r *http.Request){
    w.Header().Set("Content-Type", "application/json")
    ret := []interface{}{}
    feature := r.FormValue("feature")
    dataset := r.FormValue("dataset")
    is_100percent := r.FormValue("is_100percent")
    value_label_dis, ok := global_stat[dataset].discrete_stat[feature]
    
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

    value_sum_array := StringDoublePairList{}
    for value, sum := range value_sum {
        value_sum_array.data = append(value_sum_array.data, &(StringDoublePair{Key: value, Value: sum}))
    }
    sort.Sort(&value_sum_array)

    for label, _ := range global_stat[dataset].labels {
        record := make(map[string]interface{})
        record["key"] = feature + ": " + label
        values := []map[string]interface{}{}
        for i, vs := range value_sum_array.data{
            if i > 64 {
                break
            }
            value := vs.Key
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

func FeatureList( w http.ResponseWriter,r *http.Request ){
    w.Header().Set("Content-Type", "application/json")
    dataset := r.FormValue("dataset")
    ret := []string{}

    for key, _ := range global_stat[dataset].discrete_stat {
        ret = append(ret, key)
    }
    for key, _ := range global_stat[dataset].continuous_stat {
        ret = append(ret, key)
    }
    b, _ := json.Marshal(ret)
    fmt.Fprint(w, string(b))
}

func DataSetList(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    files, _ := ioutil.ReadDir("./data/")
    ret := []string{}
    for _, f := range files {
        if f.IsDir(){
            ret = append(ret, f.Name())
        }
    }
    b, _ := json.Marshal(ret)
    fmt.Fprint(w, string(b))
}

func main(){
    global_stat = StatAllData()
    http.Handle("/", http.FileServer(http.Dir(".")))
    http.HandleFunc( "/single_feature",SingleFeatureStat)
    http.HandleFunc("/feature_list", FeatureList)
    http.HandleFunc("/dataset_list", DataSetList)
    s := &http.Server{  
        Addr:           ":2014",
        ReadTimeout:    30 * time.Second,
        WriteTimeout:   30 * time.Second,
        MaxHeaderBytes: 1 << 20,
    }
    log.Fatal(s.ListenAndServe())
}