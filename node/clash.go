package node

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"gopkg.in/yaml.v3"
)

type Proxy struct {
	Name               string                 `yaml:"name,omitempty"`
	Type               string                 `yaml:"type,omitempty"`
	Server             string                 `yaml:"server,omitempty"`
	Port               int                    `yaml:"port,omitempty"`
	Cipher             string                 `yaml:"cipher,omitempty"`
	Password           string                 `yaml:"password,omitempty"`
	Client_fingerprint string                 `yaml:"client-fingerprint,omitempty"`
	Tfo                bool                   `yaml:"tfo,omitempty"`
	Udp                bool                   `yaml:"udp,omitempty"`
	Skip_cert_verify   bool                   `yaml:"skip-cert-verify,omitempty"`
	Tls                bool                   `yaml:"tls,omitempty"`
	Servername         string                 `yaml:"servername,omitempty"`
	Flow               string                 `yaml:"flow,omitempty"`
	AlterId            string                 `yaml:"alterId,omitempty"`
	Network            string                 `yaml:"network,omitempty"`
	Reality_opts       map[string]interface{} `yaml:"reality-opts,omitempty"`
	Ws_opts            map[string]interface{} `yaml:"ws-opts,omitempty"`
	Grpc_opts          map[string]interface{} `yaml:"grpc-opts,omitempty"`
	Auth_str           string                 `yaml:"auth_str,omitempty"`
	Auth               string                 `yaml:"auth,omitempty"`
	Up                 int                    `yaml:"up,omitempty"`
	Down               int                    `yaml:"down,omitempty"`
	Alpn               []string               `yaml:"alpn,omitempty"`
	Sni                string                 `yaml:"sni,omitempty"`
	Obfs               string                 `yaml:"obfs,omitempty"`
	Obfs_password      string                 `yaml:"obfs-password,omitempty"`
	Protocol           string                 `yaml:"protocol,omitempty"`
	Uuid               string                 `yaml:"uuid,omitempty"`
	Peer               string                 `yaml:"peer,omitempty"`
	Congestion_control string                 `yaml:"congestion_control,omitempty"`
	Udp_relay_mode     string                 `yaml:"udp_relay_mode,omitempty"`
	Disable_sni        bool                   `yaml:"disable_sni,omitempty"`
}

type ProxyGroup struct {
	Proxies []string `yaml:"proxies"`
}
type Config struct {
	Proxies      []Proxy      `yaml:"proxies"`
	Proxy_groups []ProxyGroup `yaml:"proxy-groups"`
}

// 删除opts中的空值
func DeleteOpts(opts map[string]interface{}) {
	for k, v := range opts {
		switch v := v.(type) {
		case string:
			if v == "" {
				delete(opts, k)
			}
		case map[string]interface{}:
			DeleteOpts(v)
			if len(v) == 0 {
				delete(opts, k)
			}
		}
	}
}
func convertToInt(value interface{}) (int, error) {
	switch v := value.(type) {
	case int:
		return v, nil
	case float64:
		return int(v), nil
	case string:
		return strconv.Atoi(v)
	default:
		return 0, fmt.Errorf("unexpected type %T", v)
	}
}

