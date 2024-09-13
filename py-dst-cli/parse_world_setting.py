##!/usr/bin/python3
# -*- coding: utf-8 -*-

import os
from functools import reduce
from os.path import join as pjoin
from re import compile, findall, search, sub

import lupa


def table_dict(lua_table):
    if lupa.lua_type(lua_table) == 'table':
        keys = list(lua_table)
        # 假如lupa.table为空，或keys都是整数，且从数字 1 开始以 1 为单位递增，则认为是列表，否则为字典
        if reduce(lambda x, y: x and isinstance(y, int), keys, len(keys) == 0 or keys[0] == 1):  # 为空或首项为 1，全为整数
            if all(map(lambda x, y: x + 1 == y, keys[:-1], keys[1:])):  # 以 1 为单位递增
                return list(map(lambda x: table_dict(x), lua_table.values()))
        new_dict = dict(map(lambda x, y: (x, table_dict(y)), keys, lua_table.values()))
        if 'desc' in new_dict:  # task_set 和 start_location 的 desc 是个函数，需要调用一下返回实际值
            for i, j in new_dict.items():
                if lupa.lua_type(j) == 'function':
                    new_dict[i] = {world: table_dict(j(world)) for world in new_dict.get('world', [])}
        return new_dict
    return lua_table


def dict_table(py_dict, lua_temp):  # dict 转 table。列表之类类型的转过去会有索引，table_from 的问题
    if isinstance(py_dict, dict):
        return lua_temp.table_from(
            {i: (dict_table(j, lua_temp) if isinstance(j, (dict, list, tuple, set)) else j) for i, j in py_dict.items()})
    if isinstance(py_dict, (list, tuple, set)):
        return lua_temp.table_from([(dict_table(i, lua_temp) if isinstance(i, (dict, list, tuple, set)) else i) for i in py_dict])
    return py_dict


def scan(dict_scan, num, key_set):  # 返回指定深度的 keys 集合, key_set初始传入空set
    if num != 0:
        for value in dict_scan.values():
            if isinstance(value, dict):
                key_set = key_set | scan(value, num - 1, key_set)
        return key_set
    return key_set | set(dict_scan)


def parse_po(path_po):  # 把 .po 文件按照 msgctxt: msgstr 的格式转为字典，再以 . 的深度分割 keys。这里为了效率主要转了 UI 部分的
    print('开始通过 .po 文件 获取翻译')
    with open(path_po, 'rb') as f:
        f.seek(-50000, 2)
        data = f.read()
        while b'"STRINGS.T' not in data:
            f.seek(-100000, 1)
            data = f.read(50000) + data
    data = data.decode('utf-8').replace('\r\n', '\n')
    pattern = compile(r'\nmsgctxt\s*"(.*)"\nmsgid\s*"(.*)"\nmsgstr\s*"(.*)"')

    dict_zh_split, dict_en_split = {}, {}

    print('获取中文翻译')
    dict_zh = {i[0]: i[2] for i in pattern.findall(data)}  # 因为 costomize 中有连接字符串的，所以这里不能构建成一个字典，会出错
    for i, j in dict_zh.items():
        split_key(dict_zh_split, i.split("."), j)

    # print('获取英文对照')
    # dict_en = {i[0]: i[1] for i in pattern.findall(data)}  # 因为 costomize 中有连接字符串的，所以这里不能构建成一个字典，会出错
    # for i, j in dict_en.items():
    #     split_key(dict_en_split, i.split("."), j)

    dict_split = {'zh': dict_zh_split}
    if dict_en_split:
        dict_split['en'] = dict_en_split
    return dict_split


def split_key(dict_split, list_split, value):  # 以列表值为 keys 补全字典深度。用于分割 dict 的 keys，所以叫 split
    if not list_split:
        return
    if list_split[0] not in dict_split:
        dict_split[list_split[0]] = value if len(list_split) == 1 else {}
    split_key(dict_split.get(list_split.pop(0)), list_split, value)


def creat_newdata(path_cus, new_cus):  # 删去local、不必要的require 和不需要的内容
    with open(path_cus + '.lua', 'r') as f:
        data = f.read()
    if 'local MOD_WORLDSETTINGS_GROUP' in data:
        data = data[:data.find('local MOD_WORLDSETTINGS_GROUP')]
    data = sub(r'local [^=]+?\n', '', data).replace('local ', '')
    data = sub(r'require(?![^\n]+?(?=tasksets"|startlocations"))', '', data)
    with open(new_cus + '.lua', 'w+') as f:
        f.write(data)


