
import os
import time

from inti.baseclass import BaseClass

class Controller(BaseClass):
    _name = 'Controller'

    def __init__(self, output, dmx_port='/dev/dmx0'):
        BaseClass.__init__(self, output)
        self._fd = self._setup_fd(dmx_port)
        self.debug('class initialized')

    def _setup_fd(self, dmx_port):
        if not os.path.exists(dmx_port):
            self.error('{0} does not exist'.format(dmx_port))
            return

        fd = os.open(dmx_port, os.O_WRONLY)
        return fd

    def __destroy__(self):
        if self._fd:
            os.close(self._fd)

    def send_frame(self, data):
        data = [0] + data
        if not self._fd:
            self.error('failed to write to controller')
            return

        self.debug('frame: {0}'.format(
            ' '.join('{0:02x}'.format(x) for x in data[1:])
            ))
        os.write(self._fd, ''.join(chr(x) for x in data))

if __name__ == '__main__':
    from output import Output

    o = Output(debug=True)

    controller = Controller(o)

    for i in xrange(255):
        controller.send_frame([i,0,0]*5)

    for i in reversed(xrange(255)):
        controller.send_frame([i,0,0]*5)

    for i in xrange(255):
        controller.send_frame([0,i,0]*5)

    for i in reversed(xrange(255)):
        controller.send_frame([0,i,0]*5)

    for i in xrange(255):
        controller.send_frame([0,0,i]*5)

    for i in reversed(xrange(255)):
        controller.send_frame([0,0,i]*5)


    while not o.flushed():
        time.sleep(0.1)
