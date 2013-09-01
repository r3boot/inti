package dmx

import (
    "bytes"
    "encoding/binary"
    "errors"
    "log"
    "net"
    "fmt"
    "encoding/hex"
    "strconv"
    "strings"
    "time"
)

const ArtNetProtVer uint16 = 0x0e00
const OpPoll uint16 = 0x2000
const OpPollReply uint16 = 0x2100
const OpOutput uint16 = 0x5000

// Art-Net Priority Codes
const DpLow uint8 = 0x10
const DpMed uint8 = 0x40
const DpHigh uint8 = 0x80
const DpCritical uint8 = 0xe0
const DpVolatile uint8 = 0xff

// Art-Net TalkToMe Codes
const TaArtPollReply = 0x1
const TaDiagnostics = 0x2
const TaUnicast = 0x4

// Art-Net Status1 Codes
const StUbeaPresent = 0x80
const StRDMCapable = 0x40
const StROMBoot = 0x20
const StNotImpl = 0x10
const StPortAddrProgAuthority = 0xc
const StIndicatorState = 0x3
const StWebManagement = 0x1
const StDhcpConfigured = 0x2
const StDhcpCapable = 0x4

const MAX_ARTNET_DEVICES = 1024
const ARTNET_TIMEOUT = 3

const OFFSET_OPCODE = 0x8
const ARTNET_OFFSET_LENGTH = 0x17

type ArtPollPacket struct {
    ID [8]uint8
    OpCode uint16
    ProtVer uint16
    TalkToMe uint8
    Priority uint8
}

type ArtPollReplyPacket struct {
    ID [8]uint8
    OpCode uint16
    IPAddress [4]uint8
    Port uint16
    VersInfo uint16
    NetSwitch uint8
    SubSwitch uint8
    OemHi uint8
    Oem uint8
    UbeaVersion uint8
    Status1 uint8
    EstaManLo uint8
    EstaManHi uint8
    ShortName [18]uint8
    LongName [64]uint8
    NodeReport [64]uint8
    NumPortsHi uint8
    NumPortsLo uint8
    PortTypes [4]uint8
    GoodInput [4]uint8
    GoodOutput [4]uint8
    SwIn [4]uint8
    SwOut [4]uint8
    SwVideo uint8
    SwMacro uint8
    SwRemote uint8
    Spare1 uint8
    Spare2 uint8
    Spare3 uint8
    Style uint8
    MacHi uint8
    Mac1 uint8
    Mac2 uint8
    Mac3 uint8
    Mac4 uint8
    MacLo uint8
    BindIP [4]uint8
    BindIndex uint8
    Status2 uint8
    Filler [26]uint8
}

type ArtDmxPacket struct {
    ID [8]uint8
    OpCode uint16
    ProtVer uint16
    Sequence uint8
    Physical uint8
    SubUni uint8
    Net uint8
    Length uint16
}

type ArtnetDevice struct {
    IP net.IP
    Fd net.UDPConn
    Mac net.HardwareAddr
    Name string
    Description string
    WebManagement bool
    DHCPConfigured bool
    DHCPCapable bool
}
var ArtnetDevices = make([]ArtnetDevice, MAX_ARTNET_DEVICES)
var NumArtnetDevices int = 0

var bcastSendSocket net.UDPConn
var bcastRecvSocket net.UDPConn
var lastSequence uint8 = 0

func init() {
    var cidr net.IP
    var network *net.IPNet

    addresses, err := net.InterfaceAddrs()
    if err != nil {
        log.Fatal(err)
    }

    for _,addr := range addresses {
        if cidr,network, err = net.ParseCIDR(addr.String()); err != nil {
            log.Fatal(err)
        }

        if cidr.IsLoopback() {
            continue
        }
        if strings.Contains(cidr.String(), ":") {
            continue
        }
        log.Print("Sending Art-Net Poll on "+network.String())
        ProbeArtnetDevices(*network)
    }
}

func CloseArtnetSockets() {
    for id := 0; id < NumArtnetDevices; id++ {
        ArtnetDevices[id].Fd.Close()
    }
}

// Convert uint to net.IP
func inet_ntoa(ipnr int64) net.IP {   
    var bytes [4]byte
    bytes[0] = byte(ipnr & 0xFF)
    bytes[1] = byte((ipnr >> 8) & 0xFF)
    bytes[2] = byte((ipnr >> 16) & 0xFF)
    bytes[3] = byte((ipnr >> 24) & 0xFF)

    return net.IPv4(bytes[3],bytes[2],bytes[1],bytes[0])
}

// Convert net.IP to int64
func inet_aton(ipnr net.IP) int64 {      
    bits := strings.Split(ipnr.String(), ".")
    
    b0, _ := strconv.Atoi(bits[0])
    b1, _ := strconv.Atoi(bits[1])
    b2, _ := strconv.Atoi(bits[2])
    b3, _ := strconv.Atoi(bits[3])

    var sum int64
    
    sum += int64(b0) << 24
    sum += int64(b1) << 16
    sum += int64(b2) << 8
    sum += int64(b3)
    
    return sum
}

func byteswap(in uint16) (result uint16) {
    var out uint16 = 0
    var i uint16
    var e uint16 = 2
    for i = 0; i < e; i++ {
        var shifted_right uint16 = (in >> (i * 8));
        var byte_i uint16 = shifted_right & 0xFF;
        var mirrored uint16 = (byte_i << ((e - (i + 1)) * 8));
        out |= mirrored;
    }
    return out;
}