def parse_cus(lua_cus, po):
    print('准备解析 customize.lua 文件')
    new_cus = lua_cus + '_tmp'
    creat_newdata(lua_cus, new_cus)  # 删去多余的不需要的数据并另存

    print('准备运行环境')
    lua = lupa.LuaRuntime()
    lua.execute('function IsNotConsole() return true end')  # IsNotConsole() 不是 PS4 或 XBONE 就返回 True  # for customize
    lua.execute('function IsConsole() return false end')  # IsConsole() 是 PS4 或 XBONE 就返回 True
    lua.execute('function IsPS4() return false end')  # IsPS4() 不是 PS4 就返回False  # for customize
    lua.execute('ModManager = {}')  # for startlocations
    lua.require('class')  # for util
    lua.require('util')  # for startlocations
    lua.require('constants')  # 新年活动相关
    lua.require("strict")

    dict_po = parse_po(po)
    options_list = ['WORLDGEN_GROUP', 'WORLDSETTINGS_GROUP']  # 所需数据列表
    misc_list = ['WORLDGEN_MISC', 'WORLDSETTINGS_MISC']  # 所需数据列表
    options = {}

    print('解析 customize.lua 文件，语言：%s', ', '.join(dict_po.keys()))
    # 获取英文版的应该可以通过导入strings来做？应该不需要通过 .po 文件匹配
    for lang, tran in dict_po.items():
        strings = dict_table(tran.get('STRINGS'), lua)
        if strings:
            pass
        lua.execute('STRINGS=python.eval("strings")')  # 为了翻译，也免去要先给 STRINGS 加引号之类的麻烦事
        # lua.execute('POT_GENERATION = true')
        # lua.require('strings')
        lua.require(new_cus)  # 终于开始干正事了。导入的 tasksets 会自动打印一些东西出来
        options[lang] = {'setting': {i: table_dict(lua.globals()[i]) for i in options_list if i in lua.globals()},
                         'translate': tran}
        for package in list(lua.globals().package.loaded):  # 清除加载的 customize 模块，避免下次 require 时不加载
            if 'map/' in package:
                lua.execute(f'package.loaded["{package}"]=nil')  # table.remove 不能用，显示 package.loaded 长度为0
    miscs = {i: table_dict(lua.globals()[i]) for i in misc_list if i in lua.globals()}

    print('解析 customize.lua 文件完毕')
    return options, miscs


