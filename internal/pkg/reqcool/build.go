package reqcool

import (
	"net/url"

	"github.com/imroc/req/v3"
)

type Client struct {
	*req.Client
	authCookie bool // 是否自动设置 Cookie
}

func NewGoofishApi(baseUrl string, authSetCookie bool) *Client {
	client := req.C().
		SetBaseURL(baseUrl).
		EnableInsecureSkipVerify().
		SetCookieJar(nil) // 启用默认 CookieJar，自动管理 Cookie
	return &Client{
		Client:     client,
		authCookie: authSetCookie,
	}
}

func (c *Client) Get(api string, params url.Values, result any) error {
	r := c.R()
	r.QueryParams = params

	return c.do(r, "GET", api, result)
}

func (i *Client) Post(api string, params url.Values, body map[string]any, cookie string, result any) error {
	r := i.R()
	r.QueryParams = params
	r.SetFormDataAnyType(body)

	if cookie != "" {
		r.SetHeader("Cookie", cookie)
	}

	return i.do(r, "POST", api, result)
}

func (c *Client) do(r *req.Request, method, path string, result any) error {
	resp, err := r.Send(method, path)
	if err != nil {
		return err
	}

	if resp.Cookies() != nil && c.authCookie {
		// 更新 CookieJar 中的 Cookies
		c.SetCommonCookies(resp.Cookies()...)
	}

	return resp.Into(result)
}
