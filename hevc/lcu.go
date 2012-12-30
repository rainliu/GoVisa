package hevc

import (
    //"fmt"
    //"errors"
    "bufio"
    "container/list"
)


type LCU struct {
    ParameterSet
    CUList					*list.List
    M_slice                 *Slice
	LcuAddr					int
}


func NewLCU(reader *bufio.Reader) *LCU{
	var lcu LCU
	
	lcu.ParameterSet.Reader = reader
	lcu.ParameterSet.FieldList = list.New()
	lcu.CUList = list.New()
	
    return &lcu
}

func (ps *LCU) Parse() (line string, err error)  {
    line, err = ps.ParameterSet.Parse()

    ps.LcuAddr, err = ps.ParameterSet.GetValue("lcu_address") 
    
    return line, err
}
