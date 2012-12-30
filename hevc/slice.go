package hevc

import (
    //"fmt"
    "errors"
    "bufio"
    "container/list"
)

type Slice struct {
    ParameterSet
    M_frame                 *Frame
    PicParameterSetId       int
    FirstSliceInPicFlag     int
    SliceAddr				int
}

func NewSlice(reader *bufio.Reader) *Slice{
	var slice Slice
	
	slice.ParameterSet.Reader = reader
	slice.ParameterSet.FieldList = list.New()
	
    return &slice
}

func (ps *Slice) Parse() (line string, err error)  {
    line, err = ps.ParameterSet.Parse()

    ps.PicParameterSetId = -1
    ps.FirstSliceInPicFlag = 0
    for e := ps.FieldList.Front(); e != nil; e = e.Next() {
        field := e.Value.(Field)

        if field.Name == "pic_parameter_set_id" {
            ps.PicParameterSetId = field.Value
        }else if field.Name == "first_slice_in_pic_flag" {
            ps.FirstSliceInPicFlag =  field.Value
        }
    }

    if ps.PicParameterSetId < 0 {
        err = errors.New("pic_parameter_set_id not found or invalid")
    }
    
    if ps.FirstSliceInPicFlag==1{
    	ps.SliceAddr = 0
    }else{
    	ps.SliceAddr, err = ps.ParameterSet.GetValue("slice_address") 
    }

    return line, err
}
