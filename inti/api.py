"""
.. module: api
   :platform: Linux
   :synopsis: Provides a restful api for the Inti DMX controller

.. moduleauthor: Lex van Roon <r3boot@r3blog.nl>
"""
import json
import socket
import sys

# Handle external dependencies
try:
    import bottle
except ImportError:
    print('bottle not found, please run "pip install bottle"')
    sys.exit(1)


# Path towards the media which is to be served
MEDIA = './media'


class RestAPI:
    """Class representing the API used to talk to the Inti DMX setup

    :param logger:  Object representing a configured python logger
    :type  logger:  logging.Logger
    :param host:    Ip-address or Host to listen on
    :type  host:    str
    :param port:    Port to listen on
    :type  port:    int
    """
    def __init__(self, logger, host, port, busses):
        self.log = logger
        self._host = host
        self._port = port
        self._dmx = busses

        self._app = bottle.Bottle()
        self._app.route('/', method='get', callback=self.serve_index)
        self._app.route('/css/<file>', method='get', callback=self.serve_css)
        self._app.route('/js/<file>', method='get', callback=self.serve_js)
        self._app.route('/v1/config', method='get', callback=self.serve_config)
        self._app.route('/v1/buffer', method='post', callback=self.set_buffer)

    def run(self):
        """Start the actual API
        """
        try:
            self._app.run(host=self._host, port=self._port, fast=True)
        except socket.error as err:
            self.log.error('Failed to start API: {0}'.format(err))

    @staticmethod
    def serve_index():
        """Helper function which serves the main index.html for this api
        """
        path = MEDIA + '/html'
        return bottle.static_file('index.html', root=path)

    @staticmethod
    def serve_css(file):
        """Helper function which returns a CSS file for this api if it exists
        """
        path = MEDIA + '/css'
        return bottle.static_file(file, root=path, mimetype='text/css')

    @staticmethod
    def serve_js(file):
        """Helper function which returns a JS file for this api if it exists
        """
        path = MEDIA + '/js'
        mime = 'application/javascript'
        return bottle.static_file(file, root=path, mimetype=mime)

    def serve_config(self):
        """Returns the configuration served by this api
        """
        cfg = {}
        i = 0
        for name, bus in self._dmx.items():
            cfg[name] = {
                'id': i,
                'name': bus['name'],
                'port': bus['port'],
                'buffer': bus['buffer'],
                'fixtures': {}
            }
            for f_name in bus['fixtures'].keys():
                fixture = {
                    'bus': name,
                    'address': bus['fixtures'][f_name].address,
                    'channels': bus['fixtures'][f_name].channels,
                }
                cfg[name]['fixtures'][f_name] = fixture
            i += 1
        return json.dumps(cfg)

    def set_buffer(self):
        """Dumps a buffer onto a DMX Bus
        """
        raw_data = bottle.request.body.read().decode('utf-8')
        bus_cfg = json.loads(raw_data)
        print(bus_cfg['name'])

        bus = self._dmx[bus_cfg['port']]
        if not bus:
            return bottle.HTTPResponse(status=404)
        bus.transfer(bus_cfg['buffer'])
