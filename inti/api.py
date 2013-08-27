
from gevent import monkey; monkey.patch_all()

import bottle
import json

from inti.baseclass import BaseClass
from inti.controller import Controller

class API(BaseClass):
    def __init__(self, output, controller, host='localhost', port='7231'):
        BaseClass.__init__(self, output)
        self._c = controller
        self._host = host
        self._port = port
        self._s = bottle.Bottle()
        self._setup_routing()

    def __destroy__(self):
        pass

    def _setup_routing(self):
        self._s.route('/ping',  method='GET', callback=self.ping)
        self._s.route('/frame', method='PUT', callback=self.send_frame)

    def run(self):
        self._s.run(host=self._host, port=self._port, server='gevent')

    def validate_frame(self, data):
        if not 'frame' in data.keys():
            self.error('invalid frame data: "frame" missing')
            return False

        if not 'duration' in data.keys():
            self.warning('no duration found, setting to 0ms')
            data['duration'] = 0

        #new_frame = [^chr(x) for x in data['frame']]
        #data['frame'] = new_frame

        return data

    def ping(self):
        return 'pong\r\n'

    def send_frame(self):
        frame_data = json.load(bottle.request.body)
        frame_data = self.validate_frame(frame_data)

        if not frame_data:
            return bottle.abort(404, 'Not found')

        self._c.queue_frame(frame_data['frame'], frame_data['duration'])

if __name__ == '__main__':
    from inti.output import Output
    o = Output(debug=True)
    c = Controller(o, num_spots=5)
    api = API(o, c)
    api.run()
