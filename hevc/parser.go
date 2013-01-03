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
    "bufio"
    "container/list"
    "errors"
    "fmt"
    "io"
    "strconv"
    "strings"
)

const Version = 0.91

const (
    CMD_VPS = iota
    CMD_SPS
    CMD_PPS
    CMD_FRAME
    CMD_TILE
    CMD_SLICE
    CMD_LCU
    CMD_CU
    CMD_COEF
    CMD_RESI
    CMD_PRED
    CMD_RECO
    CMD_FINAL
)

type Analyzer struct {
    VPSList   *list.List
    SPSList   *list.List
    PPSList   *list.List
    FrameList *list.List
}

var CmdStr2Id map[string]int

func init() {
    CmdStr2Id = make(map[string]int)
    CmdStr2Id["vps"] = CMD_VPS
    CmdStr2Id["sps"] = CMD_SPS
    CmdStr2Id["pps"] = CMD_PPS
    CmdStr2Id["frame"] = CMD_FRAME
    CmdStr2Id["tile"] = CMD_TILE
    CmdStr2Id["slice"] = CMD_SLICE
    CmdStr2Id["lcu"] = CMD_LCU
    CmdStr2Id["cu"] = CMD_CU
    CmdStr2Id["coef"] = CMD_COEF
    CmdStr2Id["resi"] = CMD_RESI
    CmdStr2Id["pred"] = CMD_PRED
    CmdStr2Id["reco"] = CMD_RECO
    CmdStr2Id["final"] = CMD_FINAL
}

