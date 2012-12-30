package hevc

import (
    //"fmt"
    "errors"
    "bufio"
    "container/list"
)

type SPS struct {
    ParameterSet
    M_vps                *VPS
    VidParameterSetId     int
    SeqParameterSetId     int
    
    FrameWidth			int
    FrameHeight			int
    LcuSize				int
    WidthInLcu			int
    HeightInLcu			int
}

func NewSPS(reader *bufio.Reader) *SPS{
	var sps SPS
	
	sps.ParameterSet.Reader = reader
	sps.ParameterSet.FieldList = list.New()
	
    return &sps
}

func (ps *SPS) Parse() (line string, err error)  {
    line, err = ps.ParameterSet.Parse()

    ps.VidParameterSetId = -1
    ps.SeqParameterSetId = -1
    for e := ps.FieldList.Front(); e != nil; e = e.Next() {
        field := e.Value.(Field)

        if field.Name == "video_parameter_set_id" {
            ps.VidParameterSetId = field.Value
        }else if field.Name == "seq_parameter_set_id" {
            ps.SeqParameterSetId = field.Value
        }
    }

    if ps.VidParameterSetId < 0 ||
       ps.SeqParameterSetId < 0 {
        err = errors.New("video_parameter_set_id or seq_parameter_set_id not found or invalid")
    }

	ps.FrameWidth, err  = ps.ParameterSet.GetValue("pic_width_in_luma_samples")
	ps.FrameHeight, err = ps.ParameterSet.GetValue("pic_height_in_luma_samples")
	
	var log2_min_coding_block_size_minus3, log2_diff_max_min_coding_block_size int
	log2_min_coding_block_size_minus3, err  = ps.ParameterSet.GetValue("log2_min_coding_block_size_minus3")
	log2_diff_max_min_coding_block_size,err = ps.ParameterSet.GetValue("log2_diff_max_min_coding_block_size")
	ps.LcuSize = 1<<uint(3+log2_min_coding_block_size_minus3+log2_diff_max_min_coding_block_size)

	ps.WidthInLcu = (ps.FrameWidth+ps.LcuSize-1)/ps.LcuSize
	ps.HeightInLcu = (ps.FrameHeight+ps.LcuSize-1)/ps.LcuSize
		
    return
}
