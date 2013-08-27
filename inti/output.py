
import Queue
import datetime
import threading
import time

MSG_HEADERS = {
    'info':    'I',
    'warning': 'W',
    'error':   'E',
    'debug':   'D'
}

class Output(threading.Thread):
    def __init__(self, debug=False):
        threading.Thread.__init__(self)
        self.setDaemon(True)
        self.stop = False
        self._debug = debug
        self._q = Queue.Queue()
        self.debug('Output', 'Class initialized')
        self.start()

    def run(self):
        self.debug('Output', 'mainloop started')
        while not self.stop:
            msg_data = self._q.get()
            msg = '{0} {1}[{2}] {3}'.format(
                    msg_data[0].isoformat(' '),
                    msg_data[2],
                    MSG_HEADERS[msg_data[1]],
                    msg_data[3])
            if msg_data[1] == 'debug' and not self._debug:
                continue

            print(msg)

    def _timestamp(self):
        return datetime.datetime.now().isoformat(' ')

    def _msg(self, msgtype, caller, msg):
        self._q.put([datetime.datetime.now(), msgtype, caller, msg])

    def flushed(self):
        return self._q.empty()

    def info(self, caller, msg):
        self._msg('info', caller, msg)

    def warning(self, caller, msg):
        self._msg('warning', caller, msg)

    def error(self, caller, msg):
        self._msg('error', caller, msg)

    def debug(self, caller, msg):
        self._msg('debug', caller, msg)

"""
if __name__ == '__main__':
    output = Output(debug=True)
    output.start()

    output.info('__main__', 'Informational message')
    output.warning('__main__', 'Warning message')
    output.error('__main__', 'Error message')
    output.debug('__main__', 'Debug message')

    while not output.flushed():
        time.sleep(0.1)
"""