// EncodeClash 用于生成 Clash 配置文件
func EncodeClash(urls []string, sqlconfig SqlConfig) ([]byte, error) {
	// 传入urls，解析urls，生成proxys
	// yamlfile 为模板文件
	var proxys []Proxy

	for _, link := range urls {
		Scheme := strings.Split(link, "://")[0]
		switch {
		case Scheme == "ss":
			ss, err := DecodeSSURL(link)
			if err != nil {
				log.Println(err)
				continue
			}
			// 如果没有名字，就用服务器地址作为名字
			if ss.Name == "" {
				ss.Name = fmt.Sprintf("%s:%d", ss.Server, ss.Port)
			}
			ssproxy := Proxy{
				Name:             ss.Name,
				Type:             "ss",
				Server:           ss.Server,
				Port:             ss.Port,
				Cipher:           ss.Param.Cipher,
				Password:         ss.Param.Password,
				Udp:              sqlconfig.Udp,
				Skip_cert_verify: sqlconfig.Cert,
			}
			proxys = append(proxys, ssproxy)
		case Scheme == "ssr":
			ssr, err := DecodeSSRURL(link)
			if err != nil {
				log.Println(err)
			}
			// 如果没有名字，就用服务器地址作为名字
			if ssr.Qurey.Remarks == "" {
				ssr.Qurey.Remarks = fmt.Sprintf("%s:%d", ssr.Server, ssr.Port)
			}
			ssrproxy := Proxy{
				Name:             ssr.Qurey.Remarks,
				Type:             "ssr",
				Server:           ssr.Server,
				Port:             ssr.Port,
				Cipher:           ssr.Method,
				Password:         ssr.Password,
				Obfs:             ssr.Obfs,
				Obfs_password:    ssr.Qurey.Obfsparam,
				Protocol:         ssr.Protocol,
				Udp:              sqlconfig.Udp,
				Skip_cert_verify: sqlconfig.Cert,
			}
			proxys = append(proxys, ssrproxy)
		case Scheme == "trojan":
			trojan, err := DecodeTrojanURL(link)
			if err != nil {
				log.Println(err)
				continue
			}
			// 如果没有名字，就用服务器地址作为名字
			if trojan.Name == "" {
				trojan.Name = fmt.Sprintf("%s:%d", trojan.Hostname, trojan.Port)
			}
			ws_opts := map[string]interface{}{
				"path": trojan.Query.Path,
				"headers": map[string]interface{}{
					"Host": trojan.Query.Host,
				},
			}
			DeleteOpts(ws_opts)
			trojanproxy := Proxy{
				Name:               trojan.Name,
				Type:               "trojan",
				Server:             trojan.Hostname,
				Port:               trojan.Port,
				Password:           trojan.Password,
				Client_fingerprint: trojan.Query.Fp,
				Sni:                trojan.Query.Sni,
				Network:            trojan.Query.Type,
				Flow:               trojan.Query.Flow,
				Alpn:               trojan.Query.Alpn,
				Ws_opts:            ws_opts,
				Udp:                sqlconfig.Udp,
				Skip_cert_verify:   sqlconfig.Cert,
			}
			proxys = append(proxys, trojanproxy)
		case Scheme == "vmess":
			vmess, err := DecodeVMESSURL(link)
			if err != nil {
				log.Println(err)
				continue
			}
			// 如果没有名字，就用服务器地址作为名字
			if vmess.Ps == "" {
				vmess.Ps = fmt.Sprintf("%s:%s", vmess.Add, vmess.Port)
			}
			ws_opts := map[string]interface{}{
				"path": vmess.Path,
				"headers": map[string]interface{}{
					"Host": vmess.Host,
				},
			}
			DeleteOpts(ws_opts)
			tls := false
			if vmess.Tls != "none" && vmess.Tls != "" {
				tls = true
			}
			port, _ := convertToInt(vmess.Port)
			aid, _ := convertToInt(vmess.Aid)
			vmessproxy := Proxy{
				Name:             vmess.Ps,
				Type:             "vmess",
				Server:           vmess.Add,
				Port:             port,
				Cipher:           vmess.Scy,
				Uuid:             vmess.Id,
				AlterId:          strconv.Itoa(aid),
				Network:          vmess.Net,
				Tls:              tls,
				Ws_opts:          ws_opts,
				Udp:              sqlconfig.Udp,
				Skip_cert_verify: sqlconfig.Cert,
			}
			proxys = append(proxys, vmessproxy)
		case Scheme == "vless":
			vless, err := DecodeVLESSURL(link)
			if err != nil {
				log.Println(err)
				continue
			}
			// 如果没有名字，就用服务器地址作为名字
			if vless.Name == "" {
				vless.Name = fmt.Sprintf("%s:%d", vless.Server, vless.Port)
			}
			ws_opts := map[string]interface{}{
				"path": vless.Query.Path,
				"headers": map[string]interface{}{
					"Host": vless.Query.Host,
				},
			}
			reality_opts := map[string]interface{}{
				"public-key": vless.Query.Pbk,
				"short-id":   vless.Query.Sid,
			}
			grpc_opts := map[string]interface{}{
				"grpc-mode":         "gun",
				"grpc-service-name": vless.Query.ServiceName,
			}
			if vless.Query.Mode == "multi" {
				grpc_opts["grpc-mode"] = "multi"
			}
			DeleteOpts(ws_opts)
			DeleteOpts(reality_opts)
			DeleteOpts(grpc_opts)
			tls := false
			if vless.Query.Security != "" {
				tls = true
			}
			if vless.Query.Security == "none" {
				tls = false
			}
			vlessproxy := Proxy{
				Name:               vless.Name,
				Type:               "vless",
				Server:             vless.Server,
				Port:               vless.Port,
				Servername:         vless.Query.Sni,
				Uuid:               vless.Uuid,
				Client_fingerprint: vless.Query.Fp,
				Network:            vless.Query.Type,
				Flow:               vless.Query.Flow,
				Alpn:               vless.Query.Alpn,
				Ws_opts:            ws_opts,
				Reality_opts:       reality_opts,
				Grpc_opts:          grpc_opts,
				Udp:                sqlconfig.Udp,
				Skip_cert_verify:   sqlconfig.Cert,
				Tls:                tls,
			}
			proxys = append(proxys, vlessproxy)
		case Scheme == "hy" || Scheme == "hysteria":
			hy, err := DecodeHYURL(link)
			if err != nil {
				log.Println(err)
				continue
			}
			// 如果没有名字，就用服务器地址作为名字
			if hy.Name == "" {
				hy.Name = fmt.Sprintf("%s:%d", hy.Host, hy.Port)
			}
			hyproxy := Proxy{
				Name:             hy.Name,
				Type:             "hysteria",
				Server:           hy.Host,
				Port:             hy.Port,
				Auth_str:         hy.Auth,
				Up:               hy.UpMbps,
				Down:             hy.DownMbps,
				Alpn:             hy.ALPN,
				Peer:             hy.Peer,
				Udp:              sqlconfig.Udp,
				Skip_cert_verify: sqlconfig.Cert,
			}
			proxys = append(proxys, hyproxy)
		case Scheme == "hy2" || Scheme == "hysteria2":
			hy2, err := DecodeHY2URL(link)
			if err != nil {
				log.Println(err)
				continue
			}
			// 如果没有名字，就用服务器地址作为名字
			if hy2.Name == "" {
				hy2.Name = fmt.Sprintf("%s:%d", hy2.Host, hy2.Port)
			}
			hyproxy2 := Proxy{
				Name:             hy2.Name,
				Type:             "hysteria2",
				Server:           hy2.Host,
				Port:             hy2.Port,
				Auth_str:         hy2.Auth,
				Sni:              hy2.Sni,
				Alpn:             hy2.ALPN,
				Obfs:             hy2.Obfs,
				Password:         hy2.Password,
				Obfs_password:    hy2.ObfsPassword,
				Udp:              sqlconfig.Udp,
				Skip_cert_verify: sqlconfig.Cert,
			}
			proxys = append(proxys, hyproxy2)
		case Scheme == "tuic":
			tuic, err := DecodeTuicURL(link)
			if err != nil {
				log.Println(err)
				continue
			}
			// 如果没有名字，就用服务器地址作为名字
			if tuic.Name == "" {
				tuic.Name = fmt.Sprintf("%s:%d", tuic.Host, tuic.Port)
			}
			disable_sni := false
			if tuic.Disable_sni == 1 {
				disable_sni = true
			}
			tuicproxy := Proxy{
				Name:               tuic.Name,
				Type:               "tuic",
				Server:             tuic.Host,
				Port:               tuic.Port,
				Password:           tuic.Password,
				Uuid:               tuic.Uuid,
				Congestion_control: tuic.Congestion_control,
				Alpn:               tuic.Alpn,
				Udp_relay_mode:     tuic.Udp_relay_mode,
				Disable_sni:        disable_sni,
				Sni:                tuic.Sni,
				Udp:                sqlconfig.Udp,
				Skip_cert_verify:   sqlconfig.Cert,
			}
			proxys = append(proxys, tuicproxy)
		}
	}
	// 生成Clash配置文件
	return DecodeClash(proxys, sqlconfig.Clash)
}

