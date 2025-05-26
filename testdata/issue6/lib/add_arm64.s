//go:build arm64
// +build arm64

#include "textflag.h"

// func Add(a, b int) int
TEXT ·Add(SB), NOSPLIT, $0-24  // locals=0, args=24
    MOVD    a+0(FP), R0     // a
    MOVD    b+8(FP), R1     // b
    ADD     R1, R0          // R0 = a + b
    MOVD    R0, ret+16(FP)  // return
    RET
