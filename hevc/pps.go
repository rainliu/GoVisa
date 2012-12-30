package hevc

import (
    //"fmt"
    "errors"
    "bufio"
    "container/list"
)

type PPS struct {
    ParameterSet
    M_sps                  *SPS
    SeqParameterSetId       int
    PicParameterSetId       int
}

func NewPPS(reader *bufio.Reader) *PPS{
	var pps PPS
	
	pps.ParameterSet.Reader = reader
	pps.ParameterSet.FieldList = list.New()
	
    return &pps
}

func (ps *PPS) Parse() (line string, err error)  {
    line, err = ps.ParameterSet.Parse()
	
    ps.SeqParameterSetId = -1
    ps.PicParameterSetId = -1
    for e := ps.FieldList.Front(); e != nil; e = e.Next() {
        field := e.Value.(Field)

        if field.Name == "seq_parameter_set_id" {
            ps.SeqParameterSetId = field.Value
        }else if field.Name == "pic_parameter_set_id" {
            ps.PicParameterSetId = field.Value
        }
    }

    if ps.PicParameterSetId < 0 ||
       ps.SeqParameterSetId < 0 {
        err = errors.New("pic_parameter_set_id or seq_parameter_set_id not found or invalid")
    }

    return
}
