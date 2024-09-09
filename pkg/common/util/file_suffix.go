package util

import (
	"bytes"                         // 引入bytes包，用于字节操作
	"chat-room/pkg/common/constant" // 引入常量包，定义了消息内容类型
	"encoding/hex"                  // 引入hex包，用于十六进制编码解码
	"strconv"                       // 引入strconv包，用于字符串和数字的转换
	"strings"                       // 引入strings包，用于字符串操作
	"sync"                          // 引入sync包，用于并发安全的操作

	"github.com/wxnacy/wgo/arrays" // 引入第三方数组操作库
)

// fileTypeMap 是一个并发安全的映射，用于存储文件头标识和文件类型的对应关系
var fileTypeMap sync.Map

// init 函数在包初始化时执行，初始化 fileTypeMap，存储常见文件类型的文件头标识
func init() {
	// 添加常见文件类型及其对应的文件头标识
	fileTypeMap.Store("ffd8ffe000104a464946", "jpg")  // JPEG (jpg)
	fileTypeMap.Store("89504e470d0a1a0a0000", "png")  // PNG (png)
	fileTypeMap.Store("47494638396126026f01", "gif")  // GIF (gif)
	fileTypeMap.Store("49492a00227105008037", "tif")  // TIFF (tif)
	fileTypeMap.Store("424d228c010000000000", "bmp")  // 16色位图(bmp)
	fileTypeMap.Store("424d8240090000000000", "bmp")  // 24位位图(bmp)
	fileTypeMap.Store("424d8e1b030000000000", "bmp")  // 256色位图(bmp)
	fileTypeMap.Store("41433130313500000000", "dwg")  // CAD (dwg)
	fileTypeMap.Store("3c21444f435459504520", "html") // HTML (html)
	fileTypeMap.Store("3c68746d6c3e0", "html")        // HTML (html)
	fileTypeMap.Store("3c21646f637479706520", "htm")  // HTM (htm)
	fileTypeMap.Store("48544d4c207b0d0a0942", "css")  // CSS (css)
	fileTypeMap.Store("696b2e71623d696b2e71", "js")   // JavaScript (js)
	fileTypeMap.Store("7b5c727466315c616e73", "rtf")  // Rich Text Format (rtf)
	fileTypeMap.Store("38425053000100000000", "psd")  // Photoshop (psd)
	fileTypeMap.Store("46726f6d3a203d3f6762", "eml")  // Email [Outlook Express 6] (eml)
	fileTypeMap.Store("d0cf11e0a1b11ae10000", "vsd")  // Visio 绘图 (vsd)
	fileTypeMap.Store("5374616E64617264204A", "mdb")  // MS Access (mdb)
	fileTypeMap.Store("252150532D41646F6265", "ps")   // PostScript (ps)
	// 省略部分其他文件类型的初始化...
}

// bytesToHexString 函数将字节数组转换为十六进制字符串表示
func bytesToHexString(src []byte) string {
	res := bytes.Buffer{} // 创建一个字节缓冲区
	if src == nil || len(src) <= 0 {
		return "" // 如果输入为空，返回空字符串
	}
	temp := make([]byte, 0)
	for _, v := range src {
		sub := v & 0xFF
		hv := hex.EncodeToString(append(temp, sub)) // 将字节转换为十六进制字符串
		if len(hv) < 2 {
			res.WriteString(strconv.FormatInt(int64(0), 10)) // 补充不足两位的十六进制表示
		}
		res.WriteString(hv) // 将转换后的十六进制字符串写入缓冲区
	}
	return res.String() // 返回最终的十六进制字符串
}

// GetFileType 函数根据文件的前几个字节来判断文件类型
// fSrc: 文件字节流（只需要前几个字节即可判断）
func GetFileType(fSrc []byte) string {
	var fileType string
	fileCode := bytesToHexString(fSrc) // 将文件字节流转换为十六进制字符串

	// 遍历 fileTypeMap，判断文件类型
	fileTypeMap.Range(func(key, value interface{}) bool {
		k := key.(string)
		v := value.(string)
		if strings.HasPrefix(fileCode, strings.ToLower(k)) ||
			strings.HasPrefix(k, strings.ToLower(fileCode)) {
			fileType = v
			return false // 找到匹配的文件类型后停止遍历
		}
		return true // 继续遍历
	})
	return fileType
}

// GetContentTypeBySuffix 函数根据文件后缀名判断文件内容类型
func GetContentTypeBySuffix(suffix string) int32 {
	imgList := []string{"jpeg", "jpg", "png", "gif", "tif", "bmp", "dwg"} // 图片格式列表
	exists := arrays.Contains(imgList, suffix)                            // 判断后缀名是否在图片格式列表中
	if exists >= 0 {
		return constant.IMAGE // 返回图片类型常量
	}

	audioList := []string{"mp3", "wma", "wav", "mid", "ape", "flac"} // 音频格式列表
	existAudio := arrays.Contains(audioList, suffix)                 // 判断后缀名是否在音频格式列表中
	if existAudio >= 0 {
		return constant.AUDIO // 返回音频类型常量
	}

	videoList := []string{"rmvb", "flv", "mp4", "mpg", "mpeg", "avi", "rm", "mov", "wmv", "webm"} // 视频格式列表
	existVideo := arrays.Contains(videoList, suffix)                                              // 判断后缀名是否在视频格式列表中
	if existVideo >= 0 {
		return constant.VIDEO // 返回视频类型常量
	}
	return constant.FILE // 返回文件类型常量
}
