package dmx

import (
    "errors"
    "log"
    "strconv"
    "github.com/kylelemons/go-gypsy/yaml"
)

const CHAN_FEAT_PWM uint8 = 0x1
const CHAN_FEAT_ONOFF uint8 = 0x2

const MAX_CHANNELS int = 512

type Channel struct {
    Name string
    Value uint8
    Feature uint8
}

type FixtureTemplate struct {
    Name string
    Description string
    Channels []Channel
}
var FixtureTemplates []FixtureTemplate

type Fixture struct {
    Name string
    Id int
    Channels []Channel
}
var Fixtures []Fixture

func Setup (cfg_file string, disable_dmx bool, disable_artnet bool) (err error) {
    if ! disable_dmx {
        DoDmxDiscovery()
    } else {
        log.Print("Disabling DMX discovery")
    }

    if ! disable_artnet {
        DoArtnetDiscovery()
    } else {
        log.Print("Disabling Art-Net discovery")
    }

    if err := ReadConfigFile(cfg_file); err != nil { log.Fatal(err) }
    if err := LoadFixtureTemplates(); err != nil { log.Fatal(err) }
    if err := LoadFixtures(); err != nil { log.Fatal(err) }

    log.Print("Loaded "+strconv.Itoa(len(FixtureTemplates))+" fixture template(s)")
    log.Print("Loaded "+strconv.Itoa(len(Fixtures))+" fixture(s)")

    return
}

func ReadConfigFile (file_name string) (err error) {
    config, err := yaml.ReadFile(file_name)
    if err != nil {
        log.Fatal(err)
    }
    cfgFile = *config

    return
}

func setStr (dst *string, key string) {
    var value string
    var err error
    if value, err = cfgFile.Get(key); err != nil {
        log.Fatal(err)
    }
    *dst = value
}

func setInt (dst *int, key string) {
    var value string
    var err error
    if value, err = cfgFile.Get(key); err != nil {
        log.Fatal(err)
    }
    if *dst, err = strconv.Atoi(value); err != nil {
        log.Fatal(err)
    }
}

func getFixtureTemplateId (name string) (tid int, err error) {
    for tid = 0; tid < len(FixtureTemplates); tid++ {
        if FixtureTemplates[tid].Name == name {
            return
        }
    }
    tid = -1
    err = errors.New("getFixtureTemplate: "+name+" does not exist")
    return
}

func LoadFixtureTemplates () (err error) {
    var tid, cid int = 0, 0
    for {
        base := "templates["+strconv.Itoa(tid)+"]."
        if _, err = cfgFile.Get(base + "name"); err != nil {
            err = nil
            break
        }
        var template = new(FixtureTemplate)
        var feature string
        setStr(&template.Name, base + "name")
        setStr(&template.Description, base + "description")

        for cid = 0; cid < MAX_CHANNELS; cid++ {
            c_base := base + "channels[" + strconv.Itoa(cid) + "]."
            if _, err = cfgFile.Get(c_base + "name"); err != nil {
                err = nil
                break
            }

            var channel = new(Channel)

            setStr(&channel.Name, c_base + "name")
            channel.Value = 0
            feature, err = cfgFile.Get(base + "features")
            if err != nil {
                channel.Feature = CHAN_FEAT_PWM
            } else {
                switch feature {
                default:
                    channel.Feature = CHAN_FEAT_PWM
                case "onoff":
                    channel.Feature = CHAN_FEAT_ONOFF
                }
            }
            template.Channels = append(template.Channels, *channel)
        }

        FixtureTemplates = append(FixtureTemplates, *template)
        tid += 1
    }

    return
}

func LoadFixtures () (err error) {
    var fid, tid int = 0, 0
    for {
        base := "fixtures["+strconv.Itoa(fid)+"]."
        if _, err = cfgFile.Get(base + "name"); err != nil {
            err = nil
            break
        }

        var fixture = new(Fixture)
        var template string
        setStr(&fixture.Name, base + "name")
        setInt(&fixture.Id, base + "id")
        setStr(&template, base + "template")
        if tid, err = getFixtureTemplateId(template); err != nil {
            return
        }

        fixture.Channels = FixtureTemplates[tid].Channels
        Fixtures = append(Fixtures, *fixture)
        fid += 1
    }
    return
}
