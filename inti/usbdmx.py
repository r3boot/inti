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
    def __init__(self, port, name):
        self.port = port                    # Port to use
        self.name = name
        self.buffer = [0] * MAX_CHANNELS    # Buffer representing the DMX bus
        self.fixtures = {}

    def transfer(self):
        """Transfer the contents of the DMX buffer to the DMX bus, thereby
        setting all devices to the values represented by the buffer.
        """
        packet = 0x00 + self.buffer

        # Perform DMX synchronization at 19200 baud
        bus = serial.Serial(self.port, BAUD_SYNC)
        bus.write(0x00)
        bus.close()

        # Blast the buffer onto the DMX Bus at 250000 baud
        bus = serial.Serial(self.port, BAUD_SEND)
        bus.write(packet)
        bus.close()