func (ha *Analyzer) ParseTrace(traceFile io.Reader) (err error) {
    var line string

    ha.VPSList = list.New()
    ha.SPSList = list.New()
    ha.PPSList = list.New()
    ha.FrameList = list.New()

    reader := bufio.NewReader(traceFile)
    eof := false

    line, err = reader.ReadString('\n')
    if err == io.EOF {
        err = nil
        eof = true
    } else if err != nil {
        return err
    }

    for !eof {
        if strings.Contains(line, "Video Parameter Set") {
            vps := NewVPS(reader)
            if line, err = vps.Parse(); err == nil {
                ha.VPSList.PushBack(vps)
            }
        } else if strings.Contains(line, "Sequence Parameter Set") {
            sps := NewSPS(reader)
            if line, err = sps.Parse(); err == nil {
                for v := ha.VPSList.Front(); v != nil; v = v.Next() {
                    vps := v.Value.(*VPS)
                    if vps.VidParameterSetId == sps.VidParameterSetId {
                        sps.M_vps = vps
                        break
                    }
                }
                ha.SPSList.PushBack(sps)
            }
        } else if strings.Contains(line, "Picture Parameter Set") {
            pps := NewPPS(reader)
            if line, err = pps.Parse(); err == nil {
                for v := ha.SPSList.Front(); v != nil; v = v.Next() {
                    sps := v.Value.(*SPS)
                    if sps.SeqParameterSetId == pps.SeqParameterSetId {
                        pps.M_sps = sps
                        break
                    }
                }
                ha.PPSList.PushBack(pps)
            }
        } else if strings.Contains(line, "Slice Parameter Set") {
            slice := NewSlice(reader)
            if line, err = slice.Parse(); err == nil {
                if slice.FirstSliceInPicFlag == 1 {
                    var poc int

                    if poc, err = slice.GetValue("pic_order_cnt_lsb"); err != nil{
                    	poc = 0
                    	err = nil
                    }

                    frame := NewFrame(poc)
                    for v := ha.PPSList.Front(); v != nil; v = v.Next() {
                        pps := v.Value.(*PPS)
                        if pps.PicParameterSetId == slice.PicParameterSetId {
                            frame.M_pps = pps
                            break
                        }
                    }
                    ha.FrameList.PushBack(frame)
                }
                frame := ha.FrameList.Back().Value.(*Frame)
                slice.M_frame = frame
                frame.SliceList.PushBack(slice)
            }
        } else if strings.Contains(line, "LCU Parameter Set") {
        	fmt.Printf(".")

        	lcu := NewLCU(reader)
        	if line, err = lcu.Parse(); err == nil {
        		frame := ha.FrameList.Back().Value.(*Frame)
        		slice := frame.SliceList.Back().Value.(*Slice)
        		lcu.M_slice = slice
        		frame.LCUList.PushBack(lcu)
        	}
        } else if strings.Contains(line, "CU Parameter Set") {
        	cu := NewCU(reader)
        	if line, err = cu.Parse(); err == nil {
        		frame := ha.FrameList.Back().Value.(*Frame)
        		lcu   := frame.LCUList.Back().Value.(*LCU)
        		cu.M_lcu = lcu
        		lcu.CUList.PushBack(cu)
        	}
        } else if strings.Contains(line, "PU Parameter Set") {
        	pu := NewPU(reader)
        	if line, err = pu.Parse(); err == nil {
        		frame := ha.FrameList.Back().Value.(*Frame)
        		lcu   := frame.LCUList.Back().Value.(*LCU)
        		cu    := lcu.CUList.Back().Value.(*CU)
        		pu.M_cu = cu
        		cu.PUList.PushBack(pu)
        	}
        } else if strings.Contains(line, "TU Parameter Set") {
        	tu := NewTU(reader)
        	if line, err = tu.Parse(); err == nil {
        		frame := ha.FrameList.Back().Value.(*Frame)
        		lcu   := frame.LCUList.Back().Value.(*LCU)
        		cu    := lcu.CUList.Back().Value.(*CU)
        		tu.M_cu = cu
        		cu.TUList.PushBack(tu)
        	}
        } else if strings.Contains(line, "Coefficient Parameter Set") {
        	frame := ha.FrameList.Back().Value.(*Frame)
        	lcu   := frame.LCUList.Back().Value.(*LCU)
        	cu    := lcu.CUList.Back().Value.(*CU)
        	tu	  := cu.TUList.Back().Value.(*TU)
        	coef  := NewDATA(reader, tu.Tu_color, tu.Tu_x, tu.Tu_y, tu.Tu_width, tu.Tu_height, TCOEF)
        	if line, err = coef.Parse(); err == nil {
        		//coef.M_tu = tu
        		tu.CoefList.PushBack(coef)
        	}
        } else if strings.Contains(line, "Residual Parameter Set") {
        	frame := ha.FrameList.Back().Value.(*Frame)
        	lcu   := frame.LCUList.Back().Value.(*LCU)
        	cu    := lcu.CUList.Back().Value.(*CU)
        	tu	  := cu.TUList.Back().Value.(*TU)
        	resi  := NewDATA(reader, tu.Tu_color, tu.Tu_x, tu.Tu_y, tu.Tu_width, tu.Tu_height, TPEL)
        	if line, err = resi.Parse(); err == nil {
        		//resi.M_tu = tu
        		tu.ResiList.PushBack(resi)
        	}
        } else if strings.Contains(line, "Prediction Parameter Set") {
        	frame := ha.FrameList.Back().Value.(*Frame)
        	lcu   := frame.LCUList.Back().Value.(*LCU)
        	cu    := lcu.CUList.Back().Value.(*CU)
        	tu	  := cu.TUList.Back().Value.(*TU)
        	pred  := NewDATA(reader, tu.Tu_color, tu.Tu_x, tu.Tu_y, tu.Tu_width, tu.Tu_height, TPXL)
        	if line, err = pred.Parse(); err == nil {
        		//pred.M_tu = tu
        		tu.PredList.PushBack(pred)
        	}
        } else if strings.Contains(line, "Reconstruction Parameter Set") {
        	frame := ha.FrameList.Back().Value.(*Frame)
        	lcu   := frame.LCUList.Back().Value.(*LCU)
        	cu    := lcu.CUList.Back().Value.(*CU)
        	tu	  := cu.TUList.Back().Value.(*TU)
        	reco  := NewDATA(reader, tu.Tu_color, tu.Tu_x, tu.Tu_y, tu.Tu_width, tu.Tu_height, TPXL)
        	if line, err = reco.Parse(); err == nil {
        		//reco.M_tu = tu
        		tu.RecoList.PushBack(reco)
        	}
        } else if strings.TrimSpace(line) == "" {
            eof = true
            fmt.Printf("\n")
        } else {
            eof = true
            err = errors.New("Unknown Parameter Set\n")
        }

        if err == io.EOF {
            err = nil
            eof = true
        } else if err != nil {
            eof = true
        }
    }

    return err
}

