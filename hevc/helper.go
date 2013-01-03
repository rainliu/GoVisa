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
    "fmt"
)

type Helper struct {
    cmdString   []string
    cmdLevel    int
    cmdId       int
}

func NewHelper(cmdString []string, cmdLevel int, cmdId int) (h *Helper){
    return &Helper{cmdString, cmdLevel, cmdId}
}

func (h *Helper) Perform() {
    if h.cmdLevel == 0 {
        switch h.cmdString[1] {
        case "exit":fmt.Printf("exit: exit GoHEVC analyzer\n")
        case "vps": fmt.Printf("vps [n]: show num of vps [or n-th vps]\n")
        case "sps": fmt.Printf("sps [n]: show num of sps [or n-th sps]\n")
        case "pps": fmt.Printf("pps [n]: show num of pps [or the n-th pps]\n")
        default:    fmt.Printf("frame [n]: show num of frame [or goto n-th (decoding order) frame]\n") //frame
        }
    } else if h.cmdLevel == 1 {
        switch h.cmdString[1] {
        case "exit": fmt.Printf("exit: back to up-level\n")
        case "tile": fmt.Printf("tile: show current frame's tile info\n")
        case "slice":fmt.Printf("slice [n|(x,y)]: show num of slice [or n-th slice or slice that contains pixel(x,y)]\n")
        case "lcu":  fmt.Printf("lcu [n|(x,y)]: show num of lcu [or n-th lcu or lcu that contains pixel(x,y)]\n")
        default:     fmt.Printf("cu [n|(x,y)]: show num of cu [or goto n-th cu or cu that contains pixel(x,y)]\n")
        }
    } else {
        if h.cmdId == CMD_SLICE {
            switch h.cmdString[1] {
            case "exit":fmt.Printf("exit: back to up-level\n")
            case "vps": fmt.Printf("vps: show current slice's vps info\n")
            case "sps": fmt.Printf("sps: show current slice's sps info\n")
            case "pps": fmt.Printf("pps: show current slice's pps info\n")
            default:    fmt.Printf("slice: show current slice info\n")
            }
        } else if h.cmdId == CMD_LCU {
            switch h.cmdString[1] {
            case "exit": fmt.Printf("exit: back to up-level\n")
            case "vps":  fmt.Printf("vps: show current lcu's vps info\n")
            case "sps":  fmt.Printf("sps: show current lcu's sps info\n")
            case "pps":  fmt.Printf("pps: show current lcu's pps info\n")
            case "slice":fmt.Printf("slice: show current lcu's slice info\n")
            default:     fmt.Printf("lcu: show current lcu info\n")
            }
        } else if h.cmdId == CMD_CU {
            switch h.cmdString[1] {
            case "exit": fmt.Printf("exit: back to up-level\n")
            case "vps":  fmt.Printf("vps: show current cu's vps info\n")
            case "sps":  fmt.Printf("sps: show current cu's sps info\n")
            case "pps":  fmt.Printf("pps: show current cu's pps info\n")
            case "slice":fmt.Printf("slice: show current cu's slice info\n")
            case "lcu":  fmt.Printf("lcu: show current cu's lcu info\n")
            case "cu":   fmt.Printf("cu: show current cu info\n")
            case "pu":   fmt.Printf("pu: show current cu's pu info\n")
            case "tu":   fmt.Printf("tu: show current cu's tu info\n")
            case "coef": fmt.Printf("coef: show current cu's coefficient pixels\n")
            case "resi": fmt.Printf("resi: show current cu's residual pixels\n")
            case "pred": fmt.Printf("pred: show current cu's prediction pixels\n")
            case "reco": fmt.Printf("reco: show current cu's reconstruction pixels\n")
            default:     fmt.Printf("final: show current cu's final pixels\n")
            }
        }
    }
}

