# 饥荒服务器多开

## 多个房间

目前多开方式比较简单，默认你已经配置好了环境

1. 重新使用steamcmd 生成一个饥荒服务
    ``` sh
    ./steamcmd.sh +login anonymous +force_install_dir ~/dontstarve_dedicated_server +app_update 343050 validate +quit
    ```
2. 把 libstdc++.so.6  拷贝过去
    ```
    cp ~/steamcmd/linux32/libstdc++.so.6 ~/dontstarve_dedicated_server/bin/lib32/
    ```
3. 把 dst-admin-go 复制一份

    同时把 dst-admin-go 文件夹 里面的 dst-admin-go 重新命名为 dst-admin-go2,

    同时修改 config.yml 里面的端口
    ```
    port: xxx
    ```
4. 在 ~/.klei/DoNotStarveTogether 路径下在创建一个新的存档

   ```
    ~/.klei/DoNotStarveTogether/Cluster2
   ```
5. 启动上面的 dst-admin-go2 二进制文件
    进入网页，去设置里面修改 这两处配置
    ![系统设置](./image/%E9%85%8D%E7%BD%AE.png)
6. 然后就可以启动了
    注意去 房间配置 页面 修改下 主世界的端口


## 多层世界
部署和上面是一样的，你只需要修改 cluster.ini 文件 和 server.ini 就ok 了，具体参考 https://steamcommunity.com/sharedfiles/filedetails/?id=714846590
