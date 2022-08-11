package common

import (
	"strings"

	"github.com/aluka-7/utils"
)

// AppVersion App客户端的版本信息封装，包括当前版本号、网络情况、操作平台等信息。
type AppVersion struct {
	device  string // 平台信息
	version string // 版本信息
	network string // 网络信息
}

// NewAppVersion 根据给定的App平台信息、版本信息和网络信息进行解析和存储处理，对当前请求线程都适用。
// @device 客户端设备载体
// @version 客户端版本信息
// @network 客户端请求时的网络情况
func NewAppVersion(device, version, network string) AppVersion {
	return AppVersion{device, version, network}
}

// Device 获取前端客户端的载体，返回值包括如下几种：
//  iOS-Web：基于flutter打包的iOS客户端
//  iOS-Native：为iOS平台开发的原生客户端
//  Android-Web：基于flutter打包的Android客户端
//  Android-Native：为Android平台开发的原生客户端
//  Web：基于Web浏览器的在线版本
//  Windows-Web：基于Electron打包方式的windows客户端
//  Windows-Native：为windows平台开发的原生版本
//  Mac-Web：基于Electron打包方式的Mac平台客户端
//  Mac-Native：为Mac平台开发的原生版本
//  Linux-Web：基于Electron打包方式的Linux平台客户端
//  Linux-Native：为Linux平台开发的原生版本
func (av AppVersion) Device() string {
	return av.device
}

func (av AppVersion) SimpleDevice() string {
	if av.IsClient() {
		return "Native"
	}
	return "Web"
}

// Network 获取当前App的网络标示，返回的值包括如下几种（可能还有其他值）：
//
//  2g：客户端使用的是2G网络
//  3g：客户端使用的是3G网络
//  4g：客户端使用的是4G网络
//  5g：客户端使用的是5g网络
//  wifi：客户端使用的是无线wifi网络
//  unknown：客户端的网络情况未知，如在线版时就无法获知
func (av AppVersion) Network() string {
	return av.network
}

// Version 获取当前客户端的版本信息，如：2.3.1
func (av AppVersion) Version() string {
	return av.version
}

// IsWifi 判断当前App请求是否是wifi网络环境，如果是则返回true，如果不是或前端获取不到数据则返回false。
func (av AppVersion) IsWifi() bool {
	return strings.Contains(av.network, "wifi")
}

// IsClient 判断当前App是否是客户端(包括基于flutter和原生的客户端)，如果是则返回true，否则返回false。
func (av AppVersion) IsClient() bool {
	return strings.Contains(av.device, "Native")
}

// IsWeb 判断当前请求的客户端是否是H5的在线版本，如果是则返回true，否则返回false。
func (av AppVersion) IsWeb() bool {
	return strings.Contains(av.device, "Web")
}

// CompareVersion 使用给定的版本字符串和当前版本信息进行比较，如果当前版本大于给定的版本则返回1，如果当前版本等于给定的版本则返回0，如果当前版本小于给定的版本则返回-1。
func (av AppVersion) CompareVersion(version string) int {
	myVer := strings.Split(av.version, "\\.")
	hisVer := strings.Split(version, "\\.")
	v1 := utils.StrTo(myVer[0]).MustInt64()*1000000 + utils.StrTo(myVer[1]).MustInt64()*1000 + utils.StrTo(myVer[2]).MustInt64()
	v2 := utils.StrTo(hisVer[0]).MustInt64()*1000000 + utils.StrTo(hisVer[1]).MustInt64()*1000 + utils.StrTo(hisVer[2]).MustInt64()
	if v1 == v2 {
		return 0
	} else {
		if v1 > v2 {
			return 0
		} else {
			return -1
		}
	}
}
