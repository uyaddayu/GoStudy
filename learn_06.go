package main

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"path/filepath"
	"strings"
)

const configFile = "./conf/"

var allConfig map[string]*Config
var httpServer, allServer *http.Server

type Config struct {
	HttpToHttps       bool     `json:"http_to_https"`       // 是否将HTTP强制跳转到HTTPS
	ListenHttp        bool     `json:"listen_http"`         // 是否启用HTTP
	ListenHttps       bool     `json:"listen_https"`        // 是否启用HTTPS
	AllowedMethods    []string `json:"-"`                   // json配置里忽略，但解析Str可得
	AllowedMethodsStr string   `json:"allowed_methods_str"` // GET,POST,OPTIONS
	CertFile          string   `json:"cert_file"`           // cert.pem的具体路径
	KeyFile           string   `json:"key_file"`            // privkey.pem的具体路径
	DomainName        string   `json:"domain_name"`         // 域名；eg ： a.com eg : sdf.a.com
	HostName          string   `json:"host_name"`           // eg: localhost:9090
}

func main() {
	err := updateConfig() // 开始的时候读取一遍已有配置
	if err != nil {
		fmt.Println("start error:", err.Error())
		return
	}
	fmt.Println("server start")
	// server := http.NewServeMux()
	// server.HandleFunc("/", handler)
	// xxx := new(XXX)
	// domains := []string{"tnljqn.top", "baidu4560.com"}
	//
	// xxx.Tlsconfig = new(tls.Config)
	// xxx.Tlsconfig.Certificates = make([]tls.Certificate, 0)
	// for _, v := range domains {
	// 	cc, err := tls.LoadX509KeyPair("/etc/letsencrypt/live/"+v+"/cert.pem", "/etc/letsencrypt/live/"+v+"/privkey.pem")
	// 	if err != nil {
	// 		fmt.Println(err, cc)
	// 		continue
	// 	}
	// 	xxx.Tlsconfig.Certificates = append(xxx.Tlsconfig.Certificates, cc)
	// }
	// xxx.Tlsconfig.BuildNameToCertificate()
	// // xxx.Tlsconfig.NameToCertificate=make(map[string]*tls.Certificate)
	// // xxx.Tlsconfig.NameToCertificate["www.tnljqn.top"] = &cc
	// proxy := NewMultipleHostsReverseProxy(map[string]*url.URL{
	// 	"tnljqn.top": {
	// 		Scheme: "http",
	// 		Host:   "localhost:9091",
	// 	},
	// 	"baidu4560.com": {
	// 		Scheme: "http",
	// 		Host:   "localhost:9092",
	// 	},
	// 	"gg.baidu4560.com": {
	// 		Scheme: "http",
	// 		Host:   "localhost:9091",
	// 	},
	// 	"gg.tnljqn.top": {
	// 		Scheme: "http",
	// 		Host:   "localhost:9092",
	// 	},
	// })
	allServer = &http.Server{
		Addr: "0.0.0.0:443",
		// Handler:      proxy,
		// TLSConfig:    xxx.Tlsconfig,
		TLSNextProto: make(map[string]func(*http.Server, *tls.Conn, http.Handler), 0),
	}
	httpServer = &http.Server{
		Addr: "0.0.0.0:80",
	}
	// c,_:= xxx.GetCertificate()
	// srv.TLSConfig=xxx.Tlsconfig
	//
	// srv.TLSConfig.GetCertificate = xxx.GetCertificate
	// http.ListenAndServeTLS()
	go httpServer.ListenAndServe()
	// go http.ListenAndServe(":80", http.HandlerFunc(redirect))            // http重定向到https
	go http.ListenAndServe(":81", http.HandlerFunc(updateConfigRequest)) // 更新go配置

	updateServer()
	err = allServer.ListenAndServeTLS("", "")
	// xxx.CertConfigs[domain]
	if err != nil {
		fmt.Println("server error", err)
	}
	fmt.Println("server end")
}

// http重定向到https
// func redirect(w http.ResponseWriter, req *http.Request) {
// 	target := "https://" + req.Host + req.URL.Path
// 	if len(req.URL.RawQuery) > 0 {
// 		target += "?" + req.URL.RawQuery
// 	}
//
// 	http.Redirect(w, req, target, http.StatusTemporaryRedirect)
// }

// 更新配置请求
func updateConfigRequest(w http.ResponseWriter, r *http.Request) {
	err := updateConfig()
	if err != nil {
		fmt.Fprintf(w, err.Error())
		return
	}
	updateServer()
	fmt.Fprintf(w, "reload success\n")
}

