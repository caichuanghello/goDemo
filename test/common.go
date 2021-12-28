package main

import "os"

//判断文件或者目录是否存在
func FileOrDirIsExit(path string) (b bool,err error){
	_,err=os.Stat(path)
	re:=os.IsNotExist(err)
	if os.IsNotExist(err) {
		return false,nil
	}
	return !re ,err
}