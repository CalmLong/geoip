# GeoIP

仅包含 CN 的 `IPv4` 和 `IPv6` 数据

## 运行

```bash
./geoip
```

稍等片刻输出 `geoip.dat`

* 可识别 `https_proxy` 环境变量

## 使用

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

* [github.com/gaoyifan/china-operator-ip](https://github.com/gaoyifan/china-operator-ip)
* [github.com/v2fly/v2ray-core](https://github.com/v2fly/v2ray-core)
