
#from gevent import monkey; monkey.patch_all()

import bottle
import json

from inti.baseclass import BaseClass
from inti.controller import Controller

class API(BaseClass):
    def __init__(self, output, api_queue, controller, policer, host='localhost', port='7231'):
        BaseClass.__init__(self, output)
        self.q = api_queue
        self._c = controller
        self._p = policer
        self._host = host
        self._port = port
        self._s = bottle.Bottle()
        self._setup_routing()
        self.run()

    def __destroy__(self):
        pass

    def _setup_routing(self):
        self._s.route('/ping',  method='GET', callback=self.ping)
        self._s.route('/frame', method='PUT', callback=self.send_frame)
        self._s.route('/off',   method='GET', callback=self.disable_queue)
        self._s.route('/on',   method='GET', callback=self.enable_queue)

    def run(self):
        #self._s.run(host=self._host, port=self._port, server='gevent')
        self._s.run(host=self._host, port=self._port)

    def fetch_frame(self):
        if self.q.empty():
            return False
        return self.q.get()

    def validate_frame(self, data):
        if not 'frame' in data.keys():
            self.error('invalid frame data: "frame" missing')
            return False

        i = 0
        for value in data['frame']:
            if not isinstance(value, int):
                self.error('invalid frame data: frame[{0}]: {1}'.format(
                    i, value))
                return False
            i += 1

        if not 'duration' in data.keys():
            self.warning('no duration found, setting to 0ms')
            data['duration'] = 0

        if not isinstance(data['duration'], int):
            self.error('invalid frame data: duration: {0}'.format(
                data['duration']))
            return False

        return data

    def ping(self):
        return 'pong\r\n'

    def send_frame(self):
        srcip = bottle.request.remote_addr
        if self._p.ratelimit(srcip):
            return bottle.abort(503, 'Service unavailable')

        frame_data = json.load(bottle.request.body)
        frame_data = self.validate_frame(frame_data)

        if not frame_data:
            return bottle.abort(404, 'Not found')

        self.q.put([srcip, frame_data['frame'], frame_data['duration']])

    def disable_queue(self):
        self._c.blackout(True)

    def enable_queue(self):
        self._c.blackout(False)

"""
if __name__ == '__main__':
    from inti.output import Output
    o = Output(debug=True)
    c = Controller(o, num_spots=5)
    api = API(o, c)
    api.run()
"""
