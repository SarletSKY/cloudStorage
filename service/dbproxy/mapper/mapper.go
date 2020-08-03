package mapper

import (
	"filestore-server-study/service/dbproxy/orm"
	"fmt"
	"reflect"
)

var funcs = map[string]interface{}{
	"/file/AddFileInfoTodb":      orm.AddFileInfoTodb,      // 文件上传完成，保存meta
	"/file/GetFileInfoTodb":      orm.GetFileInfoTodb,      // 从mysql获取文件元信息
	"/file/GetManyFileInfoTodb":  orm.GetManyFileInfoTodb,  // 从mysql批量获取文件元信息
	"/file/UpdateFileLocationdb": orm.UpdateFileLocationdb, // 更新文件的存储地址
	"/file/UpdateFileInfoTodb":   orm.UpdateFileInfoTodb,   // 更新数据到mysql
	"/file/DeleteFileMetaDB":     orm.DeleteFileMetaDB,     // 删除

	"/user/SignInUserdb":    orm.SignInUserdb,    // 登录用户
	"/user/SignUpUserdb":    orm.SignUpUserdb,    // 注册用户
	"/user/RegisterTokendb": orm.RegisterTokendb, // 注册token
	"/user/QueryUserInfodb": orm.QueryUserInfodb, // 查询用户信息
	"/user/GetUserToken":    orm.GetUserToken,    //  查询用户token
	"/user/UserExist":       orm.UserExist,       //  查询用户存不存在

	"/userfile/OnUserFileUploadFinshedDB": orm.OnUserFileUploadFinshedDB, // 上传文件到user_file表
	"/userfile/QueryMantUserFileDB":       orm.QueryMantUserFileDB,       // 批量查询用户文件接口
	"/userfile/UpdateUserFileInfoDB":      orm.UpdateUserFileInfoDB,      // 修改用户文件表的元信息
	"/userfile/IsUserFileUpload":          orm.IsUserFileUpload,          // 判断文件是否存在
	"/userfile/QueryUserFileDB":           orm.QueryUserFileDB,           // 查询单个用户文件元信息
	"/userfile/DeleteUserFileDB":          orm.DeleteUserFileDB,          // 删除用户文件表
	"/userfile/QueryUserFileNameExist":    orm.QueryUserFileNameExist,    // 用户文件名改之前，文件名有没有被使用
}

func FuncCall(name string, params ...interface{}) (result []reflect.Value, err error) {
	if _, ok := funcs[name]; !ok {
		err = fmt.Errorf("函数不存在")
		return
	}
	//fmt.Printf("调用了FuncCall")
	// 通过反射可以动态调用对象的导出方法
	f := reflect.ValueOf(funcs[name])
	//fmt.Println("1111")
	if len(params) != f.Type().NumIn() {
		err = fmt.Errorf("传入参数数量与被调用方法要求数量不一致")
		return
	}
	//fmt.Println("2222")

	// 构造一个Value的slice， 用作Call()方法的传入参数
	in := make([]reflect.Value, len(params))
	for k, param := range params {
		//fmt.Println("33333")
		in[k] = reflect.ValueOf(param)
		//fmt.Println(in)
	}
	//fmt.Println("5555")
	//fmt.Println(name)
	// 执行方法，并将方法结果赋值给result
	result = f.Call(in)
	//fmt.Println("66666")

	return
}
