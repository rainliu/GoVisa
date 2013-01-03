/* The copyright in this software is being made available under the BSD
 * License, included below. This software may be subject to other third party
 * and contributor rights, including patent rights, and no such rights are
 * granted under this license.
 *
 * Copyright (c) 2012-2013, H265.net
 * All rights reserved.
 *
 * Redistribution and use in source and binary forms, with or without
 * modification, are permitted provided that the following conditions are met:
 *
 *  * Redistributions of source code must retain the above copyright notice,
 *    this list of conditions and the following disclaimer.
 *  * Redistributions in binary form must reproduce the above copyright notice,
 *    this list of conditions and the following disclaimer in the documentation
 *    and/or other materials provided with the distribution.
 *  * Neither the name of the H265.net nor the names of its contributors may
 *    be used to endorse or promote products derived from this software without
 *    specific prior written permission.
 *
 * THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS"
 * AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE
 * IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE
 * ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS
 * BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR
 * CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF
 * SUBSTITUTE GOODS OR SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS
 * INTERRUPTION) HOWEVER CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN
 * CONTRACT, STRICT LIABILITY, OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE)
 * ARISING IN ANY WAY OUT OF THE USE OF THIS SOFTWARE, EVEN IF ADVISED OF
 * THE POSSIBILITY OF SUCH DAMAGE.
 */
 
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
