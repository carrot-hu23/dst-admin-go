# py-dst-cli

Python 工具集，用于解析和处理饥荒联机版（Don't Starve Together）的各类配置文件、MOD 信息和游戏资源。

## 功能概述

本工具提供以下核心功能：

1. **世界配置解析** - 解析和生成 `leveldataoverride.lua` 配置文件
2. **MOD 信息提取** - 从 Steam Workshop 获取 MOD 详细信息
3. **游戏版本查询** - 获取当前最新的饥荒服务器版本号
4. **物品数据解析** - 解析 TooManyItemPlus 等 MOD 的物品数据
5. **资源提取** - 从游戏文件中提取图片资源（worldgen_customization.tex）

## 目录结构

```
py-dst-cli/
├── main.py                              # 主入口程序
├── dst_version.py                       # 版本查询工具
├── parse_world_setting.py               # 世界配置解析器
├── parse_mod.py                         # MOD 信息解析器
├── parse_TooManyItemPlus_items.py       # 物品数据解析器
├── parse_world_webp.py                  # 世界图片资源处理
├── dst_world_setting.json               # 世界配置 JSON 数据
├── requirements.txt                     # Python 依赖列表
└── steamapikey.txt                      # Steam API Key 配置
```

## 安装依赖

### 方法一：使用虚拟环境（推荐）

```bash
# 创建虚拟环境
python3 -m venv env

# 激活虚拟环境
source env/bin/activate  # Linux/Mac
# 或
env\Scripts\activate  # Windows

# 安装依赖库
pip install -r requirements.txt

# 使用完毕后退出虚拟环境
deactivate
```

### 方法二：直接安装

```bash
pip install -r requirements.txt
```

### 更新依赖列表

如果添加了新的依赖，使用以下命令更新 `requirements.txt`：
```bash
pip freeze > requirements.txt
```

## 使用方法

### 1. 配置 Steam API Key

在 `steamapikey.txt` 文件中配置你的 Steam Web API Key：

```
YOUR_STEAM_API_KEY_HERE
```

获取 API Key: https://steamcommunity.com/dev/apikey

### 2. 运行主程序

```bash
python main.py
```

### 3. 单独使用各个工具

#### 查询游戏版本

```bash
python dst_version.py
```

输出示例：
```json
{
  "version": "573123",
  "update_time": "2024-01-15 10:30:00"
}
```

#### 解析世界配置

```bash
python parse_world_setting.py
```

生成标准的 `dst_world_setting.json` 配置文件供管理面板使用。

#### 解析 MOD 信息

```bash
python parse_mod.py <mod_id>
```

示例：
```bash
python parse_mod.py 375859599  # 解析 Global Positions MOD
```

#### 解析物品数据

```bash
python parse_TooManyItemPlus_items.py
```

从 TooManyItemPlus MOD 中提取物品列表和属性。

#### 处理世界图片资源

```bash
python parse_world_webp.py
```

将 worldgen_customization.tex 文件转换为可用的图片格式。

## 获取游戏资源文件

### 方法一：从游戏安装目录获取

如果已安装饥荒联机版，可直接从安装目录获取：

**Windows:**
```
C:\Program Files (x86)\Steam\steamapps\common\Don't Starve Together\data\images
```

**Linux:**
```
~/.steam/steam/steamapps/common/Don't Starve Together/data/images
```

**Mac:**
```
~/Library/Application Support/Steam/steamapps/common/Don't Starve Together/data/images
```

关键文件：
- `worldgen_customization.tex` - 世界生成配置界面资源

### 方法二：通过 SteamCMD 下载

```bash
# 安装 SteamCMD
# Ubuntu/Debian:
sudo apt install steamcmd

# 下载游戏文件
steamcmd +login anonymous +app_update 343050 validate +quit

# 游戏文件位置:
# Linux: ~/.steam/steam/steamapps/common/Don't Starve Together/
```

### 方法三：使用 DepotDownloader

```bash
# 下载 DepotDownloader
wget https://github.com/SteamRE/DepotDownloader/releases/latest/download/DepotDownloader.zip

# 解压并运行
unzip DepotDownloader.zip -d DepotDownloader
cd DepotDownloader

# 下载饥荒服务器文件
./DepotDownloader -app 343050 -os linux -osarch 64 -dir ./dst-server -validate
```

