# 小说服务
基于DiscoverNewNovel爬取的数据，使用beego框架提供的web json服务。

## 提供的服务
1. 分页拉取
2. 搜索
3. 章节分页
4. 章节详情

## 整体结构
小说基本信息来源于DiscoverNewNovel爬取的数据，存储于本地的mongodb中。章节详情实时拉取，为了加快速度，增加了redis缓存已经拉取过的章节详情。

## 使用的三方包

```
github.com/astaxie/beego
github.com/go-redis/redis
golang.org/x/text/transform
github.com/PuerkitoBio/goquery
golang.org/x/text/encoding/simplifiedchinese
gopkg.in/mgo.v2/bson
gopkg.in/mgo.v2
```