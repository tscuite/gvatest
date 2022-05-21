package main

import (
	"github.com/tscuite/gvatest/gvalogin"
)

const url string = "https://demo.gin-vue-admin.com"

func main() {
	gvalogin.Gvalogin(url,Map())
}

//登陆后访问地址,后期改为表格获取或数据库，或自动化接口文档，
func Map() map[string]string {
	content := `{"page": 1, "pageSize": 999}`
	countryCapitalMap := make(map[string]string)
	countryCapitalMap["https://demo.gin-vue-admin.com/api/authority/getAuthorityList"] = content
	countryCapitalMap["https://demo.gin-vue-admin.com/api/menu/getMenu"] = content
	return countryCapitalMap
}