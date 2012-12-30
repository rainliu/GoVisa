package hevc

import (
    //"fmt"
    //"errors"
    "bufio"
    "container/list"
)


type PU struct {
    ParameterSet
    M_cu               *CU
}


func NewPU(reader *bufio.Reader) *PU{
	var pu PU
	
	pu.ParameterSet.Reader = reader
	pu.ParameterSet.FieldList = list.New()
	
    return &pu
}

func (ps *PU) Parse() (line string, err error)  {
    line, err = ps.ParameterSet.Parse()

    return line, err
}