## 常见配置项说明

### 世界配置（leveldataoverride.lua）

主要配置项：
- `world_size` - 世界大小: small/medium/large/huge
- `season_start` - 起始季节: autumn/winter/spring/summer
- `day` - 昼夜长度设置
- `weather` - 天气频率
- `resources` - 资源丰富度
- `creatures` - 生物数量
- `branching` - 地图分支复杂度

### MOD 配置（modoverrides.lua）

从 Steam Workshop 获取 MOD 配置选项，支持：
- MOD 名称、描述、作者
- 配置项列表和默认值
- 订阅数、更新时间
- 兼容性标签

## 依赖库

主要 Python 库：
- `requests` - HTTP 请求，用于 Steam API 调用
- `beautifulsoup4` - HTML 解析
- `Pillow` - 图片处理
- 其他工具库

## 与 DST Admin Go 集成

本工具集的输出数据会被 DST Admin Go 的以下模块使用：
- `internal/service/levelConfig/` - 世界配置管理
- `internal/service/mod/` - MOD 管理
- `internal/service/update/` - 游戏更新检测
- `internal/api/handler/dst_map_handler.go` - 地图生成

## 输出格式

所有工具默认输出 JSON 格式数据，可直接被 DST Admin Go 管理面板使用：

```json
{
  "success": true,
  "data": { ... },
  "message": "操作成功"
}
```

## 故障排查

### 无法获取 MOD 信息

- 检查 `steamapikey.txt` 是否正确配置
- 确认网络可以访问 Steam API
- 检查 MOD ID 是否有效
- Steam API 可能有速率限制，稍后重试

### 版本查询失败

- Steam 服务器可能暂时不可用，稍后重试
- 检查防火墙是否拦截了 Steam 相关域名
- 尝试使用代理访问

### 依赖安装失败

```bash
# 升级 pip 到最新版本
pip install --upgrade pip

# 使用国内镜像源（清华大学）
pip install -r requirements.txt -i https://pypi.tuna.tsinghua.edu.cn/simple

# 或使用阿里云镜像
pip install -r requirements.txt -i https://mirrors.aliyun.com/pypi/simple/
```

### 图片资源转换失败

- 确保安装了 Pillow 库：`pip install Pillow`
- 检查源文件格式是否正确
- 某些 .tex 文件可能需要特定的转换工具

## 常见用例

### 1. 批量查询 MOD 信息

```bash
# 创建 MOD ID 列表
echo "375859599" > mod_list.txt
echo "462372013" >> mod_list.txt

# 循环查询
while read mod_id; do
  python parse_mod.py $mod_id
done < mod_list.txt
```

### 2. 定时检查游戏版本更新

```bash
# 添加到 crontab（每小时检查一次）
0 * * * * cd /path/to/py-dst-cli && python dst_version.py >> version_log.txt
```

### 3. 导出世界配置模板

```bash
python parse_world_setting.py > world_template.json
```

## 注意事项

1. **Steam API 限制**: Steam Web API 有频率限制（每天 100,000 次调用），请勿频繁调用
2. **网络要求**: 部分功能需要访问 Steam 服务器，确保网络畅通
3. **编码问题**: 处理中文 MOD 时注意字符编码（统一使用 UTF-8）
4. **虚拟环境**: 建议使用虚拟环境隔离依赖，避免污染系统 Python 环境
5. **API Key 安全**: 不要将 `steamapikey.txt` 提交到公开仓库

## 开发者信息

本工具集用于辅助 DST Admin Go 项目，提供游戏数据的离线解析和在线查询能力。

相关项目：
- [DST Admin Go](https://github.com/hujinbo23/dst-admin-go) - 主项目
- [Steam Web API](https://steamcommunity.com/dev) - Steam API 文档
- [Don't Starve Together Wiki](https://dontstarve.fandom.com/wiki/Don%27t_Starve_Together) - 游戏官方 Wiki

## 贡献指南

如果需要添加新功能：
1. 在对应的 Python 文件中添加函数
2. 更新 `main.py` 以集成新功能
3. 更新 `requirements.txt`（如果有新依赖）
4. 更新本 README 文档

