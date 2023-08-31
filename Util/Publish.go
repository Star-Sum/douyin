package Util

import (
	"douyin/Log"
	"os"
	"os/exec"
	"strconv"
	"time"
)

// CoverName 封面名称
func CoverName(AuthorName string) string {
	date := time.Now().Unix()
	// 时间戳转换
	dateString := strconv.FormatInt(date, 10)
	return AuthorName + "_" + dateString + "_Cover" + ".png"
}

// VedioName 视频名称
func VedioName(AuthorName string, vedioType string) string {
	date := time.Now().Unix()
	// 时间戳转换
	dateString := strconv.FormatInt(date, 10)
	return AuthorName + "_" + dateString + "_Publish" + "." + vedioType
}

// GetVedioCover 调用ffmpeg进行封面截图，默认截取第一张
func GetVedioCover(fileName string, AuthorName string) string {
	imageName := CoverName(AuthorName)
	vedioPath := "./Resource/Vedio/" + fileName
	imagePath := "./Resource/Image/" + imageName
	_, err := os.Stat(vedioPath)
	if err != nil {
		Log.ErrorLogWithoutPanic("Vedio Open Error!", err)
	}
	args := []string{"-i", vedioPath, "-ss", "0", "-vframes", "1", imagePath}
	cmd := exec.Command("ffmpeg", args...)
	err = cmd.Run()
	if err != nil {
		Log.ErrorLogWithoutPanic("Cover Made Error!", err)
	}
	return imageName
}

// DirToUrl 将文件路径映射为URL
func DirToUrl(filetype string, filename string) string {
	// filetype:文件的种类，有文本（text）、视频（vedio）、照片（image）三种类型
	// filename:完整的文件名，包括文件名和文件种类
	// 文件路径统一采用URL方式存储，文件存储在Resource目录下相应的子目录中
	urlString := "http://192.168.43.174:8080/douyin/"
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
