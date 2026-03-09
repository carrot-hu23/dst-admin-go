import steam.client
import steam.client.cdn
import steam.core.cm
import steam.webapi
from steam.exceptions import SteamError


def get_dst_version(steamcdn=None):

    anonymous = steam.client.SteamClient()
    anonymous.anonymous_login()
    steamcdn = steam.client.cdn.CDNClient(anonymous)

    print("开始获取dst版本文件")

    # 0 下载失败， 1 下载成功， 2 内容为空
    for _ in range(3):
        try:
            # b = next(steamcdn.get_manifests(343050, filter_func=lambda x, y: x == 343052))
            # version = next(b.iter_files('version.txt')).read().decode('utf-8').strip()
            version = [*steamcdn.iter_files(343050, 'version.txt', filter_func=lambda x, y: x == 343052)]
            if not version:
                return 0
            version = version[0].read().decode('utf-8').strip()
            # print(version)
            return int(version)
        except SteamError as e:
            print("get dst version error", SteamError)
            return 0
    return 0

if __name__ == "__main__":
    print(get_dst_version())