"""
.. module:: config
   :platform: Linux
   :synopsis: Abstraction layer around the inti configuration file

.. moduleauthor:: Lex van Roon <r3boot@r3blog.nl>
"""
import os
import sys

# Handle external dependencies
try:
    import yaml
except ImportError:
    print('yaml not found, please run "pip install pyyaml"')
    sys.exit(1)


class Config:
    """Class representing the inti configuration file

    :param logger:      Object representing a configured python logger
    :type  logger:      logger.Logger
    :param cfg_file:    Path pointing towards the yaml configuration file
    :type  cfg_file:    str
    """
    def __init__(self, logger, cfg_file):
        self.log = logger
        self._cfg_file = cfg_file
        self._cfg = {}
        self.load_config()

    def __getitem__(self, key):
        """Retrieves a key from the configuration dictionary

        :param key:     Which key to lookup
        :type key:      str
        :returns:       Item from config pointed to by key or None
        :rtype:         dict, None
        """
        try:
            return self._cfg[key]
        except KeyError:
            return None

    def __contains__(self, key):
        """Check if key is a member of the configuration dictionary

        :param key:     Key to lookup
        :type  key:     str
        :returns:       Flag indicating existence of key
        :rtype:         bool
        """
        return key in self._cfg

    def validate_config(self, cfg):
        """Helper function to validate a configuration dictionary for
        correctness

        :param cfg:     Configuration to validate
        :type  cfg:     dict
        :returns:       Flag indicating if the configuration is valid or not
        :rtype:         bool
        """
        if 'usbdmx' not in cfg:
            self.log.error('"usbdmx" not found in configuration')
            return
        i = 0
        for usbdmx in cfg['usbdmx']:
            if 'name' not in usbdmx:
                self.log.error('"name" not found in usbdmx[{0}]'.format(i))
                return
            if 'device' not in usbdmx:
                self.log.error('"device" not found in usbdmx[{0}]'.format(i))
                return
            if 'fixtures' not in usbdmx:
                self.log.error('"fixtures" not found in usbdmx[{0}]'.format(i))
            i += 1
        return True

    def load_config(self):
        """Helper function to load the configuration from disk
        """
        if not os.path.exists(self._cfg_file):
            self.log.error('{0} does not exist'.format(self._cfg_file))
            return
        raw_cfg = open(self._cfg_file, 'r').read()
        cfg = yaml.load(raw_cfg)
        if self.validate_config(cfg):
            self._cfg = cfg
        return True
