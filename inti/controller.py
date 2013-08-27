
import os
import threading
import time

from inti.baseclass import BaseClass

class Controller(BaseClass, threading.Thread):
    _name = 'Controller'

    def __init__(self, output, queue, policer, dmx_port='/dev/dmx0', num_spots=10):
        threading.Thread.__init__(self)
        self.setDaemon(True)
        self.stop = False
        BaseClass.__init__(self, output)
        self._fd = self._setup_fd(dmx_port)
        self._spots = self._setup_spots(num_spots)
        self._q = queue
        self._p = policer
        self._is_black = False
        self.debug('class initialized')
        self.start()

    def __destroy__(self):
        if self._fd:
            os.close(self._fd)

    def run(self):
        while not self.stop:
            frame_data = self._q.get()
            while self._is_black:
                time.sleep(0.1)
            self.send_frame(*frame_data)

    def _setup_fd(self, dmx_port):
        if not os.path.exists(dmx_port):
            self.error('{0} does not exist'.format(dmx_port))
            return

        fd = os.open(dmx_port, os.O_WRONLY)
        return fd

    def _setup_spots(self, num_spots):
        spots = {}
        for spot_id in xrange(num_spots):
            spots[spot_id] = [0,0,0]
        return spots

    def flushed(self):
        return self._q.empty()

    def blackout(self, value):
        if value != self._is_black and value == True:
            self.info('enabling blackout')
            self.send_frame(False, [0] * (len(self._spots.keys() * 3)))
        else:
            self.info('disabling blackout')
        self._is_black = value

    def send_frame(self, srcip, data, duration=0):
        data = [0] + data
        if not self._fd:
            self.error('failed to write to controller')
            return

        self.debug('frame: {0}'.format(
            ' '.join('{0:02x}'.format(x) for x in data[1:])
            ))
        os.write(self._fd, ''.join(chr(x) for x in data))

        if duration != 0:
            time.sleep(duration/1000.0)

        if srcip:
            self._p.decrement(srcip)

"""
if __name__ == '__main__':
    from output import Output

    o = Output(debug=True)

    controller = Controller(o, num_spots=5)

    controller.blackout()
    controller.queue_frame([255,0,0]*5, duration=1000)
    controller.queue_frame([0,255,0]*5, duration=1000)
    controller.queue_frame([0,0,255]*5, duration=1000)
    controller.blackout()

    while not controller.flushed():
        time.sleep(0.1)

    while not o.flushed():
        time.sleep(0.1)
"""
