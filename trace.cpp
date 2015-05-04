/**************************************************************************//**
 *
 * Copyright 1998-2015 NetBurner, Inc.  ALL RIGHTS RESERVED
 *   Permission is hereby granted to purchasers of NetBurner Hardware
 *   to use or modify this computer program for any use as long as the
 *   resultant program is only executed on NetBurner provided hardware.
 *
 *   No other rights to use this program or it's derivatives in part or
 *   in whole are granted.
 *
 *   It may be possible to license this or other NetBurner software for
 *   use on non-NetBurner Hardware.
 *   Please contact sales@Netburner.com for more information.
 *
 *   NetBurner makes no representation or warranties
 *   with respect to the performance of this computer program, and
 *   specifically disclaims any responsibility for any damages,
 *   special or consequential, connected with the use of this program.
 *
 *---------------------------------------------------------------------
 * NetBurner, Inc.
 * 5405 Morehouse Drive
 * San Diego, California 92121
 *
 * information available at:  http://www.netburner.com
 * E-Mail info@netburner.com
 *
 * Support is available: E-Mail support@netburner.com
 *
 *****************************************************************************/

#include <nettypes.h>
#include <sim.h>

#define TRACE_BASE 0x10000000

#define RGPIO_WR (*(rgpio_wr_struct *)0x8C000000)
#define RGPIO_RD (*(rgpio_rd_struct *)0x8C000000)

extern "C" {
    void trace_exception();
}

typedef struct {
   vuword  dir;            /* 0x8C00_0000 -> 0x8C00_0001 - (Read)  Data Direction Register
                                                            (Write) Data Direction Register                         */
   vuword  data;            /* 0x8C00_0002 -> 0x8C00_0003 - (Read)  Write Data Register
                                                            (Write) Write Data Register                             */
   vuword  enb;             /* 0x8C00_0004 -> 0x8C00_0005 - (Read)  Pin Enable Register
                                                            (Write) Pin Enable Register                             */
   vuword  clr;             /* 0x8C00_0006 -> 0x8C00_0007 - (Read)  Write Data Register
                                                            (Write) Write Data Clear Register                       */
   vuword  set;             /* 0x8C00_000A -> 0x8C00_000B - (Read)  Write Data Register
                                                            (Write) Write Data Set Register                         */
   vuword  tog;             /* 0x8C00_000E -> 0x8C00_000F - (Read)  Write Data Register
                                                            (Write) Write Data Toggle Register                      */
} rgpio_wr_struct;

typedef struct {
   vuword  dir;            /* 0x8C00_0000 -> 0x8C00_0001 - (Read)  Data Direction Register
                                                            (Write) Data Direction Register                         */
   vuword  data;            /* 0x8C00_0002 -> 0x8C00_0003 - (Read)  Write Data Register
                                                            (Write) Write Data Register                             */
   vuword  enb;             /* 0x8C00_0004 -> 0x8C00_0005 - (Read)  Pin Enable Register
                                                            (Write) Pin Enable Register                             */
   vuword  data1;
   vuword  dir1;
   vuword  data2;
   vuword  dir2;
   vuword  data3;
} rgpio_rd_struct;

void EnableTrace(int cs)
{
    // Chip selects 0 and 2 are used by the main system
    if ((cs == 0) || (cs == 2)) { return; }

    sim1.gpio.par_timer &= 0xFC; // clear T0 PAR, setting as GPIO
    sim1.gpio.pddr_e |= 0x80;
    sim1.gpio.srcr_timer |= 0x03;
    sim1.gpio.ppdsdr_e = 0x80;
    // Setup NULL catching FlexBus chip select
    sim2.cs[cs].csar = TRACE_BASE; // Set base address to 0x00xx_xxxx
    sim2.cs[cs].cscr = 0x0000FD40; // Enable internal transfer acknowledge
    sim1.ccm.misccr2 |= 0x00000002; // Set the FlexBus to run at f_sys/4
    sim2.cs[cs].csmr = 0x00000001; // Set the valid bit for this chipselect

//    RGPIO_WR.enb = RGPIO_RD.enb | 0x0010;
//    RGPIO_WR.dir = RGPIO_RD.dir | 0x0010;
//    RGPIO_WR.set = 0x0010;

    vector_base.table[9] = (uint32_t)trace_exception;

    asm("move.l %d0, %sp@-");
    asm("move.w %sr, %d0");
    asm("bset #15, %d0");
    asm("move.w %d0, %sr");
    asm("move.l %sp@+, %d0");

}

