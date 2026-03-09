import lupa
from functools import reduce

import steam.client
import steam.client.cdn
import steam.core.cm
import steam.webapi
from steam.exceptions import SteamError

import requests
import json
from io import BytesIO
from urllib.parse import urlencode
from urllib.request import Request, urlopen
import zipfile
import math

from http.server import BaseHTTPRequestHandler, HTTPServer
from urllib.parse import urlparse, parse_qs


anonymous = steam.client.SteamClient()
anonymous.anonymous_login()
steamcdn = steam.client.cdn.CDNClient(anonymous)

encoding = 'utf-8'
modinfo_lua = 'modinfo.lua'
modinfo_chs_lua = 'modinfo_chs.lua'

with open('steamapikey.txt', 'r', encoding='utf-8') as key_value:
    steamapikey = key_value.read()
    print('从文件读取 steamapikey 成功')


def search_mod_list(text='', page=1, num=25): 

    url = 'http://api.steampowered.com/IPublishedFileService/QueryFiles/v1/'
    data = {
        'page': page,
        'key': steamapikey,  # steam apikey  https://steamcommunity.com/dev/apikey
        'appid': 322330,  # 游戏id
        'language': 6,  # 0英文，6简中，7繁中
        'return_tags': True,  # 返回mod详情中的标签
        'numperpage': num,  # 每页结果
        'search_text': text,  # 标题或描述中匹配的文字
        'return_vote_data': True,
        'return_children': True
    }
    url = url + '?' + urlencode(data)
    for _ in range(2):
        try:
            req = Request(url=url)
            response = urlopen(req, timeout=10)
            mod_data = response.read().decode('utf-8')
            # print(mod_data)
            response.close()
            break
        except Exception as e:
            print('搜索mod失败\n', e )
    else:
        return   
    mod_data = json.loads(mod_data).get('response')

    total = mod_data.get('total')
    mod_info_full = mod_data.get('publishedfiledetails')

    mod_list = []
    if mod_info_full:
        for mod_info_raw in mod_info_full:
            img = mod_info_raw.get('preview_url', '')
            vote_data = mod_info_raw.get('vote_data', {})
            auth = mod_info_raw.get("creator", 0),
            mod_info = {
                'id': mod_info_raw.get('publishedfileid', ''),
                'name': mod_info_raw.get('title', ''),
                'author': f'https://steamcommunity.com/profiles/{auth}/?xml=1' if auth else '',
                'desc': mod_info_raw.get('file_description'),
                'time': int(mod_info_raw.get('time_updated', 0)),
                'sub': int(mod_info_raw.get('subscriptions', '0')),
                'img': img if 'steamuserimages' in img else '',
                # 'v': [*[i.get('tag')[8:] for i in mod_info_raw.get('tags', '') if i.get('tag', '').startswith('version:')], ''][0],
                'vote': {'star': int(vote_data.get('score') * 5) + 1 if vote_data.get('votes_up') + vote_data.get('votes_down') >= 25 else 0,
                         'num': vote_data.get('votes_up') + vote_data.get('votes_down')}
            }
            if mod_info_raw.get("num_children"):
                mod_info['child'] = list(map(lambda x: x.get('publishedfileid'), mod_info_raw.get("children")))

            mod_list.append(mod_info)
    return {'page': page, 'size': num,'total': total, 'totalPage': 1 if total/num < 1  else math.ceil(total/num),'data': mod_list}


def get_mod_base_info(modId: int):

    url = 'http://api.steampowered.com/IPublishedFileService/GetDetails/v1/'
    data = {
        'key': steamapikey,  # steam apikey  https://steamcommunity.com/dev/apikey
        'language': 6,  # 0英文，6简中，7繁中
        'publishedfileids[0]': str(modId),  # 要查询的发布文件的ID
    }
    url = url + '?' + urlencode(data)
    payload={}
    headers = {}
    response = requests.request("GET", url, headers=headers, data=payload, verify=False)

    data = json.loads(response.text)['response']['publishedfiledetails'][0]
    if data['result'] != 1:
        print("get mod error")
        return {}
    img = data.get('preview_url', '')
    author = data.get('creator')
    return {
        'id': data.get('publishedfileid'),
        'name': data.get('title'),
        'last_time': data.get('time_updated'),
        "description": data.get("file_description"),
        'auth': f'https://steamcommunity.com/profiles/{author}/?xml=1' if author else '',
        'file_url': data.get('file_url'),
        'img': f'{img}?imw=64&imh=64&ima=fit&impolicy=Letterbox&imcolor=%23000000&letterbox=true'
        if 'steamuserimages' in img else '',
        'v': [*[i.get('tag')[8:] for i in data.get('tags', '') if i.get('tag', '').startswith('version:')], ''][0],
        'creator_appid': data.get('creator_appid'),
        'consumer_appid': data.get('consumer_appid'),
    }

