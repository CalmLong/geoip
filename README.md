# GeoIP for V2Ray

适用于 V2Ray 的 IP 数据库

## 运行

```bash
./geoip
```

稍等片刻输出 `geoip.dat`

* 可识别 `https_proxy` 环境变量

## 使用

包含[这些](https://github.com/metowolf/iplist/blob/master/docs/country.md)国家的和地区的 IP 数据，对于大陆则使用[简化版](https://github.com/metowolf/iplist#%E5%A4%A7%E9%99%86-ip-%E6%AE%B5)的数据

### 示例配置

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
          "geoip:us",
          "geoip:jp"
        ],
        "outboundTag": "proxy"
      }
```

## 引用以下项目

* [github.com/metowolf/iplist](https://github.com/metowolf/iplist)
* [github.com/v2fly/v2ray-core](https://github.com/v2fly/v2ray-core)