func (ha *Analyzer) ParseYUV(yuvFile io.Reader) (err error) {
	err = nil
	reader := bufio.NewReader(yuvFile)

	for i:=0; i<ha.FrameList.Len(); i=i+1 {
		for v := ha.FrameList.Front(); v != nil; v = v.Next() {
        	frame := v.Value.(*Frame)
			if frame.FrameId == i {
				final_y := NewDATA(reader, 0, 0, 0, frame.M_pps.M_sps.FrameWidth, frame.M_pps.M_sps.FrameHeight, TPXL)
				final_y.YUVReader = yuvFile
				if err = final_y.ParseYUV(); err == nil {
        			//final_y.M_frame = frame
        			frame.YUVList.PushBack(final_y)
        		}else{
        			fmt.Printf("frame Y "+strconv.Itoa(i)+"\n")
        			return err
        		}
        		final_u := NewDATA(reader, 1, 0, 0, frame.M_pps.M_sps.FrameWidth/2, frame.M_pps.M_sps.FrameHeight/2, TPXL)
        		final_u.YUVReader = yuvFile
				if err = final_u.ParseYUV(); err == nil {
        			//final_u.M_frame = frame
        			frame.YUVList.PushBack(final_u)
        		}else{
        			fmt.Printf("frame U "+strconv.Itoa(i)+"\n")
        			return err
        		}
        		final_v := NewDATA(reader, 2, 0, 0, frame.M_pps.M_sps.FrameWidth/2, frame.M_pps.M_sps.FrameHeight/2, TPXL)
        		final_v.YUVReader = yuvFile
				if err = final_v.ParseYUV(); err == nil {
        			//final_v.M_frame = frame
        			frame.YUVList.PushBack(final_v)
        		}else{
        			fmt.Printf("frame V "+strconv.Itoa(i)+"\n")
        			return err
        		}

	            for v := frame.LCUList.Front(); v != nil; v = v.Next() {
	            	lcu := v.Value.(*LCU)
	                for e := lcu.CUList.Front(); e != nil; e = e.Next() {
                		cu := e.Value.(*CU)
                    	for t := cu.TUList.Front(); t != nil; t = t.Next() {
                			tu := t.Value.(*TU)
                			final := frame.GetFinal(cu.Cu_x, cu.Cu_y, tu.Tu_color, tu.Tu_x, tu.Tu_y, tu.Tu_width, tu.Tu_height)
                			//final.M_tu = tu
        					tu.FinalList.PushBack(final)
                		}
                    }
	            }

				break
			}
		}
    }

	return err
}