func (h *Helper) Parse() {
    switch len(h.cmdString) {
    case 1:
        if h.cmdLevel == 0 {
            fmt.Printf("help [exit|vps|sps|pps|frame]\n")
        } else if h.cmdLevel == 1 {
            fmt.Printf("help [exit|tile|slice|lcu|cu]\n")
        } else {
            if h.cmdId == CMD_SLICE {
                fmt.Printf("help [exit|vps|sps|pps|slice]\n")
            } else if h.cmdId == CMD_LCU {
                fmt.Printf("help [exit|vps|sps|pps|slice|lcu]\n")
            } else if h.cmdId == CMD_CU {
            	fmt.Printf("help [exit|vps|sps|pps|slice|lcu|cu|pu|tu|coef|resi|pred|reco|final]\n")
            }
        }
    case 2:
        if h.cmdLevel == 0 {
            if h.cmdString[1] == "exit" ||
                h.cmdString[1] == "vps" ||
                h.cmdString[1] == "sps" ||
                h.cmdString[1] == "pps" ||
                h.cmdString[1] == "frame" {
                h.Perform()
            } else {
                fmt.Printf("%s: Unknown Parameters \"%s\"\n", h.cmdString[0], h.cmdString[1])
                fmt.Printf("help [exit|vps|sps|pps|frame]\n")
            }
        } else if h.cmdLevel == 1 {
            if h.cmdString[1] == "exit" ||
                h.cmdString[1] == "tile" ||
                h.cmdString[1] == "slice" ||
                h.cmdString[1] == "lcu" ||
                h.cmdString[1] == "cu" {
                h.Perform()
            } else {
                fmt.Printf("%s: Unknown Parameters \"%s\"\n", h.cmdString[0], h.cmdString[1])
                fmt.Printf("help [exit|tile|slice|lcu|cu]\n")
            }
        } else {
            if h.cmdId == CMD_SLICE {
                if h.cmdString[1] == "exit" ||
                    h.cmdString[1] == "vps" ||
                    h.cmdString[1] == "sps" ||
                    h.cmdString[1] == "pps" ||
                    h.cmdString[1] == "slice" {
                    h.Perform()
                } else {
                    fmt.Printf("%s: Unknown Parameters \"%s\"\n", h.cmdString[0], h.cmdString[1])
                    fmt.Printf("help [exit|vps|sps|pps|slice]\n")
                }
            } else if h.cmdId == CMD_LCU {
                if h.cmdString[1] == "exit" ||
                    h.cmdString[1] == "vps" ||
                    h.cmdString[1] == "sps" ||
                    h.cmdString[1] == "pps" ||
                    h.cmdString[1] == "slice" ||
                    h.cmdString[1] == "lcu" {
                    h.Perform()
                } else {
                    fmt.Printf("%s: Unknown Parameters \"%s\"\n", h.cmdString[0], h.cmdString[1])
                    fmt.Printf("help [exit|vps|sps|pps|slice|lcu]\n")
                }
            } else if h.cmdId == CMD_CU {
                if h.cmdString[1] == "exit" ||
                    h.cmdString[1] == "vps" ||
                    h.cmdString[1] == "sps" ||
                    h.cmdString[1] == "pps" ||
                    h.cmdString[1] == "slice" ||
                    h.cmdString[1] == "lcu" ||
                    h.cmdString[1] == "cu" ||
                    h.cmdString[1] == "pu" ||
                    h.cmdString[1] == "tu" ||
                    h.cmdString[1] == "coef" ||
                    h.cmdString[1] == "resi" ||
                    h.cmdString[1] == "pred" ||
                    h.cmdString[1] == "reco" ||
                    h.cmdString[1] == "final" {
                    h.Perform()
                } else {
                    fmt.Printf("%s: Unknown Parameters \"%s\"\n", h.cmdString[0], h.cmdString[1])
                    fmt.Printf("help [exit|vps|sps|pps|slice|lcu|cu|pu|tu|coef|resi|pred|reco|final]\n")
                }
            }
        }
    default:
        fmt.Printf("Too many parameters in \"%s\"\n", h.cmdString[0])
    }
}
