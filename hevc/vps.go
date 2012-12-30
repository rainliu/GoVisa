package hevc

import (
    //"fmt"
    "errors"
    "bufio"
    "container/list"
)

type VPS struct {
    ParameterSet
    VidParameterSetId  int
}

func NewVPS(reader *bufio.Reader)  *VPS{
	var vps VPS
	
	vps.Reader = reader
	vps.ParameterSet.FieldList = list.New()
    return &vps
}

func (ps *VPS) Parse() (line string, err error)  {
    line, err = ps.ParameterSet.Parse()

    ps.VidParameterSetId = -1
    for e := ps.FieldList.Front(); e != nil; e = e.Next() {
        field := e.Value.(Field)
        if field.Name == "video_parameter_set_id" {
            ps.VidParameterSetId = field.Value
            break
        }
    }

    if ps.VidParameterSetId < 0 {
        err = errors.New("video_parameter_set_id not found or invalid")
    }

    return
}
