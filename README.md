# B站推荐视频 Tag 统计爬虫

爬取 B站 首页推荐视频，提取视频 Tag 并统计出现频次。

## 使用方法

### 1. 获取你的b站 Cookie

1. 登录 B站网页版
2. 按 F12 打开开发者工具
3. 切换到 Network -> 看起来像b站的网站 -> 复制cookie栏
4. 复制 Cookie 值，填入 config.json

### 2. 运行

直接运行：
```bash
go run main.go
```

或先编译再运行：
```bash
go build -o biliTagAnalyse.exe
./biliTagAnalyse.exe
```

### 3. 查看结果

结果默认在 `results/tags_stats.json`

## 配置

config.json 配置文件说明：

| 参数 | 说明 | 默认值 |
|------|------|--------|
| cookie | B站登录 Cookie，必填 | - |
| crawl_count | 爬取轮数 | 5 |
| crawl_interval | 每轮间隔（秒） | 300 |
| request_interval | 请求间隔（毫秒） | 500 |
| max_concurrent | 最大并发数 | 3 |
| retry_count | 失败重试次数 | 3 |
| output_file | 结果输出路径 | results/tags_stats.json |

主要修改 cookie 即可。

## 注意

- 合理设置爬取间隔，避免被封
- Cookie 有效期较短，需要定期更新
