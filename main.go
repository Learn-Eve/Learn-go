package main

import (
	"fast-learn/cmd"
	"fast-learn/utils"
	"fmt"
)

// @title Go-web开发记录
// @version 0.0.1
// @description 从零开始学go web实战
func main() {
	defer cmd.Clean()
	cmd.Start()
	// jwt相关
	token, _ := utils.GenerateToken(1, "zs")
	fmt.Println(token)

	iJwtCustClaims, err := utils.ParseToken(token + "123")
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println(iJwtCustClaims)
}