func broadcastAddress(network net.IPNet) (bcast net.IP, err error) {
    var mask int64
    ip := inet_aton(network.IP)
    if mask,err = strconv.ParseInt(network.Mask.String(), 16, 64); err != nil {
        log.Fatal(err)
    }
    bcast = inet_ntoa(ip | (0xffffffff ^ mask))

    return
}

func ProbeArtnetDevices(network net.IPNet) (err error) {
    var p ArtPollPacket
    var buf = make([]byte, 512)
    var broadcast net.IP

    if broadcast, err = broadcastAddress(network); err != nil {
        log.Fatal(err)
    }

    rs, err := net.ListenUDP("udp4", &net.UDPAddr{
        IP: broadcast,
        Port: 6454,
    })
    rs.SetDeadline(time.Now().Add(ARTNET_TIMEOUT * time.Second))
    defer rs.Close()

    ws, err := net.DialUDP("udp4", nil, &net.UDPAddr{
        IP:   broadcast,
        Port: 6454,
    })
    ws.SetDeadline(time.Now().Add(ARTNET_TIMEOUT * time.Second))
    defer ws.Close()

    p = ConstructArtPollPacket((TaArtPollReply | TaDiagnostics | TaUnicast))
    b := new(bytes.Buffer)
    if err = binary.Write(b, binary.LittleEndian, p); err != nil {
        log.Fatal(err)
    }
    ws.Write(b.Bytes())

    var opCode uint16 = 0
    var device ArtnetDevice
    var fd *net.UDPConn

    for {
        if _, _, err = rs.ReadFrom(buf); err != nil {
            break
        }

        ob := bytes.NewBuffer(buf[OFFSET_OPCODE:OFFSET_OPCODE+2])
        if err = binary.Read(ob, binary.LittleEndian, &opCode); err != nil {
            log.Fatal(err)
        }

        switch opCode {
        default:
            fmt.Print(hex.Dump(buf))
        case OpPoll:
        case OpPollReply:
            if device, err = ParseArtPollReplyPacket(buf); err != nil {
                log.Fatal(err)
            }

            fd, err = net.DialUDP("udp4", nil, &net.UDPAddr{
                IP:   device.IP,
                Port: 6454,
            })
            device.Fd = *fd

            ArtnetDevices[NumArtnetDevices] = device
            NumArtnetDevices += 1
        }


    }

    log.Print("Found "+strconv.Itoa(NumArtnetDevices)+" Art-Net device(s)")

    return
}

func ConstructArtDmxPacket() (p ArtDmxPacket) {
    p.ID[0] = 'A'
    p.ID[1] = 'r'
    p.ID[2] = 't'
    p.ID[3] = '-'
    p.ID[4] = 'N'
    p.ID[5] = 'e'
    p.ID[6] = 't'
    p.ID[7] = '\x00'
    p.OpCode = OpOutput
    p.ProtVer = ArtNetProtVer
    p.Sequence = lastSequence
    p.Physical = 0x0
    p.SubUni = 0x0
    p.Net = 0x0
    lastSequence += 1
    return
}

func ConstructArtPollPacket(ttmFlags uint8) (p ArtPollPacket) {
    p.ID[0] = 'A'
    p.ID[1] = 'r'
    p.ID[2] = 't'
    p.ID[3] = '-'
    p.ID[4] = 'N'
    p.ID[5] = 'e'
    p.ID[6] = 't'
    p.ID[7] = '\x00'
    p.OpCode = OpPoll
    p.ProtVer = ArtNetProtVer
    p.TalkToMe = ttmFlags
    p.Priority = DpHigh
    return
}

func ParseArtPollReplyPacket(buf []byte) (d ArtnetDevice, err error) {
    var p ArtPollReplyPacket
    err = binary.Read(bytes.NewBuffer(buf), binary.LittleEndian, &p)
    if err != nil {
        log.Fatal(err)
    }

    d.IP = net.IPv4(p.IPAddress[0],p.IPAddress[1],p.IPAddress[2], p.IPAddress[3])
    d.Name = fmt.Sprintf("%s", p.ShortName)
    d.Description = fmt.Sprintf("%s", p.LongName)
    d.WebManagement = (p.Status2 & StWebManagement) == StWebManagement
    d.DHCPConfigured = (p.Status2 & StDhcpConfigured) == StDhcpConfigured
    d.DHCPCapable = (p.Status2 & StDhcpCapable) == StDhcpCapable
    d.Mac, err = net.ParseMAC(fmt.Sprintf("%x0:%x:%x:%x:%x:%x", p.MacHi, p.Mac1, p.Mac2, p.Mac3, p.Mac4, p.MacLo))
    if err != nil {
        log.Fatal(err)
    }

    return
}

func GetArtnetDeviceId(name string) (id int, err error) {
    for id = 0; id < MAX_DEVICES; id++ {
        if ArtnetDevices[id].IP.String() == name {
            err = nil
            return
        }
    }
    err = errors.New("dmx.getArtnetDeviceId: No such device")

    return
}

func SendArtnetFrame(dev_id int, frame []uint8) (err error) {
    p := ConstructArtDmxPacket()
    p.Length = byteswap(uint16(len(frame)))
    b := new(bytes.Buffer)
    if err = binary.Write(b, binary.LittleEndian, p); err != nil {
        log.Fatal(err)
    }
    buf := append(b.Bytes(), frame...)
    fmt.Print(hex.Dump(buf))
    ArtnetDevices[dev_id].Fd.Write(buf)
    return
}
