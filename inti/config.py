
import time
import yaml

from inti.baseclass import BaseClass
from inti.output    import Output

class Config(BaseClass):
    def __init__(self, output, cfgfile):
        BaseClass.__init__(self, output)
        self._load_config(cfgfile)

    def __getitem__(self, key):
        try:
            return self._cfg[key]
        except KeyError:
            return None

    def _load_config(self, cfgfile):
        raw_config = open(cfgfile, 'r').read()
        cfg = (yaml.load(raw_config))
        cfg['num_spots'] = len(cfg['spots'])
        self._cfg = cfg

if __name__ == '__main__':
    output = Output(debug=True)

    cfg = Config(output, './config/controller.yaml')

    while not output.flushed():
        time.sleep(1)
