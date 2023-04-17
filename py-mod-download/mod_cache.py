
import threading
import time


class Entity:
    def __init__(self):
        self.ready = threading.Semaphore(0)
        self.data = None
        self.version = None

class ModCache:
    lock = None
    cache = None

    def __init__(self, f):
        self.lock = threading.Lock()
        self.cache = {}
        self.f = f
    def get(self, modId:str, version:str):
        key = modId+'@'+version
        self.lock.acquire()
        if key in self.cache.keys():
            mod_info =  self.cache[key]
        else:
            mod_info = None
        if mod_info is None:
            mod_info = Entity()
            self.cache[key] = mod_info
            self.lock.release()

            #TODO
            time.sleep(1)
            mod_info.data, mod_info.version = key, version
            mod_info.ready.release(1)

        else:
            self.lock.release()
            
            if mod_info.data is None:
                self.lock.acquire()
                if mod_info.data is None: 
                    print(key, "is waiting")
                    mod_info.ready.acquire(1)
                self.lock.release()

        print(mod_info.data, mod_info.version)
        return mod_info.data, mod_info.version
    

modcache = ModCache(None)

t1 = threading.Thread(target=modcache.get, args=("111", "1.1.0"))
t2 = threading.Thread(target=modcache.get, args=("111", "1.1.0"))
t3 = threading.Thread(target=modcache.get, args=("111", "1.1.2"))


t4 = threading.Thread(target=modcache.get, args=("444", "4.0.1"))
t5 = threading.Thread(target=modcache.get, args=("444", "4.0.2"))
t6 = threading.Thread(target=modcache.get, args=("666", "6.0.0"))

t1.start()
t2.start()
t3.start()
t4.start()
t5.start()
t6.start()

t1.join()
t2.join()
t3.join()
t4.join()
t5.join()
t6.join()