def parse_option(group_dict, path_base):
    print('重新组织设置格式')
    print('这里写的太什么了，不好插日志')

    result = {}
    img_info = {}
    img_name = ''
    for lang, opt in group_dict.items():
        setting, translate = opt.values()
        result[lang] = {'forest': {}, 'cave': {}}
        for group, group_value in setting.items():
            for world_type in result[lang].values():
                world_type[group] = {}
            for com, com_value in group_value.items():
                desc_val = com_value.get('desc')
                if desc_val:
                    desc_val = {i['data']: i['text'] for i in desc_val}
                for world, world_value in result.get(lang).items():
                    img_name = com_value.get('atlas', '').replace('images/', '').replace('.xml', '')
                    if img_name not in img_info:
                        with open(pjoin(path_base, com_value.get('atlas')), 'r', encoding='utf-8') as f:
                            data = f.read()
                        image_filename = search('filename="([^"]+)"', data).group(1)
                        with open(pjoin(path_base, 'images', image_filename), 'rb') as f:
                            img_data = f.read(96)
                        image_width, image_height = int(img_data[88:90].hex(), 16), int(img_data[90:92].hex(), 16)
                        img_width_start, img_width_end = search(r'u1="([^"]+?)"\s*?u2="([^"]+?)"', data).groups()
                        img_item_width = int(image_width / round(1 / (float(img_width_end) - float(img_width_start))))
                        item_num_w, item_num_h = image_width / img_item_width, image_height / img_item_width
                        img_pos = {i[0]: {'x': round(float(i[1]) * item_num_w) / item_num_w,
                                          'y': 1 - round(float(i[2]) * item_num_h) / item_num_h} for i in
                                   findall(r'<Element\s+name="([^"]+?)"\s*u1="([^"]+?)"[\d\D]*?v2="([^"]+?)"', data)}
                        img_info[img_name] = {'img_items': img_pos, 'width': image_width, 'height': image_height,
                                              'item_size': img_item_width}
                    world_value.get(group)[com] = {
                        'order': int(com_value.get('order', 0)),
                        'text': com_value.get('text', ''),
                        'atlas': {'name': img_name, 'width': img_info[img_name]['width'],
                                  'height': img_info[img_name]['height'], 'item_size': img_info[img_name]['item_size']},
                        'desc': desc_val,
                        'items': {}}
                for item, item_value in com_value['items'].items():
                    tmp = []
                    if 'forest' in item_value.get('world', '') or not item_value.get('world'):
                        tmp.append(('forest', result[lang]['forest']))
                    if 'cave' in item_value.get('world', ''):
                        tmp.append(('cave', result[lang]['cave']))
                    print('这个有问题{}\n'.format(item_value) if not tmp else '', end='')
                    for world, world_value in tmp:
                        items = world_value[group][com]['items']
                        items[item] = {
                            'image': img_info[img_name]['img_items'][item_value['image']],
                            'text': translate['STRINGS']['UI']['CUSTOMIZATIONSCREEN'].get(item.upper())}
                        if item_value.get('desc'):
                            item_desc = item_value['desc']
                            item_desc = item_desc.get(world) if isinstance(item_desc, dict) else item_desc
                            item_desc = {i['data']: i['text'] for i in item_desc}
                            items[item]['desc'] = item_desc

                        # 为带有排序优先属性 order 的项目添加 order
                        if item_value.get('order'):
                            items[item]['order'] = item_value['order']
                        # 修正地上地下使用不同 desc 时，共用的 value 不在某个的 desc 内的情况
                        tmp_desc = items[item].get('desc') or world_value[group][com]['desc']
                        tmp_value = item_value.get('value')
                        items[item]['value'] = item_value.get('value') if tmp_value in tmp_desc else list(tmp_desc)[0]

    print('格式重组完成')

    # 清理空的 items 项。并打印不同世界的项目数。
    tip_times = 0
    for lang_value in result.values():
        for world_name, world_value in lang_value.items():
            setting_num = 0
            for groups_value in world_value.values():
                for group_name, group_value in list(groups_value.items())[:]:
                    setting_num += len(group_value['items'])
                    if not group_value['items']:
                        del groups_value[group_name]
            if tip_times < len(lang_value):
                print(f'{world_name} 拥有 {setting_num} 个可设置项')
                tip_times += 1

    return result


def parse_world_setting(path_base="data"):

    print('开始解析世界设置')

    # path_base 指向饥荒程序文件夹中的 data 文件夹路径
    # path_base = "data"
    path_script = pjoin(path_base, 'databundles', 'scripts')
    po_chs = 'languages/chinese_s.po'
    lua_customize = 'map/customize'  # 务必用正斜杠避免问题。lua 内部 require 会用正斜杠，两个不一样的话操作对应模块时会有坑

    # if not os.path.exists(path_cus := pjoin(path_script, lua_customize + '.lua')):
    #     print('未找到 %s 文件，准备退出', path_cus)
    #     return
    # if not os.path.exists(path_po := pjoin(path_script, po_chs)):
    #     print('未找到 %s 文件，准备退出', path_po)
    #     return
    path_cus = pjoin(path_script, lua_customize + '.lua')
    if not os.path.exists(path_cus):
        print('未找到 %s 文件，准备退出' % path_cus)
        return

    path_po = pjoin(path_script, po_chs)
    if not os.path.exists(path_po):
        print('未找到 %s 文件，准备退出' % path_po)
        return

    cwd_now = os.getcwd()
    os.chdir(path_script)
    settings = None

    try:
        options_raw, misc = parse_cus(lua_customize, po_chs)
        # print(options_raw)
        os.chdir(cwd_now)
        settings = parse_option(options_raw, path_base)
        print('解析世界设置完成')
    except Exception as e:
        os.chdir(cwd_now)
        print('解析世界设置失败')
        print(e, exc_info=True)

    return settings


if __name__ == "__main__":
    datadata = parse_world_setting("C:\\Program Files (x86)\\Steam\steamapps\\common\\Don't Starve Together\\data")
    print(datadata)
    import json
    with open('C:\\Users\\paratera\\Desktop\\我的\\饥荒面板\\dst-admin-go\\py-dst-cli\\dst_world_setting.json', 'w', encoding='utf-8') as f:
        f.write(json.dumps(datadata, ensure_ascii=False))