def check_is_dst_mod(mod_info):
    creator_appid, consumer_appid = mod_info.get('creator_appid'), mod_info.get('consumer_appid')
    if not (creator_appid == 245850 and consumer_appid == 322330):
        print('%s 不是饥荒联机版 mod', mod_info.get("id"))
        return False
    return True

def get_mod_info(modId: int):

    # 获取基础信息
    mod_info = get_mod_base_info(modId)

    if not check_is_dst_mod(mod_info):
        return {}
    
    file_url = mod_info['file_url']
    if file_url:
        mod_type = 'v1'
        mod_config = get_mod_config_file_by_url(file_url=file_url)
    else:
        mod_config = get_mod_config_file_by_steamcmd(modId=modId)
        mod_type = 'v2'
    # 获取 mod 文件

    mod_info['mod_type'] = mod_type
    mod_info['mod_config'] = mod_config

    return mod_info

def get_mod_config_file_by_url(file_url: str):
    status, modinfo = 0, {'modinfo': b'', 'modinfo_chs': b''}
    resp_input_io = BytesIO()
    isSuccess = False
    for i in range(3):
        try:
            req = Request(url=file_url)
            res = urlopen(req, timeout=10)
            resp_input_io.write(res.read())
            res.close()
            isSuccess = True
            break
        except Exception as e:
            print("下载失败 \n", e, file_url)
    if not isSuccess:
        return {}
    with zipfile.ZipFile(resp_input_io) as file_zip:
        namelist = file_zip.namelist()
        if modinfo_lua in namelist:
            modinfo['modinfo'] = file_zip.read(modinfo_lua)
        if modinfo_chs_lua in namelist:
            modinfo['modinfo_chs'] = file_zip.read(modinfo_chs_lua)
    
    resp_input_io.close()
    data = modinfo['modinfo'].decode(encoding)
    return lua_runtime(data)

def get_mod_config_file_by_steamcmd(modId: int):
    return get_mod_info_dict(modId=modId)

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
# mod_info_dict = get_mod_info_dict(modId)
# print(mod_info_dict)

#print(get_mod_base_info(modId=modId))
# 1216718131

# print(get_mod_info(modId=modId))

def get_dst_version():
    # 0 下载失败， 1 下载成功， 2 内容为空
    for _ in range(3):
        try:
            # b = next(steamcdn.get_manifests(343050, filter_func=lambda x, y: x == 343052))
            # version = next(b.iter_files('version.txt')).read().decode('utf-8').strip()
            version = [*steamcdn.iter_files(343050, 'version.txt', filter_func=lambda x, y: x == 343052)]
            if not version:
                return 0
            print(version)
            version = version[0].read().decode('utf-8').strip()
            # print(version)
            return version
        except SteamError as e:
            print('获取版本失败', e)
    return 0


class ModHandler(BaseHTTPRequestHandler):
    
    def do_GET(self):
        if self.path.startswith('/py/mod/'):
            self.mod_api()
        elif self.path.startswith('/py/search/mod'):
            self.search_api()
        elif self.path.startswith('/py/version'):
            self.dst_version_api()
        else:
            self.send_error(404)

    def mod_api(self):
        mod_id = self.path.split('/')[3]
        mod_info = get_mod_info(int(mod_id))
        response = {'modInfo': mod_info}
        self.send_response(200)
        self.send_header('Content-type', 'application/json')
        self.end_headers()
        self.wfile.write(json.dumps(response).encode())

    def search_api(self):
        
        query = parse_qs(urlparse(self.path).query)
        search_text = query['text'][0]
        page = query['page'][0]
        size = query['size'][0]

        print(search_text, page, size)

        response = search_mod_list(search_text, int(page), int(size))

        self.send_response(200)
        self.send_header('Content-type', 'application/json')
        self.end_headers()
        self.wfile.write(json.dumps(response).encode())
    
    def dst_version_api(self):
        response = get_dst_version()
        self.send_response(200)
        self.send_header('Content-type', 'application/json')
        self.end_headers()
        self.wfile.write(json.dumps(response).encode())

PORT = 8000
httpd = HTTPServer(('localhost', PORT), ModHandler)
print(f'Serving on port {PORT}')
httpd.serve_forever()
