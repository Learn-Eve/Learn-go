package test

import (
	"encoding/json"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

var stRootDir02 string
var stSeperator02 string
var iRootNode Node

const stJsonFileName02 = "dir.json"

type Node struct {
	Text     string `json:"text"`
	Children []Node `json:"children"`
}

func loadJson02() {
	stSeperator02 = string(filepath.Separator)
	stWorkDir, _ := os.Getwd()
	stRootDir02 = stWorkDir[:strings.LastIndex(stWorkDir, stSeperator02)]

	gnJsonFileBytes, _ := os.ReadFile(stWorkDir + stSeperator02 + stJsonFileName02)
	err := json.Unmarshal(gnJsonFileBytes, &iRootNode)
	if err != nil {
		panic("Load Json Error:" + err.Error())
	}
}

func parseNode(iNode Node, stParentDir string) {
	if iNode.Text != "" {
		createDir02(iNode, stParentDir)
	}

	if stParentDir != "" {
		stParentDir = stParentDir + stSeperator02
	}

	if iNode.Text != "" {
		stParentDir = stParentDir + iNode.Text
	}

	for _, v := range iNode.Children {
		parseNode(v, stParentDir)
	}
}

func createDir02(iNode Node, stParentDir string) {
	stDirPath := stRootDir02 + stSeperator02
	if stParentDir != "" {
		stDirPath = stDirPath + stParentDir + stSeperator02
	}
	stDirPath = stDirPath + iNode.Text
	//fmt.Println(stDirPath)

	err := os.MkdirAll(stDirPath, fs.ModePerm)
	if err != nil {
		panic("Create Dir Error:" + err.Error())
	}
}

func TestGenerateDir02(t *testing.T) {
	loadJson02()
	parseNode(iRootNode, "")
}
