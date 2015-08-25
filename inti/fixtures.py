"""
.. module: fixtures
   :platform: Linux
   :synopsis: Abstraction layer around DMX fixtures

.. moduleauthor:: Lex van Roon <r3boot@r3blog.nl>
"""

# Various constants used to denote the channels
RED = 'red'
GREEN = 'green'
BLUE = 'blue'
PAN = 'pan'
TILT = 'tilt'
PROG = 'prog'


class Fixture:
    """Base class representing a fixture without channels. It's purpose is to
    provide generic functions for derived fixtures.

    :param bus:     Object representing the bus this fixture is attached to
    :type  bus:     usbdmx.UsbDmxDevice
    :param address: Address of the fixture on the bus
    :type  address: int
    :param name:    Descriptive name for the fixture
    :type  name:    str
    """
    bus = None
    address = None
    name = 'Fixture'
    channels = {}

    def __init__(self, bus, address, name):
        self.bus = bus
        self.address = address
        self.name = name

    def __getitem__(self, channel):
        """Get the value from a channel if it is defined

        :param channel:     Name of the channel
        :type  channel:     str
        :returns:           Current value of the channel or None
        :rtype:             int or None
        """
        if channel not in self.channels:
            return

        addr = self.address + self.channels[channel]
        return self.bus.buffer[addr]

    def __setitem__(self, channel, value):
        """Set the value of a channel if it is defined. Value will be capped
        between 0 and 255

        :param channel:     Name of the channel
        :type  channel:     str
        :param value:       New value for this channel
        :type  value:       int
        """
        if channel not in self.channels:
            return

        if value < 0:
            value = 0
        elif value > 255:
            value = 255

        addr = self.address + self.channels[channel]
        self.bus.buffer[addr] = value

    def __repr__(self):
        """Helper function which returns the name of the fixture
        """
        return self.name


class NurdNode(Fixture):
    """Class representing a NurdNode (see http://nodes.nurdspace.nl/)

    A NurdNode provides a 3-channel 3W RGB led combined with servo's for
    pan&tilt control.
    """
    channels = {
        RED:   0,
        GREEN: 1,
        BLUE:  2,
        PAN:   3,
        TILT:  4,
        PROG:  5,
    }

    def __init__(self, bus, address, name):
        Fixture.__init__(self, bus, address, name)
