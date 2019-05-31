## 环境变量说明

- `PORT`: 服务器监听端口，*默认值为8080*
- `BASE_URL`: 网站静态资源根级URL，**必须提供且不能为空**
- `INDEX_PATH`: 首页网页相对路径，*默认值为`index.html`*
- `FETCH_TIMEOUT`: 获取静态资源超时时间，单位为秒，*默认值为`3`*
- `REFRESH_INTERVAL`: 刷新缓存的最小时间间隔，单位为秒。在刷新间隔内即使远端内容有变化也不会刷新缓存。*默认值为`30`*