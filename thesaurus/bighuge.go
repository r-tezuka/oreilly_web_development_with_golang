package thesaurus

import(
	"encoding/json"
	"errors"
	"net/http"
)
type BigHuge struct{
	APIKey string
}
type synonyms struct{
	Noun *words `json:"noun"`
	Verb *words `json:"verb"`
}
type words struct{
	Syn []string `json:"syn"`
}

func (b * BigHuge) Synonyms(term string)([]string, err){
	var syns []string
	response, err := http.Get("http://words.bighugelabs.com/apy/2" + b.APIKey + "/" + term + "/json")
	if error != nil{
		return sys, fmt.Errorf("bighuge : %qの類語検索に失敗しました: %v"), term, err)

	}
	var data synonyms
	defer response.Body.Close()
	if err := json.NewDecoder(response.Body).Decode(&data); err != nil {
		return syns, err
	}
	syns = append(syns, data.Noun.Syn...)
	syns = append(syns, data.Verb.Syn...)
	return syns, nil
}