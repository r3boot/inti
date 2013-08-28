
import time
import yaml

from inti.baseclass import BaseClass
from inti.output    import Output

class Config(BaseClass):
    _name = 'Config'

    def __init__(self, output, cfgfile):
        BaseClass.__init__(self, output)
        self._cfg = self._load_config(cfgfile)

    def __getitem__(self, key):
        try:
            return self._cfg[key]
        except KeyError:
            return None

    def _load_config(self, cfgfile):
        self.debug('loading {0}'.format(cfgfile))
        raw_config = open(cfgfile, 'r').read()
        cfg = (yaml.load(raw_config))

        cfg['num_spots'] = len(cfg['spots'])
        cfg['spot_keys'] = cfg['spots'].keys()
        cfg['spot_keys'].sort()

        if not 'groups' in cfg.keys():
            cfg['groups'] = {}

        cfg['groups']['all'] = {
            'description': 'All slots on the controller',
            'location': 'Global',
            'spots': cfg['spot_keys']
        }
        cfg['num_groups'] = len(cfg['groups'])
        cfg['group_keys'] = cfg['groups'].keys()
        cfg['group_keys'].sort()

        for spot in cfg['spot_keys']:
            for key, value in cfg['spots'][spot].items():
                self.debug('spot[{0}] {1}: {2}'.format(
                    spot, key, value))

        for group in cfg['group_keys']:
            for key, value in cfg['groups'][group].items():
                self.debug('group[{0}] {1}: {2}'.format(
                    group, key, value))

        return cfg

    def dump_config(self):
        return self._cfg

"""
if __name__ == '__main__':
    output = Output(debug=True)

    cfg = Config(output, './config/controller.yaml')

    while not output.flushed():
        time.sleep(1)
"""