func (ha *Analyzer) ParseCmd(prefix string, reader *bufio.Reader, out string) {
    cmdLevel := 0
    cmdId := 0

    for {
        fmt.Printf("%s>", prefix)
        line, err := reader.ReadString('\n')
        if err != nil {
            break
        }

        cmdString := strings.Fields(strings.TrimSpace(line))

        if len(cmdString) == 0 {
            continue
        } else if cmdString[0] == "exit" {
            break
        } else if cmdString[0] == "help" {
            NewHelper(cmdString, cmdLevel, cmdId).Parse()
        } else if cmdString[0] == "vps" ||
            cmdString[0] == "sps" ||
            cmdString[0] == "pps" {
            if len(cmdString) == 1 {
                if cmdString[0] == "vps" {
                    fmt.Printf("Total Number of %s: %d\n", cmdString[0], ha.VPSList.Len())
                } else if cmdString[0] == "sps" {
                    fmt.Printf("Total Number of %s: %d\n", cmdString[0], ha.SPSList.Len())
                } else { //cmdString[0] == "pps"
                    fmt.Printf("Total Number of %s: %d\n", cmdString[0], ha.PPSList.Len())
                }
            } else if len(cmdString) == 2 {
                if n, err := strconv.ParseInt(cmdString[1], 0, 0); err != nil {
                    fmt.Printf("%s: Invalid Parameters \"%s\"\n", cmdString[0], cmdString[1])
                } else {
                    if cmdString[0] == "vps" {
                        if int(n) < ha.VPSList.Len() {
                            for v := ha.VPSList.Front(); v != nil; v = v.Next() {
                                vps := v.Value.(*VPS)
                                if vps.VidParameterSetId == int(n) {
                                    vps.ShowInfo()
                                    break
                                }
                            }
                        } else {
                            fmt.Printf("Warning: can't find %d-th %s\n", n, cmdString[0])
                        }
                    } else if cmdString[0] == "sps" {
                        if int(n) < ha.SPSList.Len() {
                            for v := ha.SPSList.Front(); v != nil; v = v.Next() {
                                sps := v.Value.(*SPS)
                                if sps.SeqParameterSetId == int(n) {
                                    sps.ShowInfo()
                                    break
                                }
                            }
                        } else {
                            fmt.Printf("Warning: can't find %d-th %s\n", n, cmdString[0])
                        }
                    } else { //cmdString[0] == "pps"
                        if int(n) < ha.PPSList.Len() {
                            for v := ha.PPSList.Front(); v != nil; v = v.Next() {
                                pps := v.Value.(*PPS)
                                if pps.PicParameterSetId == int(n) {
                                    pps.ShowInfo()
                                    break
                                }
                            }
                        } else {
                            fmt.Printf("Warning: can't find %d-th %s\n", n, cmdString[0])
                        }
                    }
                }
            } else {
                fmt.Printf("Too many parameters in command \"%s\"\n", cmdString[0])
            }
        } else if cmdString[0] == "frame" {
            if len(cmdString) == 1 {
                fmt.Printf("Total Number of %s: %d", cmdString[0], ha.FrameList.Len())
                fmt.Printf("%s\n", out)
            } else if len(cmdString) == 2 {
                if n, err := strconv.ParseInt(cmdString[1], 0, 0); err != nil {
                    fmt.Printf("%s: Invalid Parameters \"%s\"\n", cmdString[0], cmdString[1])
                } else {
                    if int(n) < ha.FrameList.Len() {
                        for v := ha.FrameList.Front(); v != nil; v = v.Next() {
                            frame := v.Value.(*Frame)
                            if frame.FrameId == int(n) {
                                ParseCmdFrame(prefix+"\\"+cmdString[0]+" "+cmdString[1], reader, cmdLevel+1, CmdStr2Id[cmdString[0]], frame)
                                break
                            }
                        }
                    } else {
                        fmt.Printf("Warning: can't find %d-th %s\n", n, cmdString[0])
                    }
                }
            } else {
                fmt.Printf("Too many parameters in command \"%s\"\n", cmdString[0])
            }
        } else {
            fmt.Printf("Unknown Command \"%s\"\n", cmdString[0])
        }
    }
}

