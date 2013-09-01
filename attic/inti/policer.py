
import threading
import time

from inti.baseclass import BaseClass

NUM_REQUESTS = 0
RATELIMITED = 1
PENALTY = 60
GLOBAL_PENALTY = 120

class Policer(BaseClass, threading.Thread):
    def __init__(self, output, max_requests=3600, max_per_ip=1800):
        BaseClass.__init__(self, output)
        threading.Thread.__init__(self)
        self.setDaemon(True)
        self._event = threading.Event()
        self._event.set()
        self.stop = False
        self._total_requests = 0
        self._global_ratelimit = False
        self._max_requests = max_requests
        self._max_per_ip = max_per_ip
        self._p = {}
        self.start()

    def run(self):
        while not self.stop:
            now = time.time()
            self._event.clear()
            self._cleanup_ratelimits(now)
            self._event.set()
            time.sleep(1.0)

    def _cleanup_ratelimits(self, now):
        for srcip in self._p.keys():
            if not self._p[srcip][RATELIMITED]:
                continue
            elif now - self._p[srcip][RATELIMITED] > PENALTY:
                self.info('clearing ratelimit for {0}'.format(srcip))
                self._p[srcip][NUM_REQUESTS] = 0
                self._p[srcip][RATELIMITED] = False

        if self._global_ratelimit and now - self._global_ratelimit > GLOBAL_PENALTY:
            self.info('global ratelimit disabled')
            self._total_requests = 0
            self._global_ratelimit = False

    def decrement(self, srcip):
        self._event.wait()
        if srcip in self._p.keys():
            self._p[srcip][NUM_REQUESTS] -= 1
            self._total_requests -= 1

    def ratelimit(self, srcip):
        self._event.wait()
        now = time.time()
        if self._global_ratelimit:
            return True

        if self._total_requests > self._max_requests:
            self._global_ratelimit = now
            self.warning('global ratelimit enabled')
            return True

        if srcip not in self._p.keys():
            self._p[srcip] = [0, False]

        if self._p[srcip][RATELIMITED]:
            return True

        if self._p[srcip][NUM_REQUESTS] > (self._max_per_ip-1):
            self.warning('ratelimiting {0}'.format(srcip))
            self._p[srcip][RATELIMITED] = now
            return True

        self._p[srcip][NUM_REQUESTS] += 1
        self._total_requests += 1
