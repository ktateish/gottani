// testdata/asmextern/add_amd64.s
#include "textflag.h"

// func Add(a, b int) int
TEXT ·Add(SB), NOSPLIT, $0-24  // locals=0, args=24
    MOVQ  a+0(FP), AX    // a
    MOVQ  b+8(FP), BX    // b
    ADDQ BX, AX          // AX = a + b
    MOVQ AX, ret+16(FP)  // return
    RET
