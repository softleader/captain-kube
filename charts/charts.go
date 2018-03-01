package charts

import (
	"io/ioutil"
	"gopkg.in/yaml.v2"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

/**
parse yaml 檔案放入 map 回傳
 */
func Load(filename string) map[interface{}]interface{} {

	dat, err := ioutil.ReadFile(filename)
	check(err)

	m := make(map[interface{}]interface{})
	yaml.Unmarshal(dat, &m)

	return m

}

type KeywordValues struct {
	MappingValues []string
}

func FindKeywordFromMap(m map[interface{}]interface{}, keyword string, kv *KeywordValues) {
	for key, value := range m {
		switch value.(type) {
		case map[interface{}]interface{}:
			FindKeywordFromMap(value.(map[interface{}]interface{}), keyword, kv)

		default:
			if key == keyword {
				//log.Println("match => Key:", key, "Value:", value)
				kv.MappingValues = append(kv.MappingValues, value.(string))
			}
		}
	}
}
