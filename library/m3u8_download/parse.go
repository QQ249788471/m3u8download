package m3u8download

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/grafov/m3u8"
	"io/ioutil"
	"net/http"
)

type M3u8Format struct {

	// m3u8 url
	Url string `json:"url"`

	// http 请求 m3u8 url 得到的内容
	Body string `json:"body"`

	// 解析 m3u8 的错误，或者 http 请求的错误
	Err error `json:"err"`

	// ts url 列表，与 M3u8List 互斥
	TsList []string `json:"ts_list"`

	// 子 m3u8 url 列表，与 TsList 互斥
	M3u8List []M3u8Format `json:"m3u8_list"`
}

func NewM3u8Format(url string) M3u8Format {
	return M3u8Format{
		Url:      url,
		TsList:   make([]string, 0),
		M3u8List: make([]M3u8Format, 0),
	}
}

// 解析 m3u8 文件
func (m *M3u8Format) Parse() bool {

	// 获取 m3u8 url 的内容
	resp, err := http.Get(m.Url)
	if err != nil {
		m.Err = err
		return false
	}
	defer resp.Body.Close()

	// 获取 body
	body, err := ioutil.ReadAll(resp.Body)
	m.Body = fmt.Sprintf("??\n%s\n??", string(body))
	if err != nil {
		m.Err = err
		return false
	}

	// 解析 m3u8
	p, listType, err := m3u8.DecodeFrom(bytes.NewReader(body), true)
	if err != nil {
		m.Err = err
		return false
	}

	// 由于 m3u8 可以有主播放列表，它类似于文件夹，因此这里递归获取所有的 ts 列表
	switch listType {
	case m3u8.MEDIA:
		// 如果是 ts 列表
		mediapl := p.(*m3u8.MediaPlaylist)
		for _, ts := range mediapl.Segments {
			if ts != nil {
				m.TsList = append(m.TsList, ts.URI)
			}
		}
	case m3u8.MASTER:
		// 如果是主播放列表
		masterpl := p.(*m3u8.MasterPlaylist)
		for _, v := range masterpl.Variants {
			if v == nil {
				continue
			}
			mf := NewM3u8Format(v.URI)
			success := mf.Parse()
			m.M3u8List = append(m.M3u8List, mf)

			if !success {
				return false
			}
		}
	}

	return true
}

func (m *M3u8Format) String() string {
	s, _ := json.Marshal(m)
	return string(s)
}
