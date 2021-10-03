package sms

import background "github.com/mrcook/smstilemap/image"

type SMS struct {
	videoRAM  *VRAM
	colourRAM *CRAM
}

// FromBackgroundImage converts a tile mapped image to SMS data.
func FromBackgroundImage(bg *background.Background) (*SMS, error) {
	sms := SMS{
		videoRAM:  &VRAM{},
		colourRAM: &CRAM{},
	}

	// TODO: need to generate colour palette data from the tiles
	//       should SMS do that, or should it be done in `image` first?
	//       Perhaps in image it could be part of the _validations_?

	// convert all background tiles to planar data and add to tilemap
	for i := 0; i < bg.TileCount(); i++ {
		if tile, err := bg.GetTile(i); err != nil {
			return nil, err
		} else {
			sms.videoRAM.addCharacter(i, tile)
			sms.videoRAM.addTilemapEntry(tile)
		}
	}

	return &sms, nil
}
