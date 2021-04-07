# 微博签到/app任务/超话签到

修改自 [happy888888/WeiboTask](https://github.com/happy888888/WeiboTask) ，没有 forked 因为基本上整个结构都不一样了。forked 起不到 up to date 的作用。

十分感谢原作者对微博app的解析。代码看了，很详细，参数拼接很细致。

1. 增加 go mod 支持
2. 重新规范命名
3. 解决原有 app 签到无法使用的问题
4. 超话过多时无法正常打卡的问题

其实代码中还有一些不够细致的地方，就懒得改了，能用就行。

## Docker 编译

作者原有的 docker 镜像也无法使用，在这里给懒得搭 go 环境的人一个 docker 编译方案
// btw 原作者的 ipk 也用不了，因为本人不用 windows ，那个能不能用就不知道了

```bash
docker run -v /home/user/docker/signup-weibo:/app -w /app golang:latest go build *.go
```

## 运行

因为 go 是针对平台编译，因此如果是 windows 就不能用上面的编译方式了，自己装个 go 环境编译吧。

默认读取 /etc/weibo/config.json ，如果不存在则读当前目录下的 config.json

## 参数

ALC 参数参考[原作者](https://github.com/happy888888/WeiboTask)

C 只有 android，iphone 下抓包拿不到正确的 S 参数。

通过 charles 抓包微博app，任一url为 api.weibo.com 下的链接后的get参数中均有 `s`/`c`/`from` 值。
