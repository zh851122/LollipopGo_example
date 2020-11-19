package _go

import (
	"fmt"
	lua "github.com/yuin/gopher-lua"
)

var Gvariable_proto map[string]*STvariable_proto

type STvariable_proto struct {
	Sid   string
	Data1 string
	Data2 string
	Data3 string
	Data4 string
	Data5 string
}

func init() {
	Gvariable_proto = make(map[string]*STvariable_proto)
}

func Initvariable_proto() {
	L := lua.NewState()
	variable_prototmp := make(map[string]*STvariable_proto)
	defer L.Close()
	err := L.DoFile("./lua/variable_proto.lua")
	if err != nil {
		fmt.Println("err", err.Error())
		panic(err)
	}
	variable_proto := L.GetGlobal("variable_proto")
	if tbl, ok := variable_proto.(*lua.LTable); ok {
		fmt.Println(L.ObjLen(tbl))
		tbl.ForEach(func(value lua.LValue, value2 lua.LValue) {
			data := new(STvariable_proto)
			data.Sid = value.String()
			if tbl2, ok2 := value2.(*lua.LTable); ok2 {
				fmt.Println(L.ObjLen(tbl2))
				tbl2.ForEach(func(value3 lua.LValue, value4 lua.LValue) {
					if "data1" == value3.String() {
						data.Data1 = value4.String()
					}
					if "data2" == value3.String() {
						data.Data2 = value4.String()
					}
					if "data3" == value3.String() {
						data.Data3 = value4.String()
					}
					if "data4" == value3.String() {
						data.Data4 = value4.String()
					}
					if "data5" == value3.String() {
						data.Data5 = value4.String()
					}
				})
			}
			variable_prototmp[data.Sid] = data
		})
	}
	Gvariable_proto = variable_prototmp
}
