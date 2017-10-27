import requests

proxy = dict(
    http='socks5:127.0.0.1:1080',
    https='socks5:127.0.0.1:1080',
)

req = requests.get('http://ip.cn', proxies=proxy)
print(req.text)
