
#from gevent import monkey; monkey.patch_all()

import bottle
import json
import os

from inti.baseclass import BaseClass
from inti.controller import Controller

class API(BaseClass):
    def __init__(self, output, config, api_queue, controller, policer, host='localhost', port='7231'):
        BaseClass.__init__(self, output)
        self.q = api_queue
        self._cfg = config
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
        self._s.route('/',           method='GET', callback=self.serve_webapp)
        self._s.route('/js/<name>',  method='GET', callback=self.serve_js)
        self._s.route('/css/<name>', method='GET', callback=self.serve_css)
        self._s.route('/img/<name>', method='GET', callback=self.serve_img)
        self._s.route('/config',     method='GET', callback=self.get_config)

        self._s.route('/ping',       method='GET', callback=self.ping)
        self._s.route('/frame',      method='PUT', callback=self.send_frame)
        self._s.route('/off',        method='GET', callback=self.disable_queue)
        self._s.route('/on',         method='GET', callback=self.enable_queue)

    def run(self):
        #self._s.run(host=self._host, port=self._port, server='gevent')
        self._s.run(host=self._host, port=self._port)

    def _validate_frame(self, data):
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
            data['duration'] = 0

        if not isinstance(data['duration'], int):
            self.error('invalid frame data: duration: {0}'.format(
                data['duration']))
            return False

        return data

    def _serve_file(self, filetype, file):
        filename = '{0}/{1}/{2}'.format(
            self._cfg['webapp']['media'],
            filetype,
            file)
        if not os.path.exists(filename):
            self.warning('{0} does not exist'.format(filename))
            return bottle.abort(404, 'File not found')

        return open(filename, 'r').read()

    def serve_js(self, name):
        bottle.response.content_type = 'text/javascript; charset=UTF-8'
        return self._serve_file('js', name)

    def serve_css(self, name):
        bottle.response.content_type = 'text/css; charset=UTF-8'
        return self._serve_file('css', name)

    def serve_img(self, name):
        bottle.response.content_type = 'image/png'
        return self._serve_file('img', name)

    def serve_webapp(self):
        return self._serve_file('html', 'app.html')

    def get_config(self):
        bottle.response.content_type = 'application/json; charset=UTF-8'
        cfg = {
            'spots':  self._cfg['spots'],
            'groups': self._cfg['groups'],
        }
        return json.dumps(cfg)

    def ping(self):
        return 'pong\r\n'

    def send_frame(self):
        srcip = bottle.request.remote_addr
        if self._p.ratelimit(srcip):
            return bottle.abort(503, 'Service unavailable')

        frame_data = json.load(bottle.request.body)
        frame_data = self._validate_frame(frame_data)

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
