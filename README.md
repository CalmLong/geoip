# GeoIP for V2Ray

适用于 V2Ray 的 IP 数据库

## 运行

将 `MaxMind` 秘钥添加到环境变量

```bash
export MAX_MIND_KEY=YOUR_KEY
```

运行

```bash
./geoip
```

稍等片刻输出 `geoip.dat`

* 可识别 `http_proxy` 环境变量

## 使用

和官方版本不同的是使用本工具输出的 `geoip.dat` 仅包含以下

* `cn` 仅包含 CN 大陆内的 IP 数据
* `ncn` 包含除 CN 大陆以外的 IP 数据
* `private` 私有地址


```json
      {
        "type": "field",
        "ip": [
          "geoip:private",
          "geoip:cn"
        ],
        "outboundTag": "direct"
      },
      {
        "type": "field",
        "ip": [
          "geoip:ncn"
        ],
        "outboundTag": "proxy"
      }
```

## 引用以下项目

* [github.com/v2fly/geoip](https://github.com/v2fly/geoip)