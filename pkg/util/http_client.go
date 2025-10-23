package util

import (
	"context"
	"fmt"
	"time"

	"github.com/go-resty/resty/v2"
)

// Client 是封装结构体
type HttpClient struct {
	restyClient *resty.Client
	baseURL     string
}

// NewHttpClient 使用 Resty 创建 Client
func NewHttpClient(baseURL string, timeout time.Duration) *HttpClient {
	r := resty.New().
		SetBaseURL(baseURL).
		SetTimeout(timeout).
		SetHeader("Accept", "*/*")

	// 可以添加默认中间件，比如请求/响应日志
	r.OnBeforeRequest(func(c *resty.Client, req *resty.Request) error {
		return nil
	})
	r.OnAfterResponse(func(c *resty.Client, resp *resty.Response) error {
		// 响应后日志、错误处理
		return nil
	})

	// 重试机制示例
	r.SetRetryCount(3).
		SetRetryWaitTime(2 * time.Second).
		SetRetryMaxWaitTime(10 * time.Second)

	return &HttpClient{
		restyClient: r,
		baseURL:     baseURL,
	}
}

// SetHeader 设置一个默认 Header
func (c *HttpClient) SetHeader(key, value string) {
	c.restyClient.SetHeader(key, value)
}

// Get 执行 GET 请求并解析 JSON 到 result
func (c *HttpClient) Get(ctx context.Context, path string, query map[string]string, result interface{}) error {
	resp, err := c.restyClient.R().
		SetContext(ctx).
		SetQueryParams(query).
		SetResult(result).
		Get(path)
	if err != nil {
		return err
	}
	if resp.IsError() {
		return fmt.Errorf("error: request status=%d\tresponse=%s",
			resp.StatusCode(), resp.String())
	}
	return nil
}

func (c *HttpClient) Post(ctx context.Context, path string, body interface{}, result interface{}) error {
	resp, err := c.restyClient.R().
		SetContext(ctx).
		SetBody(body).
		SetResult(result).
		Post(path)
	if err != nil {
		return err
	}
	if resp.IsError() {
		return fmt.Errorf("error: status=%d\tresponse=%s",
			resp.StatusCode(), resp.String())
	}
	return nil
}
