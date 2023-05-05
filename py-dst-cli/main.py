import json
from http.server import BaseHTTPRequestHandler, HTTPServer

import dst_version
import parse_world_setting
import parse_world_webp

class ModHandler(BaseHTTPRequestHandler):
    
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
httpd = HTTPServer(('127.0.0.1', PORT), ModHandler)
print(f'Serving on port {PORT}')
httpd.serve_forever()

