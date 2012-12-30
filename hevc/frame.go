package hevc

import (
    //"errors"
    "fmt"
    "container/list"
)

type Frame struct {
    SliceList           *list.List
    LCUList				*list.List
    YUVList				*list.List
    M_pps               *PPS
    FrameId             int
}


func NewFrame(fid int)  *Frame{
	var frame Frame

	frame.FrameId = fid
	frame.SliceList = list.New()
	frame.LCUList = list.New()
	frame.YUVList = list.New()

    return &frame
}

func (frame *Frame) GetFinal (cu_x, cu_y, tu_color, tu_x, tu_y, tu_width, tu_height int)  *DATA{
	var final DATA

	final.M_color = tu_color
	final.M_tu_x = tu_x
	final.M_tu_y = tu_y
	final.M_width = tu_width
	final.M_height = tu_height
	final.M_type = TPXL

	final.M_data = make([][]int, tu_height)
	for j:= 0; j<tu_height; j=j+1 {
		final.M_data[j] = make([]int, tu_width)
	}

	var chroma uint
	if tu_color!=0{
		chroma = 1
	}else{
		chroma = 0
	}

	//fmt.Printf("%d, %d, %d, %d, %d, %d, %d, %d\n",frame.FrameId, cu_x, cu_y, tu_color, tu_x, tu_y, tu_width, tu_height)

	for v := frame.YUVList.Front(); v != nil; v = v.Next() {
    	data := v.Value.(*DATA)
    	if data.M_color == tu_color{
    		//fmt.Printf("%d, %d\n",	data.M_width, data.M_height)

			for j:=0; j<tu_height; j=j+1{
				for i:=0; i<tu_width; i=i+1{
					final.M_data[j][i] = data.M_data[(cu_y>>chroma)+tu_y+j][(cu_x>>chroma)+tu_x+i]
				}
			}
		}
	}

	return &final
}

func (frame *Frame) ShowTileSummary() {
    fmt.Printf("============================================================================\n",)

    tiles_enabled_flag_type, tiles_enabled_flag, err := frame.M_pps.GetTypeValue("tiles_enabled_flag")
    if err!=nil{
        tiles_enabled_flag_type, tiles_enabled_flag, err = frame.M_pps.GetTypeValue("tiles_or_entropy_coding_sync_idc")

        if err!=nil {
    	    fmt.Printf("tiles_enabled_flag not found\n")
    	    return
    	}else{
    	    fmt.Printf("%-62s %-6s: %4d\n", "tiles_or_entropy_coding_sync_idc", tiles_enabled_flag_type, tiles_enabled_flag)
    	}
    }else{
        fmt.Printf("%-62s %-6s: %4d\n", "tiles_enabled_flag", tiles_enabled_flag_type, tiles_enabled_flag)
    }

    if tiles_enabled_flag==1{
    	num_tile_columns_minus1_type, num_tile_columns_minus1, _ := frame.M_pps.GetTypeValue("num_tile_columns_minus1")
    	fmt.Printf("%-62s %-6s: %4d\n", "num_tile_columns_minus1", num_tile_columns_minus1_type, num_tile_columns_minus1)

    	num_tile_rows_minus1_type, num_tile_rows_minus1, _ := frame.M_pps.GetTypeValue("num_tile_rows_minus1")
    	fmt.Printf("%-62s %-6s: %4d\n", "num_tile_rows_minus1", num_tile_rows_minus1_type, num_tile_rows_minus1)

    	uniform_spacing_flag_type, uniform_spacing_flag, _ := frame.M_pps.GetTypeValue("uniform_spacing_flag")
    	fmt.Printf("%-62s %-6s: %4d\n", "uniform_spacing_flag", uniform_spacing_flag_type, uniform_spacing_flag)

    	if uniform_spacing_flag != 1 {
    		for i := 0; i < num_tile_columns_minus1; i = i+1 {
    			t, v, _ := frame.M_pps.GetTypeValue("column_width_minus1[i]")
    			fmt.Printf("%-62s %-6s: %4d\n", "column_width_minus1[i]", t, v)
			}

			for i := 0; i < num_tile_rows_minus1; i = i+1 {
				t, v, _ := frame.M_pps.GetTypeValue("row_height_minus1[i]")
    			fmt.Printf("%-62s %-6s: %4d\n", "row_height_minus1[i]", t, v)
			}
    	}

    	loop_filter_across_tiles_enabled_flag_type, loop_filter_across_tiles_enabled_flag, _ := frame.M_pps.GetTypeValue("loop_filter_across_tiles_enabled_flag")
    	fmt.Printf("%-62s %-6s: %4d\n", "loop_filter_across_tiles_enabled_flag", loop_filter_across_tiles_enabled_flag_type, loop_filter_across_tiles_enabled_flag)
    }

    fmt.Printf("============================================================================\n")
}
