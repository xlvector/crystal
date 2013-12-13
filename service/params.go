package service

import(
    "net/http"
)

type Parameters struct {
	Feature, Dataset string
}

func GetParameters(w http.ResponseWriter, r *http.Request) *Parameters {
	ret := Parameters{}
	ret.Feature = r.FormValue("feature")
	ret.Dataset = r.FormValue("dataset")
	return &ret
}