package meet

import (
	"encoding/json"
	"fmt"
	"io"
	"lyworder/global"
	"net/http"
	"strings"
)

type Robot struct {
	ToUserName []string `json:"toUserName,omitempty"`
	ToUsers    []string `json:"toUsers,omitempty"`
	ToUser     string   `json:"touser,omitempty"`
	Content    string   `json:"content"`
}

func Send(r *Robot) error {

	jsonBody, err := json.Marshal(r)
	if err != nil {
		return fmt.Errorf("JSON序列化失败: %v", err)
	}

	// 创建POST请求
	req, err := http.NewRequest("POST", fmt.Sprintf("%s%s", global.MeetSetting.Host, global.MeetSetting.RoBotKey), strings.NewReader(string(jsonBody)))
	if err != nil {
		return fmt.Errorf("创建请求失败: %v", err)
	}

	// 设置请求头
	req.Header.Set("Content-Type", "application/json")

	// 发送请求
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("发送请求失败: %v", err)
	}
	defer resp.Body.Close()

	// 读取响应
	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("请求失败: 状态码 %d, 响应: %s", resp.StatusCode, string(body))
	}

	return nil
}
