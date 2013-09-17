package config

import (
    "log"
    "strconv"
    "github.com/kylelemons/go-gypsy/yaml"
)

var cfgFile yaml.File

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

type Group struct {
    Name string
    Description string
    Fixtures []*Fixture
}
var Groups []Group

type GroupMember struct {
    Fixture *Fixture
    Groups []*Group
}
var GroupMembership []GroupMember

func Setup (file_name string) (err error) {
    if err = ReadConfigFile(file_name); err != nil { return }
    if err = LoadFixtureTemplates(); err != nil { return }
    if err = LoadFixtures(); err != nil { return }
    if err = LoadGroups(); err != nil { return }

    log.Print("Loaded "+strconv.Itoa(len(FixtureTemplates))+" fixture template(s)")
    log.Print("Loaded "+strconv.Itoa(len(Fixtures))+" fixture(s)")
    log.Print("Loaded "+strconv.Itoa(len(Groups)-1)+" group(s)")
    return
}

func ReadConfigFile (file_name string) (err error) {
    config, err := yaml.ReadFile(file_name)
    if err != nil { return }
    cfgFile = *config
    return
}

func setStr (dst *string, key string) {
    var value string
    var err error
    if value, err = cfgFile.Get(key); err != nil { return }
    *dst = value
}

func setInt (dst *int, key string) {
    var value string
    var err error
    if value, err = cfgFile.Get(key); err != nil { return }
    if *dst, err = strconv.Atoi(value); err != nil { return }
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
        if tid, err = GetFixtureTemplateId(template); err != nil {
            return
        }

        fixture.Channels = FixtureTemplates[tid].Channels
        Fixtures = append(Fixtures, *fixture)
        fid += 1
    }
    return
}

func LoadGroups () (err error) {
    var all_group = new(Group)
    all_group.Name = "All"
    all_group.Description = "All fixtures"

    for fid := 0; fid < len(Fixtures); fid++ {
        all_group.Fixtures = append(all_group.Fixtures, &Fixtures[fid])
    }
    Groups = append(Groups, *all_group)

    for gid := 0; gid < MAX_GROUPS; gid++ {
        base := "groups["+strconv.Itoa(gid)+"]."
        if _, err = cfgFile.Get(base + "name"); err != nil {
            err = nil
            break
        }
        var group = new(Group)
        var fid int
        setStr(&group.Name, base + "name")
        setStr(&group.Description, base + "description")

        var fixture_name string
        for mid := 0; mid < MAX_GROUP_MEMBERS; mid++ {
            f_base := base + "fixtures["+strconv.Itoa(mid)+"]."
            if fixture_name, err = cfgFile.Get(f_base + "name"); err != nil {
                err = nil
                break
            }
            if fid, err = GetFixtureId(fixture_name); err != nil { return }

            group.Fixtures = append(group.Fixtures, &Fixtures[fid])

        }

        Groups = append(Groups, *group)
    }

    return
}
