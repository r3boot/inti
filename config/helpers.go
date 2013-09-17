package config

import (
    "errors"
)

func GetFixtureTemplateId (name string) (tid int, err error) {
    for tid = 0; tid < len(FixtureTemplates); tid++ {
        if FixtureTemplates[tid].Name == name {
            return
        }
    }
    tid = -1
    err = errors.New("getFixtureTemplate: "+name+" does not exist")
    return
}

func GetFixtureId (name string) (fid int, err error) {
    for fid = 0; fid < len(Fixtures); fid++ {
        if Fixtures[fid].Name == name {
            return
        }
    }
    fid = -1
    err = errors.New("getFixtureId: "+name+" does not exist")
    return
}

func GetMembershipId (name string) (id int, err error) {
    for gid := 0; gid < len(Groups); gid++ {
        if GroupMembership[gid].Fixture.Name == name {
            gid = id
            err = nil
            return 
        }
    }
    return
}
