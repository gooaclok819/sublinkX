package node

import (
	"encoding/base64"
	"encoding/json"
	"testing"

	"gopkg.in/yaml.v3"
)

// encodeOne 把单个节点链接转成 clash 配置并取回第一个节点的字段表
// （DecodeClash 需要一份模板文件，用仓库自带的 ./template/clash.yaml）
func encodeOne(t *testing.T, link string) map[string]interface{} {
	t.Helper()
	out, err := EncodeClash([]string{link}, SqlConfig{Clash: "../template/clash.yaml"})
	if err != nil {
		t.Fatalf("EncodeClash 出错: %v", err)
	}
	var doc struct {
		Proxies []map[string]interface{} `yaml:"proxies"`
	}
	if err := yaml.Unmarshal(out, &doc); err != nil {
		t.Fatalf("clash YAML 解析失败: %v", err)
	}
	if len(doc.Proxies) != 1 {
		t.Fatalf("期望 1 个节点，实际 %d 个", len(doc.Proxies))
	}
	return doc.Proxies[0]
}

func vmessLink(t *testing.T, obj map[string]interface{}) string {
	t.Helper()
	b, err := json.Marshal(obj)
	if err != nil {
		t.Fatalf("vmess json 编码失败: %v", err)
	}
	return "vmess://" + base64.StdEncoding.EncodeToString(b)
}

// 每个节点的传输层参数必须写在对应字段里：
// xhttp -> xhttp-opts，ws -> ws-opts，grpc -> grpc-opts，且不能互相串门。
func TestClashTransportOpts(t *testing.T) {
	baseVmess := map[string]interface{}{
		"v": "2", "ps": "vm", "add": "example.com", "port": "443",
		"id": "b831381d-6324-4d53-ad4f-8cda48b30811", "aid": "0", "scy": "auto",
		"type": "none", "tls": "tls", "sni": "example.com",
	}

	xhttpVmess := map[string]interface{}{}
	for k, v := range baseVmess {
		xhttpVmess[k] = v
	}
	xhttpVmess["net"], xhttpVmess["path"], xhttpVmess["host"] = "xhttp", "/xhp", "example.com"

	wsVmess := map[string]interface{}{}
	for k, v := range baseVmess {
		wsVmess[k] = v
	}
	wsVmess["net"], wsVmess["path"], wsVmess["host"] = "ws", "/ws", "example.com"

	grpcVmess := map[string]interface{}{}
	for k, v := range baseVmess {
		grpcVmess[k] = v
	}
	grpcVmess["net"], grpcVmess["path"] = "grpc", "vm-grpc"

	cases := []struct {
		desc        string
		link        string
		wantNetwork string
		wantOpts    string
		wantPath    string
		absent      []string
	}{
		{
			desc:        "vless + xhttp",
			link:        "vless://b831381d-6324-4d53-ad4f-8cda48b30811@example.com:8444?type=xhttp&security=tls&path=%2Fxhp&host=example.com&mode=auto&sni=example.com#n",
			wantNetwork: "xhttp", wantOpts: "xhttp-opts", wantPath: "/xhp",
			absent: []string{"ws-opts", "grpc-opts"},
		},
		{
			desc:        "vless + ws",
			link:        "vless://b831381d-6324-4d53-ad4f-8cda48b30811@example.com:443?type=ws&security=tls&path=%2Fws&host=example.com&sni=example.com#n",
			wantNetwork: "ws", wantOpts: "ws-opts", wantPath: "/ws",
			absent: []string{"xhttp-opts"},
		},
		{
			desc:        "vless + grpc",
			link:        "vless://b831381d-6324-4d53-ad4f-8cda48b30811@example.com:443?type=grpc&security=tls&serviceName=gs&sni=example.com#n",
			wantNetwork: "grpc", wantOpts: "grpc-opts", wantPath: "",
			absent: []string{"ws-opts", "xhttp-opts"},
		},
		{
			desc:        "vless + reality + tcp（回归）",
			link:        "vless://b831381d-6324-4d53-ad4f-8cda48b30811@example.com:443?type=tcp&security=reality&fp=chrome&pbk=EXAMPLEfakePublicKeyForTestOnly0000&sid=0123456789abcdef&sni=example.com&flow=xtls-rprx-vision#n",
			wantNetwork: "tcp", wantOpts: "reality-opts", wantPath: "",
			absent: []string{"ws-opts", "grpc-opts", "xhttp-opts"},
		},
		{
			desc:        "vmess + xhttp",
			link:        vmessLink(t, xhttpVmess),
			wantNetwork: "xhttp", wantOpts: "xhttp-opts", wantPath: "/xhp",
			absent: []string{"ws-opts", "grpc-opts"},
		},
		{
			desc:        "vmess + ws",
			link:        vmessLink(t, wsVmess),
			wantNetwork: "ws", wantOpts: "ws-opts", wantPath: "/ws",
			absent: []string{"xhttp-opts"},
		},
		{
			desc:        "vmess + grpc",
			link:        vmessLink(t, grpcVmess),
			wantNetwork: "grpc", wantOpts: "grpc-opts", wantPath: "",
			absent: []string{"ws-opts", "xhttp-opts"},
		},
		{
			desc:        "trojan + xhttp",
			link:        "trojan://pass@example.com:8443?type=xhttp&security=tls&path=%2Ftp&host=example.com&sni=example.com#n",
			wantNetwork: "xhttp", wantOpts: "xhttp-opts", wantPath: "/tp",
			absent: []string{"ws-opts", "grpc-opts"},
		},
		{
			desc:        "trojan + ws",
			link:        "trojan://pass@example.com:443?type=ws&security=tls&path=%2Fws&host=example.com&sni=example.com#n",
			wantNetwork: "ws", wantOpts: "ws-opts", wantPath: "/ws",
			absent: []string{"xhttp-opts"},
		},
	}

	for _, c := range cases {
		t.Run(c.desc, func(t *testing.T) {
			p := encodeOne(t, c.link)
			if got := p["network"]; got != c.wantNetwork {
				t.Errorf("network = %v，期望 %v", got, c.wantNetwork)
			}
			opts, ok := p[c.wantOpts].(map[string]interface{})
			if !ok {
				t.Fatalf("缺少 %s（实际字段: %v）", c.wantOpts, keys(p))
			}
			if c.wantPath != "" {
				if got := opts["path"]; got != c.wantPath {
					t.Errorf("%s.path = %v，期望 %v", c.wantOpts, got, c.wantPath)
				}
			}
			for _, a := range c.absent {
				if _, bad := p[a]; bad {
					t.Errorf("不该出现 %s（实际值 %v）", a, p[a])
				}
			}
		})
	}
}

func keys(m map[string]interface{}) []string {
	out := make([]string, 0, len(m))
	for k := range m {
		out = append(out, k)
	}
	return out
}
