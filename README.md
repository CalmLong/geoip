# GeoIP

包含 CN 的 `IPv4` 和 `IPv6` 数据

## 获取

Github Action 每日自动构建

**geoip.dat**: [https://raw.githubusercontent.com/CalmLong/geoip/release/geoip.dat](https://raw.githubusercontent.com/CalmLong/geoip/release/geoip.dat)

### 示例配置

```json
      {
        "type": "field",
        "ip": [
          "geoip:private",
          "geoip:cn"
        ],
        "outboundTag": "direct"
      }
```

## 运行

```bash
./geoip
```

稍等片刻输出 `geoip.dat`

* 可识别 `https_proxy` 环境变量

## 引用以下项目

* [github.com/gaoyifan/china-operator-ip](https://github.com/gaoyifan/china-operator-ip)
* [github.com/17mon/china_ip_list](https://github.com/17mon/china_ip_list)
* [github.com/v2fly/v2ray-core](https://github.com/v2fly/v2ray-core)
* [APNIC](https://ftp.apnic.net/stats/apnic/)
