package hevc

import (
    //"fmt"
    //"errors"
    "bufio"
    "container/list"
)


type TU struct {
    ParameterSet
    CoefList		   *list.List
    ResiList		   *list.List
    PredList		   *list.List
    RecoList		   *list.List
    FinalList		   *list.List
    M_cu               *CU
    Tu_color			int
    Tu_x				int
    Tu_y				int
    Tu_width			int
	Tu_height			int
}


func NewTU(reader *bufio.Reader) *TU{
	var tu TU
	
	tu.ParameterSet.Reader = reader
	tu.ParameterSet.FieldList = list.New()
	tu.CoefList = list.New()
	tu.ResiList = list.New()
	tu.PredList = list.New()
	tu.RecoList = list.New()
	tu.FinalList= list.New()
	
    return &tu
}

func (ps *TU) Parse() (line string, err error)  {
    line, err = ps.ParameterSet.Parse()

	ps.Tu_color, err = ps.ParameterSet.GetValue("tu_color") 
    if err != nil {
    	return line, err
    }
    
	ps.Tu_x, err = ps.ParameterSet.GetValue("tu_x") 
    if err != nil {
    	return line, err
    }
    
    ps.Tu_y, err = ps.ParameterSet.GetValue("tu_y") 
    if err != nil {
    	return line, err
    }
    
	ps.Tu_width, err = ps.ParameterSet.GetValue("tu_width") 
    if err != nil {
    	return line, err
    }
    
    ps.Tu_height, err = ps.ParameterSet.GetValue("tu_height") 
    
    return line, err
}


