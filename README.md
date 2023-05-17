# dst-admin-go
> 饥荒联机版管理后台Go版本
>
> Date: 2023/05/11


## 预览

在线预览地址 http://1.12.223.51:8082/

![首页效果](./doc/登录.png)
![首页效果](./doc/房间.png)
![首页效果](./doc/mod.png)
![首页效果](./doc/mod配置.png)
![统计效果](./doc/统计.png)
![面板效果](./doc/面板.png)
![日志效果](./doc/日志.png)

## 部署
> 注意目录必须要有读写权限

1. 从release下载 dst-admin-go.tgz

2. 解压，上传到服务器

3. 修改config.yml 配置（端口）
    | 配置              | 解释                      | 是否必须|
    | ----------------- | ------------------------- | -------|
    | port          | 端口          | 是 |
    | db | 数据库名称（可以随便叫啥）              | 是 |
    | path           | 监听日志的路径(`/$HOME/.klei/DoNotStarveTogether/$Cluster1`)          | 是 |

    **参考配置**
    ```yml
    port: 8082
    db: dst-db
    path: /root/.klei/DoNotStarveTogether/MyDediServer
    ```

4. 修改dis_config 配置文件
   
    | 配置              | 解释                      | 是否必须|
    | ----------------- | ------------------------- | -------|
    | steamcmd          | steamcmd安装路径          | 是 |
    | force_install_dir | 饥荒安装路径              | 是 |
    | cluster           | 要启动房间的名称          | 是 |
    | backup            | 存档备份位置              | 是 |
    | mod_download_path | mod下载位置(路径要求存在) | 是 |
    
    
    
    如果你的启动路径不是 ` /$HOME/.klei/DoNotStarveTogether `, 请修改

    | 配置                          | 解释 |
    | ----------------------------- | ---- |
    | persistent_storage_root       |      |
    | donot_starve_server_directory |      |
    
    [参考命令](https://dontstarve.fandom.com/zh/wiki/%E5%A4%9A%E4%BA%BA%E7%89%88%E9%A5%A5%E8%8D%92%E7%8B%AC%E7%AB%8B%E6%9C%8D%E5%8A%A1%E5%99%A8?variant=zh#%E5%90%AF%E5%8A%A8%E5%8F%82%E6%95%B0)
    
    ```sh
    dontstarve_dedicated_server_nullrenderer -console -persistent_storage_root " + persistent_storage_root + "-conf_dir " + donot_starve_server_directory + " -cluster " + cluster + " -shard " + DST_CAVES + " ;"
    ```
    
    
    
5. 启动
    
    ```
    nohup ./dst-admin-go >dst-admin-go/log &
    ```
    
    

## 运行

**修改config.yml**
```
#端口
port: 8082
db: dst-db
#监听日志路径
path: /root/.klei/DoNotStarveTogether/MyDediServer
```

**修改dst_config**（也可以通过页面修改）
```
# steamcmd 位置
steamcmd=/root/steamcmd/
# steamcmd 饥荒安装的位置
force_install_dir=/root/dst/
# 要启动的服务器
cluster=cluster2
# 游戏备份的路径
backup=C:\Users\xm\Desktop\饥荒配置文件和建家截图\饥荒存储备份\backup
# 游戏mod下载的路径
mod_download_path=/download_mod
```

运行
```
go mod tidy
go run main.go
```

## 打包

### Linux 打包
```sh
go build
```

### window 打包

window 下打包 Linux 二进制 （由于sqlite受操作系统影响，Linux二进制请在Linux环境build）

***也有可能是我环境原因**

```
打开 cmd
set GOARCH=amd64
set GOOS=linux

go build
```
