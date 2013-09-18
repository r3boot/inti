package config

type FixtureData struct {
    C []uint8
    I int
}

type FrameData struct {
    F []FixtureData
    D int
}

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

type Controller struct {
    Name string
    Device string
    Id int
    Channels int
}
var Controllers []Controller
