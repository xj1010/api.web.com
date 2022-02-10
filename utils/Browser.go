/**
1.可设置代理
2.可设置 cookie
3.自动保存并应用响应的 cookie
4.自动为重新向的请求添加 cookie
*/
package utils

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/http/cookiejar"
	"net/textproto"
	"net/url"
	"strconv"
	_ "strconv"
	"strings"
	"time"
)

type Browser struct {
	cookies []*http.Cookie
	headers map[string]string
	client *http.Client
}

//初始化
func NewBrowser() *Browser {
	hc := &Browser{}

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	hc.headers = make(map[string]string,0)

	hc.client = &http.Client{Timeout : time.Duration(5)* time.Second, Transport:tr}
	headers := map[string]string {
		"Accept-Language" : "zh-CN",
		"User-Agent" : "Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/75.0.3770.100 Safari/537.36",
		"Referer" : "https://www.lagou.com",
		"Content-Type" : "application/x-www-form-urlencoded",
		"Host" : "https://www.baidu.com",
		"Origin": "https://www.lagou.com",
		"Connection": "keep-alive",

	}
	hc.AddHeader(headers)

	rand.Seed(time.Now().UnixNano())
	num := rand.Intn(100)
	if num > 50 {
		//hc.SetProxyUrl("socks5h://localhost:1080")
	}
	//hc.SetProxyUrl("socks5h://localhost:1080")
	hc.client.Jar, _ = cookiejar.New(nil)
	//为所有重定向的请求增加cookie
	hc.client.CheckRedirect = func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse
	}

	return hc
}

//设置代理地址
func (self *Browser) SetProxyUrl(proxyUrl string)  {
	proxy := func(_ *http.Request) (*url.URL, error) {
		return url.Parse(proxyUrl)
	}
	transport := &http.Transport{Proxy:proxy}
	self.client.Transport = transport
}

//获取当前所有的cookie
func (self *Browser) GetCookie() ([]*http.Cookie) {
	return self.cookies
}

//设置请求cookie
func (self *Browser) AddCookie(cookies map[string]string)  {
	for k, v := range cookies {
		cookie := &http.Cookie{Name : k, Value : v}
		self.cookies = append(self.cookies, cookie)
	}
}

//设置请求header
func (self *Browser) AddHeader(headers map[string]string)  {
	for k, v := range headers {
		self.headers[k] = v
	}
}

//发送Get请求
func (self *Browser) Get(requestUrl string) ([]byte, int) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println(r)
		}
	}()

	request, _ := http.NewRequest("GET", requestUrl, nil)
	self.setRequestCookie(request)
	self.setRequestHeader(request)
	response,_ := self.client.Do(request)
	defer response.Body.Close()

	data, _ := ioutil.ReadAll(response.Body)
	return data, response.StatusCode
}

//发送Get请求
func (self *Browser) GetRedirect(requestUrl string) ([]byte, int, map[string]string) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println(r)
		}
	}()

	redirectMap := make(map[string]string,3)

	request, _ := http.NewRequest("GET", requestUrl, nil)
	self.setRequestCookie(request)
	self.setRequestHeader(request)
	response,_ := self.client.Do(request)
	defer response.Body.Close()

	Redirect:
	data, _ := ioutil.ReadAll(response.Body)
	if response.StatusCode == 302 || response.StatusCode == 301 {
		location := response.Header.Get("location")
		if len(location) > 0 {
			domain, _ := url.Parse(location)
			self.setCookieWhenRedirect(self.client, domain, response)

			request, _ := http.NewRequest("GET", location, nil)
			response,_ = self.client.Do(request)
			for _, v := range self.client.Jar.Cookies(request.URL){
				redirectMap[v.Name] = v.Value
			}
			goto Redirect
		}

	}

	return data, response.StatusCode, redirectMap
}

func (self *Browser) setCookieWhenRedirect(client *http.Client, domain *url.URL, resp *http.Response) {
	cookieInHeader := resp.Header[textproto.CanonicalMIMEHeaderKey("Set-Cookie")]
	if len(cookieInHeader) == 0 {
		return
	}

	cookieList := make([]*http.Cookie, 0)
	for _, cookieStr := range cookieInHeader {
		cookie := &http.Cookie{}
		parts := strings.Split(cookieStr, ";")
		for _, item := range parts {
			if item == "" {
				continue
			}
			item = strings.Trim(item, " ")
			if item == "HttpOnly" {
				cookie.HttpOnly = true
				continue
			}

			subItems := strings.Split(item, "=")
			if len(subItems) != 2 {
				continue
			}
			if subItems[0] == "Path" || subItems[0] == "path" {
				cookie.Path = subItems[1]
				continue
			}
			if subItems[0] == "Max-Age" {
				maxAge, _ := strconv.Atoi(subItems[1])
				cookie.MaxAge = maxAge
				continue
			}
			if subItems[0] == "Domain"  || subItems[0] == "domain" {
				cookie.Domain = subItems[1]
				continue
			}
			cookie.Name 	= subItems[0]
			cookie.Value 	= subItems[1]
		}
		cookieList = append(cookieList, cookie)
	}

	client.Jar.SetCookies(domain, cookieList)
}

//发送Post请求
func (self *Browser) Post(requestUrl string, params map[string]string) ([]byte) {
	defer func() {
		if r := recover(); r != nil {
		}
	}()
	postData := self.encodeParams(params)
	request, _ := http.NewRequest("POST", requestUrl, strings.NewReader(postData))
	self.setRequestCookie(request)
	self.setRequestHeader(request)

	response, _ := self.client.Do(request)
	defer response.Body.Close()

	//保存响应的 cookie
	self.SetResponseCookie(response.Cookies())
	data, _ := ioutil.ReadAll(response.Body)
	return data
}

func (self *Browser) PostJson(requestUrl string, params map[string]interface{}) []byte {
	defer func() {
		if r := recover(); r != nil {
		}
	}()
	str, err := json.Marshal(params)
	if err != nil {
		return nil
	}

	request, _ := http.NewRequest("POST", requestUrl, strings.NewReader(string(str)))
	self.headers["Content-Type"] = "application/json"
	self.setRequestCookie(request)
	self.setRequestHeader(request)

	response, _ := self.client.Do(request)
	defer response.Body.Close()

	//保存响应的 cookie
	self.SetResponseCookie(response.Cookies())
	data, _ := ioutil.ReadAll(response.Body)
	return data
}

//为请求设置 cookie
func (self *Browser) setRequestCookie(request *http.Request)  {
	for _,v := range self.cookies{
		request.AddCookie(v)
	}
}

func (self *Browser) SetResponseCookie(cookies []*http.Cookie)  {
	self.cookies = cookies
}

func (self *Browser) GetResponseCookie() []*http.Cookie {
	return self.cookies
}

func (self *Browser) setRequestHeader(request *http.Request)  {
	for k , v := range self.headers {
		request.Header.Set(k, v)
	}
}

//参数 encode
func (self *Browser) encodeParams(params map[string]string) string {
	paramsData := url.Values{}
	for k,v := range params {
		paramsData.Set(k,v)
	}
	return paramsData.Encode()
}

