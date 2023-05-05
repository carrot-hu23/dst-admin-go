# -*- coding: utf-8 -*-
import os
import zipfile
from shutil import rmtree

import gevent
import steam.client
import steam.client.cdn
import steam.core.cm
import steam.webapi
from steam.exceptions import SteamError

def download_dst_scripts(steamcdn=None):

    anonymous = steam.client.SteamClient()
    anonymous.anonymous_login()
    steamcdn = steam.client.cdn.CDNClient(anonymous)

    print('开始下载世界设置所需饥荒文件')
    # 0:失败, 1:下载成功, 2:没有全部下载成功
    if not os.path.exists('data'):
        os.mkdir('data')
    if not os.path.exists('data/databundles'):
        os.mkdir('data/databundles')
    if not os.path.exists('data/images'):
        os.mkdir('data/images')
    code = 0
    status = {'scripts': False, 'set_tex': False, 'set_xml': False, 'gen_tex': False, 'gen_xml': False, 'ver': False}
    for i in range(5):
        try:
            # b = next(steamcdn.get_manifests(343050, filter_func=lambda x, y: x == 343052))
            # version = next(b.iter_files('version.txt')).read().decode('utf-8').strip()
            for file in steamcdn.iter_files(343050, filter_func=lambda x, y: x == 343052):
                if 'scripts.zip' in file.filename and not status['scripts']:
                    scripts_con = file.read()
                    with open('data/databundles/scripts.zip', 'wb') as f:
                        f.write(scripts_con)
                    with zipfile.ZipFile('data/databundles/scripts.zip') as file_zip:
                        if os.path.exists('data/databundles/scripts'):
                            rmtree('data/databundles/scripts')
                        file_zip.extractall('data/databundles')
                    print('获取 scripts.zip 成功')
                    status['scripts'] = True
                elif 'worldsettings_customization.tex' in file.filename and not status['set_tex']:
                    with open('data/images/worldsettings_customization.tex', 'wb') as f:
                        f.write(file.read())
                    print('获取 worldsettings_customization.tex 成功')
                    status['set_tex'] = True
                elif 'worldsettings_customization.xml' in file.filename and not status['set_xml']:
                    with open('data/images/worldsettings_customization.xml', 'wb') as f:
                        f.write(file.read())
                    print('获取 worldsettings_customization.xml 成功')
                    status['set_xml'] = True
                elif 'worldgen_customization.tex' in file.filename and not status['gen_tex']:
                    with open('data/images/worldgen_customization.tex', 'wb') as f:
                        f.write(file.read())
                    print('获取 worldgen_customization.tex 成功')
                    status['gen_tex'] = True
                elif 'worldgen_customization.xml' in file.filename and not status['gen_xml']:
                    with open('data/images/worldgen_customization.xml', 'wb') as f:
                        f.write(file.read())
                    print('获取 worldgen_customization.xml 成功')
                    status['gen_xml'] = True
                elif 'version.txt' in file.filename and not status['ver']:
                    with open('data/version.txt', 'wb') as f:
                        f.write(file.read())
                    print('获取 version.txt 成功')
                    status['ver'] = True

                code = 1 if all(status.values()) else 2
            if code == 1:
                break
        except SteamError as e:
            print(e, exc_info=True)
        except gevent.timeout.Timeout:
            print('超时了')
    return code, status

if __name__ == "__main__":
    print(download_dst_scripts())