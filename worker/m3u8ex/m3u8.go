package m3u8ex

import (
	"errors"
	"fmt"
	m3u8download "github.com/lyd2/goff-net-task/library/m3u8_download"
	"net/url"
	"strings"
)

type TsInfo struct {

	// 此 ts 的 uri，它是一个绝对地址
	Uri string

	// 此 ts 的保存文件名
	Filename string

	// 是否解析成功
	Status bool

	// 解析此 ts uri 的错误
	Error string
}

func NewTsInfo(uri, filename string, err error) TsInfo {

	status := true
	errinfo := ""
	if err != nil {
		status = false
		errinfo = err.Error()
	}

	return TsInfo{
		Uri:      uri,
		Filename: filename,
		Status:   status,
		Error:    errinfo,
	}

}

// 下载并解析 m3u8 文件
func DownloadAndParseM3u8(url string) (m3u8f m3u8download.M3u8Format, err error) {
	m3u8f = m3u8download.NewM3u8Format(url)
	if !m3u8f.Parse() {
		err = errors.New("m3u8 file parse error")
	}
	return
}

// 搜索 ts 列表
func SearchTs(format *m3u8download.M3u8Format) []TsInfo {

	tsInfoList := []TsInfo{}

	// ts 列表
	if len(format.M3u8List) == 0 {
		for _, v := range format.TsList {
			tsInfoList = append(tsInfoList, NewTsInfo(buildTsName(format.Url, v)))
		}
		return tsInfoList
	}

	// m3u8 master
	for _, v := range format.M3u8List {
		tsInfoList = append(tsInfoList, SearchTs(&v)...)
	}
	return tsInfoList
}

// 生成 ts 的 uri
// 因为有的 m3u8 里的 ts uri 是不带有域名的，这里会判断并拼接上域名
func buildTsName(m3u8, ts string) (uri, filename string, err error) {

	tsUrl, err := url.Parse(ts)
	if err != nil {
		return ts, "", err
	}

	// 是绝对 url 路径
	if tsUrl.IsAbs() {
		// url 编码解析
		uri, err = url.QueryUnescape(tsUrl.String())
		if err != nil {
			return tsUrl.String(), "", err
		}
		return uri, strings.Replace(tsUrl.Path, "/", "_", -1), nil
	}

	// 是相对路径
	m3u8Url, err := url.Parse(m3u8)
	if err != nil {
		return m3u8, "", err
	}

	m3u8AndTsUrl, err := m3u8Url.Parse(ts)
	if err != nil {
		return fmt.Sprintf("%s--%s", m3u8, ts), "", err
	}

	// url 编码解析
	uri, err = url.QueryUnescape(m3u8AndTsUrl.String())
	if err != nil {
		return m3u8AndTsUrl.String(), "", err
	}
	return uri, strings.Replace(m3u8AndTsUrl.Path, "/", "_", -1), nil

}
