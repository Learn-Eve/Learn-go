package main

import "fast-learn/cmd"

// @title Go-web开发记录
// @version 0.0.1
// @description 从零开始学go web实战
func main() {
	defer cmd.Clean()
	cmd.Start()
}
