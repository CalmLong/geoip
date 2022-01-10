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

## 引用以下项目

* [github.com/reflect2/china-ip-list](https://github.com/reflect2/china-ip-list)
* [github.com/v2fly/v2ray-core](https://github.com/v2fly/v2ray-core)
