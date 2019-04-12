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

type Config struct {
	HttpToHttps       bool     `json:"http_to_https"`
	ListenHttp        bool     `json:"listen_http"`
	ListenHttps       bool     `json:"listen_https"`
	AllowedMethods    []string `json:"-"`
	AllowedMethodsStr string   `json:"allowed_methods_str"`
	CertFile          string   `json:"cert_file"`
	KeyFile           string   `json:"key_file"`
	DomainName        string   `json:"domain_name"`
}

func main() {
	fmt.Println("server start")
	// server := http.NewServeMux()
	// server.HandleFunc("/", handler)
	xxx := new(XXX)
	domains := []string{"tnljqn.top", "baidu4560.com"}

	xxx.Tlsconfig = new(tls.Config)
	xxx.Tlsconfig.Certificates = make([]tls.Certificate, 0)
	for _, v := range domains {
		cc, err := tls.LoadX509KeyPair("/etc/letsencrypt/live/"+v+"/cert.pem", "/etc/letsencrypt/live/"+v+"/privkey.pem")
		if err != nil {
			fmt.Println(err, cc)
			continue
		}
		xxx.Tlsconfig.Certificates = append(xxx.Tlsconfig.Certificates, cc)
	}
	xxx.Tlsconfig.BuildNameToCertificate()
	// xxx.Tlsconfig.NameToCertificate=make(map[string]*tls.Certificate)
	// xxx.Tlsconfig.NameToCertificate["www.tnljqn.top"] = &cc
	proxy := NewMultipleHostsReverseProxy(map[string]*url.URL{
		"tnljqn.top": {
			Scheme: "http",
			Host:   "localhost:9091",
		},
		"baidu4560.com": {
			Scheme: "http",
			Host:   "localhost:9092",
		},
		"gg.baidu4560.com": {
			Scheme: "http",
			Host:   "localhost:9091",
		},
		"gg.tnljqn.top": {
			Scheme: "http",
			Host:   "localhost:9092",
		},
	})
	srv := &http.Server{
		Addr:         "0.0.0.0:443",
		Handler:      proxy,
		TLSConfig:    xxx.Tlsconfig,
		TLSNextProto: make(map[string]func(*http.Server, *tls.Conn, http.Handler), 0),
	}
	// c,_:= xxx.GetCertificate()
	// srv.TLSConfig=xxx.Tlsconfig
	//
	// srv.TLSConfig.GetCertificate = xxx.GetCertificate
	// http.ListenAndServeTLS()
	go http.ListenAndServe(":80", http.HandlerFunc(redirect))
	go http.ListenAndServe(":81", http.HandlerFunc(updateConfig))
	err := srv.ListenAndServeTLS("", "")
	// xxx.CertConfigs[domain]
	if err != nil {
		fmt.Println("server error", err)
	}
	fmt.Println("server end")
}

// http重定向到https
func redirect(w http.ResponseWriter, req *http.Request) {
	target := "https://" + req.Host + req.URL.Path
	if len(req.URL.RawQuery) > 0 {
		target += "?" + req.URL.RawQuery
	}

	http.Redirect(w, req, target, http.StatusTemporaryRedirect)
}

// 更新配置
func updateConfig(w http.ResponseWriter, r *http.Request) {
	if allConfig == nil {

	}
	fmt.Println("config reload")
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
		if config.DomainName == "" {
			return nil // 无视了
		}
		allConfig[config.DomainName] = config
		println(path)
		return nil
	})
	if err != nil {
		fmt.Fprintf(w, fmt.Sprintf("filepath.Walk() returned %v\n", err))
		return
	}
	fmt.Fprintf(w, "reload success")
}

type XXX struct {
	// CertConfigs map[string]*CertificateConfig
	Tlsconfig *tls.Config
}

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
			req.URL.Path = v.Path
			return
		}
		var s = ""
		if a := strings.Index(req.Host, "."); a > 0 {
			s = req.Host[a+1:]
		}
		if v, ok := targets[s]; ok { // 判断泛域名
			req.URL.Scheme = v.Scheme
			req.URL.Host = v.Host
			req.URL.Path = v.Path
			return
		}

		req.Close = true // 关闭请求
		// target := targets[""]
	}
	return &httputil.ReverseProxy{Director: director}
}