func ParseCmdFrame(prefix string, reader *bufio.Reader, cmdLevel int, cmdId int, frame *Frame) {
    for {
        fmt.Printf("%s>", prefix)
        line, err := reader.ReadString('\n')
        if err != nil {
            break
        }

        cmdString := strings.Fields(strings.TrimSpace(line))

        if len(cmdString) == 0 {
            continue
        } else if cmdString[0] == "exit" {
            break
        } else if cmdString[0] == "help" {
            NewHelper(cmdString, cmdLevel, cmdId).Parse()
        } else if cmdString[0] == "tile" {
            if len(cmdString) == 1 {
                frame.ShowTileSummary()
            } else {
                fmt.Printf("Too many parameters in command \"%s\"\n", cmdString[0])
            }
        } else if cmdString[0] == "slice" ||
            cmdString[0] == "lcu" ||
            cmdString[0] == "cu" {
            if len(cmdString) == 1 {
                if cmdString[0] == "slice" {
                    fmt.Printf("Total Number of %s: %d\n", cmdString[0], frame.SliceList.Len())
                }else if cmdString[0] == "lcu" {
                	fmt.Printf("Total Number of %s: %d\n", cmdString[0], frame.LCUList.Len()) //frame.M_pps.M_sps.WidthInLcu*frame.M_pps.M_sps.HeightInLcu)
                }else{//"cu"
                	n := 0;
                	for v := frame.LCUList.Front(); v != nil; v = v.Next() {
                		lcu := v.Value.(*LCU)
                		n = n + lcu.CUList.Len();
                	}
                    fmt.Printf("Total Number of %s: %d\n", cmdString[0], n)
                }
            } else if len(cmdString) == 2 {
                if cmdString[0]=="slice" {
                	if strings.HasPrefix(cmdString[1], "(") && strings.HasSuffix(cmdString[1], ")") {
                		xy := strings.Split(strings.TrimRight(strings.TrimLeft(cmdString[1], "("), ")"), ",");
                		x, err1:= strconv.ParseInt(strings.TrimSpace(xy[0]), 0, 0)
                		y, err2:= strconv.ParseInt(strings.TrimSpace(xy[1]), 0, 0)
                		if err1 !=nil || err2 !=nil || int(x) < 0 || int(x) >= frame.M_pps.M_sps.FrameWidth || int(y) < 0 || int(y) >= frame.M_pps.M_sps.FrameHeight{
                			fmt.Printf("%s: Invalid Parameters \"%s\"\n", cmdString[0], cmdString[1])
                		}else{
                			xInLcu, yInLcu := int(x)/frame.M_pps.M_sps.LcuSize, int(y)/frame.M_pps.M_sps.LcuSize
                			LcuAddr := yInLcu*frame.M_pps.M_sps.WidthInLcu+xInLcu;
                			for v := frame.LCUList.Front(); v != nil; v = v.Next() {
                				lcu := v.Value.(*LCU)
                                if  LcuAddr == lcu.LcuAddr {
                                	slice := lcu.M_slice
                                    ParseCmdSliceLcuCu(prefix+"\\"+cmdString[0]+" "+strconv.Itoa(LcuAddr), reader, cmdLevel+1, CmdStr2Id[cmdString[0]], slice)
                                    break
                                }
                            }
                		}
                	} else if n, err := strconv.ParseInt(cmdString[1], 0, 0); err != nil {
                        fmt.Printf("%s: Invalid Parameters \"%s\"\n", cmdString[0], cmdString[1])
                    } else {
                        if int(n) < frame.SliceList.Len() {
                            for v, i := frame.SliceList.Front(), 0; v != nil; v, i = v.Next(), i+1 {
                                if i == int(n) {
                                    slice := v.Value.(*Slice)
                                    ParseCmdSliceLcuCu(prefix+"\\"+cmdString[0]+" "+cmdString[1], reader, cmdLevel+1, CmdStr2Id[cmdString[0]], slice)
                                    break
                                }
                            }
                        } else {
                            fmt.Printf("Warning: can't find %d-th %s\n", n, cmdString[0])
                        }
                    }
                }else if cmdString[0]=="lcu" {
                	if strings.HasPrefix(cmdString[1], "(") && strings.HasSuffix(cmdString[1], ")") {
                		xy := strings.Split(strings.TrimRight(strings.TrimLeft(cmdString[1], "("), ")"), ",");
                		x, err1:= strconv.ParseInt(strings.TrimSpace(xy[0]), 0, 0)
                		y, err2:= strconv.ParseInt(strings.TrimSpace(xy[1]), 0, 0)
                		if err1 !=nil || err2 !=nil || int(x) < 0 || int(x) >= frame.M_pps.M_sps.FrameWidth || int(y) < 0 || int(y) >= frame.M_pps.M_sps.FrameHeight{
                			fmt.Printf("%s: Invalid Parameters \"%s\"\n", cmdString[0], cmdString[1])
                		}else{
                			xInLcu, yInLcu := int(x)/frame.M_pps.M_sps.LcuSize, int(y)/frame.M_pps.M_sps.LcuSize
                			LcuAddr := yInLcu*frame.M_pps.M_sps.WidthInLcu+xInLcu;
                			for v := frame.LCUList.Front(); v != nil; v = v.Next() {
                				lcu := v.Value.(*LCU)
                                if  LcuAddr == lcu.LcuAddr {
                                    ParseCmdSliceLcuCu(prefix+"\\"+cmdString[0]+" "+strconv.Itoa(LcuAddr), reader, cmdLevel+1, CmdStr2Id[cmdString[0]], lcu)
                                    break
                                }
                            }
                		}
                	} else if n, err := strconv.ParseInt(cmdString[1], 0, 0); err != nil {
                        fmt.Printf("%s: Invalid Parameters \"%s\"\n", cmdString[0], cmdString[1])
                    } else {
                        if int(n) < frame.LCUList.Len() {
                            for v := frame.LCUList.Front(); v != nil; v = v.Next() {
                            	lcu := v.Value.(*LCU)
                                if int(n) == lcu.LcuAddr  {
                                    ParseCmdSliceLcuCu(prefix+"\\"+cmdString[0]+" "+cmdString[1], reader, cmdLevel+1, CmdStr2Id[cmdString[0]], lcu)
                                    break
                                }
                            }
                        } else {
                            fmt.Printf("Warning: can't find %d-th %s\n", n, cmdString[0])
                        }
                    }
                }else{//"cu"
                    if strings.HasPrefix(cmdString[1], "(") && strings.HasSuffix(cmdString[1], ")") {
                		xy := strings.Split(strings.TrimRight(strings.TrimLeft(cmdString[1], "("), ")"), ",");
                		x, err1:= strconv.ParseInt(strings.TrimSpace(xy[0]), 0, 0)
                		y, err2:= strconv.ParseInt(strings.TrimSpace(xy[1]), 0, 0)
                		if err1 !=nil || err2 !=nil || int(x) < 0 || int(x) >= frame.M_pps.M_sps.FrameWidth || int(y) < 0 || int(y) >= frame.M_pps.M_sps.FrameHeight{
                			fmt.Printf("%s: Invalid Parameters \"%s\"\n", cmdString[0], cmdString[1])
                		}else{
                			xInCu, yInCu := int(x), int(y)
                			xInLcu, yInLcu := xInCu/frame.M_pps.M_sps.LcuSize, yInCu/frame.M_pps.M_sps.LcuSize
                			LcuAddr := yInLcu*frame.M_pps.M_sps.WidthInLcu+xInLcu;
                			for v := frame.LCUList.Front(); v != nil; v = v.Next() {
                				lcu := v.Value.(*LCU)
                                if  LcuAddr == lcu.LcuAddr {
                                	for e := lcu.CUList.Front(); e != nil; e = e.Next() {
                						cu := e.Value.(*CU)
                						if xInCu >= cu.Cu_x && xInCu < cu.Cu_x+cu.Cu_size &&
                						   yInCu >= cu.Cu_y && yInCu < cu.Cu_y+cu.Cu_size {
                                    		ParseCmdSliceLcuCu(prefix+"\\"+cmdString[0]+" ("+strconv.Itoa(cu.Cu_x)+","+strconv.Itoa(cu.Cu_y)+")", reader, cmdLevel+1, CmdStr2Id[cmdString[0]], cu)
                                    		break
                                    	}
                                    }
                                    break
                                }
                            }
                		}
                	} else if n, err := strconv.ParseInt(cmdString[1], 0, 0); err != nil {
                        fmt.Printf("%s: Invalid Parameters \"%s\"\n", cmdString[0], cmdString[1])
                    } else {
	                    nInCu := 0;
	                	for v := frame.LCUList.Front(); v != nil; v = v.Next() {
	                		lcu := v.Value.(*LCU)
	                		if int(n) >= nInCu && int(n) < nInCu+lcu.CUList.Len(){
	                			for e := lcu.CUList.Front(); e != nil; e = e.Next() {
                					if  int(n) == nInCu {
                                    	cu := e.Value.(*CU)
                                    	ParseCmdSliceLcuCu(prefix+"\\"+cmdString[0]+" "+cmdString[1], reader, cmdLevel+1, CmdStr2Id[cmdString[0]], cu)
                                    	break
                               		}else{
                               			nInCu = nInCu + 1
                               		}
                                }
                                break
	                		}else{
	                			nInCu = nInCu + lcu.CUList.Len();
	                		}
	                	}

                        if int(n) != nInCu {
                        	fmt.Printf("Warning: can't find %d-th %s\n", n, cmdString[0])
                        }
                    }
                }
            } else {
                fmt.Printf("Too many parameters in command \"%s\"\n", cmdString[0])
            }
        } else {
            fmt.Printf("Unknown command \"%s\"\n", cmdString[0])
        }
    }
}

