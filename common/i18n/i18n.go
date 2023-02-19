package i18n

import (
	"AynaLivePlayer/common/config"
	"AynaLivePlayer/common/util"
	"encoding/json"
	"os"
)

const FILENAME = "translation.json"

type Translation struct {
	Languages []string
	Messages  map[string]map[string]string
}

func (t *Translation) HasLanguage(lang string) bool {
	for _, l := range t.Languages {
		if l == lang {
			return true
		}
	}
	return false
}

var TranslationMap Translation
var CurrentLanguage string

func init() {
	TranslationMap = Translation{make([]string, 0), make(map[string]map[string]string)}
	file, err := os.ReadFile(config.GetAssetPath(FILENAME))
	if err == nil {
		_ = json.Unmarshal([]byte(file), &TranslationMap)
	}
	LoadLanguage(config.General.Language)
}

func LoadLanguage(lang string) {
	CurrentLanguage = lang
	if TranslationMap.HasLanguage(lang) {
		return
	}
	TranslationMap.Languages = append(TranslationMap.Languages, lang)
	for id, m := range TranslationMap.Messages {
		m[lang] = id
	}
}

func T(id string) string {
	if x, ok := TranslationMap.Messages[id]; ok {
		return x[CurrentLanguage]
	}
	TranslationMap.Messages[id] = make(map[string]string)
	for _, l := range TranslationMap.Languages {
		TranslationMap.Messages[id][l] = id
	}
	return id
}

func SaveTranslation() {
	content, _ := util.MarshalIndentUnescape(TranslationMap, "", "  ")
	_ = os.WriteFile(config.GetAssetPath(FILENAME), []byte(content), 0666)
}
