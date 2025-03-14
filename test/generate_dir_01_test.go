package test

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

var stRootDir string
var stSeparator string
var iJsonData map[string]any

const stJsonFileName = "dir.json"

func loadJson() {
	stSeparator = string(filepath.Separator)
	stWorkDir, _ := os.Getwd()
	stRootDir = stWorkDir[:strings.LastIndex(stWorkDir, stSeparator)]

	//fmt.Println(stWorkDir)
	//fmt.Println(stRootDir)

	gnJsonBytes, _ := os.ReadFile(stWorkDir + stSeparator + stJsonFileName)

	err := json.Unmarshal(gnJsonBytes, &iJsonData)

	if err != nil {
		panic("Load Json Error:" + err.Error())
	}
}

func parseMap(mapData map[string]any, stParantDir string) {
	for key, value := range mapData {
		switch value.(type) {
		case string:
			{
				path, _ := value.(string)
				if path == "" {
					continue
				}

				if stParantDir != "" {
					path = stParantDir + stSeparator + path
					if key == "text" {
						stParantDir = path
					}
				} else {
					stParantDir = path
				}

				createDir(path)
			}
		case []any:
			{
				parseArray(value.([]any), stParantDir)
			}
		}

	}
}

func parseArray(giJsonData []any, stParentDir string) {
	for _, v := range giJsonData {
		mapV, _ := v.(map[string]any)
		parseMap(mapV, stParentDir)
	}
}

func createDir(path string) {
	if path == "" {
		return
	}
	fmt.Println(path)

	err := os.MkdirAll(stRootDir+stSeparator+path, fs.ModePerm)
	if err != nil {
		panic("Create Dir Error:" + err.Error())
	}
}

func TestGenerateDir01(t *testing.T) {
	loadJson()
	parseMap(iJsonData, "")
}