// DecodeClash 用于解析 Clash 配置文件
func DecodeClash(proxys []Proxy, yamlfile string) ([]byte, error) {
	// 读取 YAML 文件
	var data []byte
	var err error
	if strings.Contains(yamlfile, "://") {
		resp, err := http.Get(yamlfile)
		if err != nil {
			log.Println("http.Get error", err)
			return nil, err
		}
		defer resp.Body.Close()
		data, err = io.ReadAll(resp.Body)
		if err != nil {
			log.Printf("error: %v", err)
			return nil, err
		}
	} else {
		data, err = os.ReadFile(yamlfile)
		if err != nil {
			log.Printf("error: %v", err)
			return nil, err
		}
	}
	// 解析 YAML 文件
	config := make(map[interface{}]interface{})
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		log.Printf("error: %v", err)
		return nil, err
	}

	// 检查 "proxies" 键是否存在于 config 中
	proxies, ok := config["proxies"].([]interface{})
	if !ok {
		// 如果 "proxies" 键不存在，创建一个新的切片
		proxies = []interface{}{}
	}
	// 定义一个代理列表名字
	ProxiesNameList := []string{}
	// 添加新代理
	for _, p := range proxys {
		ProxiesNameList = append(ProxiesNameList, p.Name)
		proxies = append(proxies, p)
	}
	// proxies = append(proxies, newProxy)
	config["proxies"] = proxies
	// 往ProxyGroup中插入代理列表
	// ProxiesNameList := []string{"newProxy", "ceshi"}
	proxyGroups := config["proxy-groups"].([]interface{})
	for i, pg := range proxyGroups {
		proxyGroup, ok := pg.(map[string]interface{})
		if !ok {
			continue
		}
		// 如果 proxyGroup["proxies"] 是 nil，初始化它为一个空的切片
		if proxyGroup["proxies"] == nil {
			proxyGroup["proxies"] = []interface{}{}
		}
		// 如果为链式代理的话则不插入返回
		// log.Print("代理类型为:", proxyGroup["type"])
		if proxyGroup["type"] == "relay" {
			break
		}
		// 清除 nil 值
		var validProxies []interface{}
		for _, p := range proxyGroup["proxies"].([]interface{}) {
			if p != nil {
				validProxies = append(validProxies, p)
			}
		}
		// 添加新代理
		for _, newProxy := range ProxiesNameList {
			validProxies = append(validProxies, newProxy)
		}
		proxyGroup["proxies"] = validProxies
		proxyGroups[i] = proxyGroup
	}

	config["proxy-groups"] = proxyGroups

	// 将修改后的内容写回文件
	newData, err := yaml.Marshal(config)
	if err != nil {
		log.Printf("error: %v", err)
	}
	return newData, nil
}
