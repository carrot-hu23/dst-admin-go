import lupa
from functools import reduce

import steam.client
import steam.client.cdn
import steam.core.cm
import steam.webapi
from steam.exceptions import SteamError

anonymous = steam.client.SteamClient()
anonymous.anonymous_login()
steamcdn = steam.client.cdn.CDNClient(anonymous)

encoding = 'utf-8'

#TODO 解决并发问题
'''
lock
使用一个cache{modId: xxx, semaphore: xxx, mod: xxx}
init cache
初始化信号量
unlock

执行IO

其他等待 mod 数据返回

'''
def get_mod_info_dict(modId:int):

    status, modinfo = 0, {'modinfo':b'', 'modinfo_chs': b''}
    mod_item = steamcdn.get_manifest_for_workshop_item(modId)
    modinfo_names = ['modinfo.lua', 'modinfo_chs.lua']
    modinfo_list = list(filter(lambda x: x.filename in modinfo_names, mod_item.iter_files() ))
    if not modinfo:
        print("未找到modinfo", modId)

    for info in modinfo_list:
        modinfo[info.filename[:-4]] = info.read()
        status = 1
        
    data = modinfo['modinfo'].decode(encoding)
    
    return lua_runtime(data)
    

def lua_runtime(data: bytes):
    lua = lupa.LuaRuntime()
    lang='zh'
    # lupa 全局变量
    lupag = ['locale', 'folder_name', 'ChooseTranslationTable'] + ['print', 'rawlen', 'loadfile', 'rawequal', 'pairs', '_VERSION', 'select', 'pcall', 'debug', 'io', 'getmetatable', 'assert', 'package', 'os', 'warn', 'next', 'load', 'tostring', 'setmetatable', 'rawget', 'coroutine', 'tonumber', 'error', 'collectgarbage', 'python', 'utf8', 'math', 'ipairs', 'rawset', 'type', 'xpcall', '_G', 'table', 'dofile', 'require', 'string']

    # 模拟运行环境
    lua.execute(f'locale = "{lang}"')
    lua.execute(f'folder_name = "workshop-{modId}"')
    lua.execute(f'ChooseTranslationTable = function(tbl) return tbl["{lang}"] or tbl[1] end')

    # 开始处理数据
    lua.execute(data)

    # 选取需要的值并转为 python 对象
    g = lua.globals()
    # 选择白名单或黑名单模式
    info_dict = {key: table_dict(g[key]) for key in filter(lambda x: x not in lupag, g)}
    # info_dict = {key: table_dict(g[key]) for key in filter(lambda x: x in info_list_full, g)}

    # 去除空值
    info_dict = {i: j for i, j in info_dict.items() if j or j is False or j == 0}
    
    return info_dict

def table_dict(lua_table):
    if lupa.lua_type(lua_table) == 'table':
        keys = list(lua_table)
        # 假如lupa.table为空，或keys都是整数，且从数字 1 开始以 1 为单位递增，则认为是列表，否则为字典
        if reduce(lambda x, y: x and isinstance(y, int), keys, len(keys) == 0 or keys[0] == 1):  # 为空或首项为 1，全为整数
            if all(map(lambda x, y: x + 1 == y, keys[:-1], keys[1:])):  # 以 1 为单位递增
                return list(map(lambda x: table_dict(x), lua_table.values()))
        return dict(map(lambda x, y: (x, table_dict(y)), keys, lua_table.values()))
    # 由于需要用于解析 modinfo 为 json 格式，所以不支持函数，这里直接删掉
    if lupa.lua_type(lua_table) == 'function':
        return '这里原本是个函数，不过已经被我干掉了'
    return lua_table

modId = 2505341606
mod_info_dict = get_mod_info_dict(modId)

print(mod_info_dict)