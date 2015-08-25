"""
.. module:: usbdmx
   :platform: Linux
   :synopsis: Abstraction of a Usb-to-Dmx bus

... moduleauthor:: Lex van Roon <r3boot@r3blog.nl>
"""

import sys

# Handle external dependencies
try:
    import serial
except ImportError:
    print('pyserial not found, please run "pip install pyserial"')
    sys.exit(1)

MAX_CHANNELS = 512          # Maximum number of channels on a DMX Bus
BAUD_SYNC = 19200           # Baud-rate used during DMX synchronization
BAUD_SEND = 250000          # Baud-rate used during DMX buffer transfer


class UsbDmxDevice:
    """Class representing a single USB-to-DMX device.

    :param port:    Port to use
    :type  port:    str
    :param name:    Descriptive name for the USB-to-DMX device
    :type  name:    str
    """
    _cfg = {}

    def __init__(self, port, name):
        self._cfg = {
            'port': port,                   # Port to use
            'name': name,                   # Descriptive name for this bus
            'buffer': [0] * MAX_CHANNELS,   # Buffer representing this bus
            'fixtures': {},                 # Fixtures attached to this bus
        }

    def __getitem__(self, key):
        """Helper function to return the current value for a key
        in the local configuration dictionary

        :param key:     Key to lookup
        :type  key:     str
        :returns:       Value of the configuration pointed to by key or None
        :rtype:         obj or None
        """
        try:
            return self._cfg[key]
        except KeyError:
            return None

    def __setitem__(self, key, value):
        """Helper function to set a value for a key in the local
        configuration dictionary

        :param key:     Key to configure
        :type  key:     str
        :param value:   New value for this key
        :type  value:   obj
        """
        if key not in self._cfg:
            return
        self._cfg[key] = value

    def __contains__(self, key):
        """Helper function to check if a value is present in the local
        configuration dictionary

        :param key:     Key to lookup
        :type  key:     str
        :returns:       Flag indicating existence of the key
        :rtype:         bool
        """
        return key in self._cfg

    def __repr__(self):
        """Helper function which returns a textual representation of this bus

        :returns:   Name of this bus
        :rtype:     str
        """
        return self._cfg['name']

    def asdict(self):
        """Helper function which returns the local configuration dictionary
        for this bus

        :returns:   Dictionary representing the configuration for this bus
        :rtype:     dict
        """
        cfg = self._cfg
        fixtures = {}
        for k,v in cfg['fixtures'].items():
            fixtures[k] = v.asdict()
        cfg['fixtures'] = fixtures
        return cfg

    def transfer(self):
        """Transfer the contents of the DMX buffer to the DMX bus, thereby
        setting all devices to the values represented by the buffer.
        """
        packet = [0x00] + self._cfg['buffer']

        # Perform DMX synchronization at 19200 baud
        bus = serial.Serial(self._cfg['port'], BAUD_SYNC)
        bus.write([0x00])
        bus.close()

        # Blast the buffer onto the DMX Bus at 250000 baud
        bus = serial.Serial(self._cfg['port'], BAUD_SEND)
        bus.write(packet)
        bus.close()
