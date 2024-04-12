; Sega Master System example program
;
; Display a splash screen on the SMS using background tiles.
; Based on the tutorial: https://www.smspower.org/maxim/HowToProgram/
;
; Modified 2024-04-10 by Michael R. Cook
;
;==============================================================
; WLA-DX banking setup
;==============================================================
.memorymap
defaultslot 0
slotsize $8000
slot 0 $0000
.endme

.rombankmap
bankstotal 1
banksize $8000
banks 1
.endro

;==============================================================
; SMS defines
;==============================================================
.define VDPControl $bf
.define VDPData $be
.define VRAMWrite $4000
.define CRAMWrite $c000

;==============================================================
; SDSC tag and SMS rom header
;==============================================================
.sdsctag 1.0,"smstilemap example","Example program to display a splash screen.","Maxim & Michael R. Cook"

.bank 0 slot 0
.org $0000
;==============================================================
; Boot section
;==============================================================
    di              ; disable interrupts
    im 1            ; Interrupt mode 1
    jp main         ; jump to main program

.org $0066
;==============================================================
; Pause button handler
;==============================================================
    ; Do nothing
    retn

;==============================================================
; Main program
;==============================================================
main:
    ld sp, $dff0

    ;==============================================================
    ; Set up VDP registers
    ;==============================================================
    ld hl, VDPInitData
    ld b, VDPInitDataEnd-VDPInitData
    ld c, VDPControl
    otir

    ;==============================================================
    ; Clear VRAM
    ;==============================================================
    ; 1. Set VRAM write address to $0000
    ld hl, $0000 | VRAMWrite
    call SetVDPAddress
    ; 2. Output 16KB of zeroes
    ld bc, $4000     ; Counter for 16KB of VRAM
-:  xor a
    out (VDPData), a ; Output to VRAM address, which is auto-incremented after each write
    dec bc
    ld a, b
    or c
    jr nz, -

    ;==============================================================
    ; Load palette
    ;==============================================================
    ; 1. Set VRAM write address to CRAM (palette) address 0
    ld hl, $0000 | CRAMWrite
    call SetVDPAddress
    ; 2. Output colour data
    ld hl, PaletteData
    ld bc, PaletteDataEnd-PaletteData
    call CopyToVDP

    ;==============================================================
    ; Load tiles
    ;==============================================================
    ; 1. Set VRAM write address to tile index 0
    ld hl, $0000 | VRAMWrite
    call SetVDPAddress
    ; 2. Output tile data
    ld hl, TileData              ; Location of tile data
    ld bc, TileDataEnd-TileData  ; Counter for number of bytes to write
    call CopyToVDP

    ;==============================================================
    ; Write splash screen tiles to name table
    ;==============================================================
    ; 1. Set VRAM write address to tilemap index 0
    ld hl, $3800 | VRAMWrite
    call SetVDPAddress
    ; 2. Output tilemap data
    ld hl, Tilemap
    ld bc, TilemapEnd-Tilemap ; Counter for number of bytes to write
    call CopyToVDP

    ; Turn screen on
    ld a,%01000000
;          ||||||`- Zoomed sprites -> 16x16 pixels
;          |||||`-- Doubled sprites -> 2 tiles per sprite, 8x16
;          ||||`--- Mega Drive mode 5 enable
;          |||`---- 30 row/240 line mode
;          ||`----- 28 row/224 line mode
;          |`------ VBlank interrupts
;          `------- Enable display
    out (VDPControl),a
    ld a,$81
    out (VDPControl),a

    ; Infinite loop to stop program
-:  jr -

;==============================================================
; Helper functions
;==============================================================

SetVDPAddress:
; Sets the VDP address
; Parameters: hl = address
    push af
        ld a,l
        out (VDPControl),a
        ld a,h
        out (VDPControl),a
    pop af
    ret

CopyToVDP:
; Copies data to the VDP
; Parameters: hl = data address, bc = data length
; Affects: a, hl, bc
-:  ld a,(hl)    ; Get data byte
    out (VDPData),a
    inc hl       ; Point to next letter
    dec bc
    ld a,b
    or c
    jr nz,-
    ret

;==============================================================
; Data
;==============================================================

; VDP initialisation data
VDPInitData:
.db $04,$80 ; Register $00 - Mode Control No. 1
.db $00,$81 ; Register $01 - Mode Control No. 2
.db $ff,$82 ; Register $02 - Name Table Base Address; &ff = address at &3800
.db $ff,$83 ; Register $03 - Color Table Base Address (no effect)
.db $ff,$84 ; Register $04 - Background Pattern Generator Base Address (no effect)
.db $ff,$85 ; Register $05 - Sprite Attribute Table Base Address
.db $ff,$86 ; Register $06 - Sprite Pattern Generator Base Address
.db $ff,$87 ; Register $07 - Overscan/Backdrop Color (using sprite palette)
.db $00,$88 ; Register $08 - Background X Scroll
.db $00,$89 ; Register $09 - Background Y Scroll
.db $ff,$8a ; Register $0A - Line counter (turn off line interupt requests)
VDPInitDataEnd:

; include the screen data; tilemap, palette, tiles.
.include "jetpac.asm"
