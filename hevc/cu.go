package hevc

import (
    //"fmt"
    //"errors"
    "bufio"
    "container/list"
)


type CU struct {
    ParameterSet
    PUList			    *list.List
    TUList				*list.List
    M_lcu               *LCU
	Cu_x				int
	Cu_y				int
	Cu_size				int
}


func NewCU(reader *bufio.Reader) *CU{
	var cu CU
	
	cu.ParameterSet.Reader = reader
	cu.ParameterSet.FieldList = list.New()
	cu.PUList = list.New()
	cu.TUList = list.New()
	
    return &cu
}

func (ps *CU) Parse() (line string, err error)  {
    line, err = ps.ParameterSet.Parse()

    ps.Cu_x, err = ps.ParameterSet.GetValue("cu_x") 
    if err != nil {
    	return line, err
    }
    
    ps.Cu_y, err = ps.ParameterSet.GetValue("cu_y") 
    if err != nil {
    	return line, err
    }
    
    ps.Cu_size, err = ps.ParameterSet.GetValue("cu_size")
    
    return line, err
}
