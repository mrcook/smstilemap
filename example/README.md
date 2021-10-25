# PNG to Master System ROM image example

This Linux Bash shell script can be used for generating a Sega Master System
ROM file, which will display the image on a SMS emulator. 

This script requires the [WLA DX](https://github.com/vhelin/wla-dx) compiler
and linker to be in your `$PATH`.

Change into this `example` directory and run the command:

    $ sh build.sh

This will produce a `splash.sms` ROM image in the same directory, which can be
loaded into an emulator.


## The Sample Image

The provided image is the loading screen from the ZX Spectrum game JETPAC.
This image has had some minor adjustments to remove the colour clash and align
various graphical elements along 8-pixel boundaries to reduce the number of
needed tiles.

In its original format, it would require 433 tiles to reproduce the image on
the SMS screen. After the changes this was reduced to just 426 unique tiles.
With further editing this could certainly be improved upon.

A reduction of 7 tiles may not seem much, but with a maximum of 448 tiles
available in the Master System VDP memory, every tile counts!


## The Assembly Source Code

The `splash.asm` assembly code is the same code used from Maxim's _How To Program_
tutorial on the [SMS Power](https://www.smspower.org/maxim/HowToProgram/) website,
with a few minor changes to use the generated tile/tilemap data.
