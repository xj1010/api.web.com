package utils

import (
	"github.com/axgle/mahonia"
	"os"
)

func InArray(target string, strArr []string) bool {
	for _, element := range strArr{
		if target == element{
			return true
		}
	}
	return false
}

func ConvertToString(src string, srcCode string, tagCode string) string {
	srcCoder := mahonia.NewDecoder(srcCode)
	srcResult := srcCoder.ConvertString(src)
	tagCoder := mahonia.NewDecoder(tagCode)
	_, cdata, _ := tagCoder.Translate([]byte(srcResult), true)
	result := string(cdata)
	return result

}

/**
判断文件是否存在
*/
func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}

	if os.IsNotExist(err) {
		return false, nil
	}

	return false, err
}

/**
创建文件夹
*/
func CreateDir(dirName string) bool {
	err := os.Mkdir(dirName,755)
	if err != nil{
		return false
	}

	return true
}


