# smstilemap - a Sega Master System tile/sprite library

Most tools for creating sprites and background tiles for the Sega Master System
(SMS) are only available on the Windows operating system, making it challenging
for macOS and Linux users.

This small [Go language](https://golang.org) library was created to aid
developers converting raster images (e.g. PNG) to the SMS tile, palette, and
tilemap (name table) data.

The source image is converted into 8x8 pixel tiles, duplicate tiles are removed
(using both vertical and horizontal flipping comparison), SMS colour palette
data is generated from the image, and a tilemap is generated.

Note: if the colours in the source image do not match exactly those on the SMS,
a nearest match conversion will be attempted. This can have an undesirable
effect, so it's recommend to follow the image generation guide below.


## Usage

* TODO


## Guide to generating Sega Master System compatible images

Any graphics or sprite editor that can control the colour palette and export to
the PNG can be used.

Master System specs:

* 256x192 pixels screen size (viewport)
* 64 predefined colours
* 32 palette colours: 16 for background tiles and 16 for sprites or background tiles
* tile size: 8x8 pixels
* maximum of 448 unique tiles

When creating the source image, the above dimensions, colours, and palette size
should be used. See [sms/colour.go] for a list of the RGB colours the SMS
supports.

The Master System screen viewport would require 768 unique tiles to fill it.
As the SMS can only hold a maximum of 448 tiles, images need to be crafted for
tile re-use. Careful alignment along 8 pixel boundaries and utilising flipped
tiles (vertical and horizontal) can help to achieve this maximum tile usage.


## TODO

* support SMS sprite sheet by using only palette #2
* validate the source image before processing:
  - colours matches the 64 available on the SMS
  - colour palette does not exceed 32 colours, or 16 for sprite sheets


## LICENSE

Copyright (c) 2021 Michael R. Cook. All rights reserved.

This work is licensed under the terms of the MIT license.
For a copy, see <https://opensource.org/licenses/MIT>.
