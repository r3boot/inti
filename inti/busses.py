"""
.. module: busses
   :platform: Linux
   :synopsis: Abstraction class representing all DMX busses on a system

.. moduleauthor: Lex van Roon <r3boot@r3blog.nl>
"""

import os

from inti import usbdmx
from inti import fixtures

MAX_USBDMX_BUSSES = 10      # Maximum number of Usb-to-Dmx devices


class Busses:
    """Class representing all DMX/Art-Net busses on this system

    :param logger:  Object containing a configured python logger
    :type  logger:  logger.Logger
    :param cfg:     Dictionary containing the configuration file
    :type  cfg:     dict
    """
    def __init__(self, logger, cfg):
        self.log = logger
        self.cfg = cfg
        self._busses = {}
        self.discover_usbdmx_devices()

    def __getitem__(self, port):
        """Get the object built for a DMX port

        :param port:    Name of the port
        :type  port:    str
        :returns:       Object representing the bus or None
        :rtype:         usbdmx.UsbDmxDevice or None
        """
        try:
            return self._busses[port]
        except KeyError:
            return None

    def __contains__(self, port):
        """Check if port is a member of the available busses

        :param port:    Name of the port
        :type  port:    str
        :returns:       Flag indicating membership
        :rtype:         bool
        """
        return port in self._busses

    def items(self):
        """Helper function which returns an iterable which loops over all
        busses in the system
        """
        return self._busses.items()

    def asdict(self):
        """Helper function which returns the configuration of this bus
        """
        busses = {}
        for bus in self._busses:
            name = self._busses[bus]['name']
            data = self._busses[bus].asdict()
            busses[name] = data
        return busses

    def discover_usbdmx_devices(self):
        """Probe all available local Usb-to-Dmx devices, up to a maximum
        of MAX_USBDMX_BUSSES controllers. All discovered controllers will be
        added to the local dictionary of busses
        """
        for i in range(MAX_USBDMX_BUSSES):
            device_name = '/dev/ttyUSB{0}'.format(i)
            if not os.path.exists(device_name):
                break

            # Check if the device is configured
            port_cfg = self.cfg.get_usbdmx_byport(device_name)
            if not port_cfg:
                self.log.warning('"{0}" is not configured'.format(device_name))
                continue

            # Found a device, configure it
            name = port_cfg['name']
            device = usbdmx.UsbDmxDevice(device_name, name)
            self._busses[device_name] = device

            # Now add all fixtures attached to this device
            for f_cfg in port_cfg['fixtures']:
                if f_cfg['template'] == 'NurdNode':
                    name = f_cfg['name']
                    address = f_cfg['address']
                    fixture = fixtures.NurdNode(device, address, name)
                    self._busses[device_name]['fixtures'][name] = fixture
                else:
                    self.log.warning('Unknown fixture template "{0}"'.format(
                        f_cfg['template']))
                    continue

    def discover_artnet_devices(self):
        """Probe all Art-Net devices on the network
        """
        self.log.info("Art-Net support not yet implemented")
