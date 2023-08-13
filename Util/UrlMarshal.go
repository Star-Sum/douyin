package Util

// DirToUrl 将文件路径映射为URL
func DirToUrl(filetype string, filename string) string {
	// filetype:文件的种类，有文本（text）、视频（vedio）、照片（image）三种类型
	// filename:完整的文件名，包括文件名和文件种类
	// 文件路径统一采用URL方式存储，文件存储在Resource目录下相应的子目录中
	urlString := "http://localhost:8080/douyin/"
	switch filetype {
	case "text":
		urlString += "text/" + filename
	case "vedio":
		urlString += "vedio/" + filename
	case "image":
		urlString += "image/" + filename
	}
	return urlString
}
