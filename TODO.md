# TODO new feature



## 添加申请权限页面

### 房间密码
对于某些房间，需要用户提供， 用户名，科雷ID，SteamId，才能申请密码

并且管理员可以设置房间密码的有效时间

对于黑名单用户，不可申请

### 申请管理员权限
由于某些用户比较活跃，可以给予某些房间的管理员身份，帮忙维护房间管理

这种可以管理手动批准，也可以设置游戏时长达到某个时长

### 申请管理控制台权限
由于某些用户在群里比较活跃，平时维护群，可以给与管理控制台权限，但是权限是受限制的

由于比较特殊，需要管理员手动批准申请

## 添加举报熊孩子举报页面
对于某些熊孩子，玩家可以在页面 填写熊孩子游戏名称、steamId，科雷Id 来举报熊孩子，后台将审核

## 房间生成页面


## 是否考虑支持一台服务器支持多个游戏房间呢？



## 自动更新游戏


## 游戏crash掉了，自动重启饥荒服务

由于某些原因，导致房间炸档，服务器download了

起个后台监控，当饥荒服务器挂掉后，自动重启，


## 修改启动方式 (暂时pass)

本地写个文件 start_world
```
Master

Caves
```

根据 start_world 里面的内容来启动世界

```
./dontstarve_dedicated_server_nullrenderer -console -cluster " + cluster + " -shard " + DST_MASTER + "  ;"
```

启动之前需要杀掉全部的 Master 和 Caves 世界

配置文件详情
https://steamcommunity.com/sharedfiles/filedetails/?id=1616647350


## 多开

一般的启动流程

启动
1. 检查是否启动了
    ps -ef | grep -v grep | grep 'Master' | sed -n '1P' | awk '{print $2}'
2. 杀掉Master 进程
    "ps -ef | grep -v grep |grep '" + DST_MASTER + "' |sed -n '1P'|awk '{print $2}' |xargs kill -9"

3. screen -d -m -S 启动进程

./dontstarve_dedicated_server_nullrenderer -console -cluster Cluster1 -shard Master


### 设置
1. 用户创建存档 Cluster2
2. 配置世界（可视化界面操作）
3. 后台在 ~/.klei/DontStarveTogether/Cluser2 创建相应的文件

|-- Cluster1
|   |-- adminlist.txt
|   |-- blocklist.txt
|   |-- Caves
|   |   |-- leveldataoverride.lua
|   |   |-- modoverrides.lua
|   |   `-- server.ini
|   |-- cluster.ini
|   |-- cluster_token.txt
|   |-- Master
|   |   |-- leveldataoverride.lua
|   |   |-- modoverrides.lua
|   |   `-- server.ini


4. 返回用户 ~/.klei/DontStarveTogether 创建的世界 集合 [Cluster1, Cluster2]
ClusterList:[
    {
    name: Cluster1,
    masterStatus: false,
    cavesStatus: true,
    mem: 1200,
    cpu: 50%,
    onlinePlayers: 8/10,
    days: 1,
    season: spring
    },
    {
    name: Cluster1,
    masterStatus: false,
    cavesStatus: true,
    mem: 1200,
    cpu: 50%,
    onlinePlayers: 8/10,
    days: 1,
    season: spring
    },
]

启动
return "cd " + dst_install_dir + "/bin ; screen -d -m -S \"" + SCREEN_WORK_MASTER_NAME + "\"  ./dontstarve_dedicated_server_nullrenderer -console -cluster " + cluster + " -shard " + DST_MASTER + "  ;"





 ./dontstarve_dedicated_server_nullrenderer -console -cluster MyCluster3 -shard Caves