func ParseCmdSliceLcuCu(prefix string, reader *bufio.Reader, cmdLevel int, cmdId int, slc interface{}) {
    for {
        fmt.Printf("%s>", prefix)
        line, err := reader.ReadString('\n')
        if err != nil {
            break
        }

        cmdString := strings.Fields(strings.TrimSpace(line))

        if len(cmdString) == 0 {
            continue
        } else if cmdString[0] == "exit" {
            break
        } else if cmdString[0] == "help" {
            NewHelper(cmdString, cmdLevel, cmdId).Parse()
        } else if cmdString[0] == "vps" ||
            cmdString[0] == "sps" ||
            cmdString[0] == "pps" ||
            cmdString[0] == "slice" {
            var slice   *Slice
            switch slc.(type) {
            case *Slice: slice = slc.(*Slice)
            case *LCU: slice = slc.(*LCU).M_slice
            case *CU:  slice = slc.(*CU).M_lcu.M_slice
            default: break
            }
            if len(cmdString) == 1 {
                if cmdString[0] == "slice" {
                    slice.ShowInfo()
                }else if cmdString[0] == "pps" {
                    slice.M_frame.M_pps.ShowInfo()
                }else if cmdString[0] == "sps" {
                    slice.M_frame.M_pps.M_sps.ShowInfo()
                }else{//vps
                    slice.M_frame.M_pps.M_sps.M_vps.ShowInfo()
                }
            } else {
                fmt.Printf("Too many parameters in command \"%s\"\n", cmdString[0])
            }
        } else if cmdString[0] == "lcu" && (cmdId == CMD_LCU ||cmdId==CMD_CU) {
            var lcu   *LCU
            switch slc.(type) {
            case *LCU: lcu = slc.(*LCU)
            case *CU:  lcu = slc.(*CU).M_lcu
            default: break
            }
            if len(cmdString) == 1 {
            	lcu.ShowInfo()
            } else {
                fmt.Printf("Too many parameters in command \"%s\"\n", cmdString[0])
            }
        } else if (cmdString[0] == "cu" ||
            cmdString[0] == "pu" ||
            cmdString[0] == "tu" ||
            cmdString[0] == "coef" ||
            cmdString[0] == "resi" ||
            cmdString[0] == "pred" ||
            cmdString[0] == "reco" ||
            cmdString[0] == "final") &&
            cmdId == CMD_CU {
            cu := slc.(*CU)
            if len(cmdString) == 1 {
            	if cmdString[0] == "cu" {
            		cu.ShowInfo()
            	}else if cmdString[0] == "pu" {
            		for e := cu.PUList.Front(); e != nil; e = e.Next() {
                		pu := e.Value.(*PU)
                		pu.ShowInfo()
                	}
            	}else if cmdString[0] == "tu" {
            		for e := cu.TUList.Front(); e != nil; e = e.Next() {
                		tu := e.Value.(*TU)
                		tu.ShowInfo()
                	}
                }else if cmdString[0] == "coef" {
            		for e := cu.TUList.Front(); e != nil; e = e.Next() {
                		tu := e.Value.(*TU)
                		for v := tu.CoefList.Front(); v != nil; v = v.Next() {
                			coef := v.Value.(*DATA)
                			coef.ShowInfo()
                		}
                	}
                }else if cmdString[0] == "resi" {
            		for e := cu.TUList.Front(); e != nil; e = e.Next() {
                		tu := e.Value.(*TU)
                		for v := tu.ResiList.Front(); v != nil; v = v.Next() {
                			resi := v.Value.(*DATA)
                			resi.ShowInfo()
                		}
                	}
                }else if cmdString[0] == "pred" {
            		for e := cu.TUList.Front(); e != nil; e = e.Next() {
                		tu := e.Value.(*TU)
                		for v := tu.PredList.Front(); v != nil; v = v.Next() {
                			pred := v.Value.(*DATA)
                			pred.ShowInfo()
                		}
                	}
                }else if cmdString[0] == "reco" {
            		for e := cu.TUList.Front(); e != nil; e = e.Next() {
                		tu := e.Value.(*TU)
                		for v := tu.RecoList.Front(); v != nil; v = v.Next() {
                			reco := v.Value.(*DATA)
                			reco.ShowInfo()
                		}
                	}
            	}else{
            		for e := cu.TUList.Front(); e != nil; e = e.Next() {
                		tu := e.Value.(*TU)
                		for v := tu.FinalList.Front(); v != nil; v = v.Next() {
                			final := v.Value.(*DATA)
                			final.ShowInfo()
                		}
                	}
                	//fmt.Printf("Perform command: \"%s\"\n", cmdString[0])
                }
            } else {
                fmt.Printf("Too many parameters in command \"%s\"\n", cmdString[0])
            }
        } else {
            fmt.Printf("Unknown command \"%s\"\n", cmdString[0])
        }
    }
}