// 更新配置
func updateConfig() error {
	fmt.Println("config reload")
	if allConfig == nil {
		allConfig = make(map[string]*Config)
	}
	err := filepath.Walk(configFile, func(path string, f os.FileInfo, err error) error {
		if f == nil {
			return err
		} // 报错则终止
		if f.IsDir() {
			return nil
		} // 文件夹忽略

		config := new(Config)
		if contents, err := ioutil.ReadFile(path); err == nil {
			err := json.Unmarshal(contents, &config)
			if err != nil {
				return err
			}
		} else {
			return err
		}
		if config.DomainName == "" || // 域名为空
			(!config.ListenHttp && !config.ListenHttps) || // 既不启用http也不启用https
			(config.ListenHttps && (config.CertFile == "" || config.KeyFile == "")) { // 启用了https却少了证书
			return nil // 无视了
		}
		if config.AllowedMethodsStr != "" { // 限制了请求
			strs := strings.Split(config.AllowedMethodsStr, ",")
			config.AllowedMethods = make([]string, 0)
			for _, v := range strs {
				ts := strings.ToUpper(strings.TrimSpace(v))
				if ts != "" {
					config.AllowedMethods = append(config.AllowedMethods, ts)
				}
			}
		}

		allConfig[config.DomainName] = config
		println(path)
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}

// 更新服务
func updateServer() {
	// 搞定https的
	tlsconfig := new(tls.Config)
	tlsconfig.Certificates = make([]tls.Certificate, 0)
	urls := make(map[string]*url.URL)
	// 将证书转成tlsconfig
	for _, v := range allConfig {
		if !v.ListenHttps {
			continue
		}
		cc, err := tls.LoadX509KeyPair(v.CertFile, v.KeyFile)
		if err != nil {
			fmt.Println(err, cc)
			continue
		}
		tlsconfig.Certificates = append(tlsconfig.Certificates, cc)
		// 反向代理的对应端口地址
		urls[v.DomainName] = &url.URL{Scheme: "http", Host: v.HostName}
	}
	tlsconfig.BuildNameToCertificate()
	allServer.TLSConfig = tlsconfig
	// proxy := NewMultipleHostsReverseProxy(urls)
	// 反向代理
	allServer.Handler = NewMultipleHostsReverseProxy(urls)
	// 搞定http的
	httpServer.Handler = &HTTPReverseProxy{rp: NewMultipleHostsReverseProxy(urls)}
}

// type XXX struct {
// 	// CertConfigs map[string]*CertificateConfig
// 	Tlsconfig *tls.Config
// }
// func (cm XXX) GetCertificate(clientInfo *tls.ClientHelloInfo) (*tls.Certificate, error) {
// 	fmt.Println("in GetCertificate")
// 	if x509Cert, ok := cm.Tlsconfig.NameToCertificate[clientInfo.ServerName]; ok {
// 		return x509Cert, nil
// 	}
// 	if a := strings.Index(clientInfo.ServerName, "."); a > 0 {
// 		clientInfo.ServerName = clientInfo.ServerName[a+1:]
// 	}
// 	if x509Cert, ok := cm.Tlsconfig.NameToCertificate[clientInfo.ServerName]; ok {
// 		return x509Cert, nil
// 	}
// 	clientInfo.Conn.Close()
// 	return nil, nil
// }
//
// type CertificateConfig struct {
// }

// 反向代理
func NewMultipleHostsReverseProxy(targets map[string]*url.URL) *httputil.ReverseProxy {
	director := func(req *http.Request) {
		// req.Header.Set("X-Forwarded-For", req.RemoteAddr)
		// fmt.Println("X-Forwarded-For",req.Header.Get("X-Forwarded-For"))
		if v, ok := targets[req.Host]; ok { // 判断固定域名
			req.URL.Scheme = v.Scheme
			req.URL.Host = v.Host
			// req.URL.Path = v.Path
			return
		}
		var s = ""
		if a := strings.Index(req.Host, "."); a > 0 {
			s = req.Host[a+1:]
		}
		if v, ok := targets[s]; ok { // 判断泛域名
			req.URL.Scheme = v.Scheme
			req.URL.Host = v.Host
			// req.URL.Path = v.Path
			return
		}

		req.Close = true // 没找到则关闭请求
		// target := targets[""]
	}
	return &httputil.ReverseProxy{Director: director}
}

// func HTTPNewMultipleHostsReverseProxy() *httputil.ReverseProxy {
// 	director := func(req *http.Request) {
// 		// 判断是否需要跳https
// 		// var config *Config
// 		// if v, ok := allConfig[req.Host]; ok {
// 		// 	config = v
// 		// }
// 		// if config == nil {
// 		// 	var s = ""
// 		// 	if a := strings.Index(req.Host, "."); a > 0 {
// 		// 		s = req.Host[a+1:]
// 		// 	}
// 		// 	if v, ok := allConfig[s]; ok {
// 		// 		config = v
// 		// 	}
// 		// }
// 		// if config == nil {
// 		// 	req.Close = true
// 		// 	return
// 		// }
// 		// if config.HttpToHttps {
// 		// 	req.URL.Scheme = "https"
// 		// 	req.URL.Host = req.Host
// 		// 	fmt.Println("HttpToHttps", req.URL.Host, req.URL.Scheme, req.URL.Path)
// 		// 	return
// 		// }
// 		req.URL.Scheme = "http"
// 		req.URL.Host = config.HostName
// 	}
// 	return &httputil.ReverseProxy{Director: director}
// }

func (p *HTTPReverseProxy) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	// 判断是否需要跳https
	var config *Config
	if v, ok := allConfig[req.Host]; ok {
		config = v
	}
	if config == nil {
		var s = ""
		if a := strings.Index(req.Host, "."); a > 0 {
			s = req.Host[a+1:]
		}
		if v, ok := allConfig[s]; ok {
			config = v
		}
	}
	if config == nil {
		return
	}
	if config.HttpToHttps {
		target := "https://" + req.Host + req.URL.Path
		if len(req.URL.RawQuery) > 0 {
			target += "?" + req.URL.RawQuery
		}

		http.Redirect(rw, req, target, http.StatusTemporaryRedirect)
		return
	}
	// req.URL.Scheme = "http"
	// req.URL.Host = config.HostName
	p.rp.ServeHTTP(rw, req) // 走反向代理通道
}

type HTTPReverseProxy struct {
	rp *httputil.ReverseProxy
}
