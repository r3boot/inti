
class BaseClass:
    _name = 'BaseClass'

    def __init__(self, output):
        self._output = output

    def info(self, msg):
        self._output.info(self._name, msg)

    def warning(self, msg):
        self._output.info(self._name, msg)

    def error(self, msg):
        self._output.info(self._name, msg)

    def debug(self, msg):
        self._output.info(self._name, msg)
