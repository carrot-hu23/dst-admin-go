import json
from http.server import BaseHTTPRequestHandler, HTTPServer

import dst_version
import parse_world_setting
import parse_world_webp

import sys
import os
import sched
import time

world_gen_path = sys.argv[1]

# 初始化调度器
scheduler = sched.scheduler(time.time, time.sleep)

def gen_world_iamge_job(path):
    print("定时任务--生成世界misc文件正在执行, 参数为:", path)
    parse_world_webp.download_dst_scripts()

def gen_world_setting_job(path):
    print("定时任务--生成世界json文件正在执行, 参数为:", path)
    datadata = parse_world_setting.parse_world_setting()
    import json
    if not os.path.exists('data'):
        os.mkdir('data')
    with open('data/dst_world_setting.json', 'w', encoding='utf-8') as f:
        f.write(json.dumps(datadata, ensure_ascii=False))

def run():

    scheduler.enter(10, 1, gen_world_iamge_job, (world_gen_path,))

    scheduler.enter(20, 1, gen_world_setting_job, (world_gen_path,))

    scheduler.run()

# 循环执行任务
while True:
    run()



'''
class DstWorldHandler(BaseHTTPRequestHandler):
    
    def do_GET(self):
        if self.path.startswith('/py/dst/version/'):
            self.getDstVersion()
        elif self.path.startswith('/py/dst/world/setting/webp'):
            self.getDstWorldSettingWebp()
        elif self.path.startswith('/py/dst/world/setting/json'):
            self.getDstWorldSettingJson()
        else:
            self.send_error(404)

    def getDstVersion(self):
        version = dst_version.get_dst_version()
        self.send_response(200)
        self.send_header('Content-type', 'application/json')
        self.end_headers()
        self.wfile.write(json.dumps(version).encode())

    def getDstWorldSettingWebp(self):
        parse_world_webp.download_dst_scripts()
        self.send_response(200)
        self.send_header('Content-type', 'application/json')
        self.end_headers()
        self.wfile.write(json.dumps("ok").encode())
    
    def getDstWorldSettingJson(self):
        response = parse_world_setting.parse_world_setting()
        self.send_response(200)
        self.send_header('Content-type', 'application/json')
        self.end_headers()
        self.wfile.write(json.dumps(response).encode())

PORT = 8000
httpd = HTTPServer(('127.0.0.1', PORT), DstWorldHandler)
print(f'Serving on port {PORT}')
httpd.serve_forever()
'''
