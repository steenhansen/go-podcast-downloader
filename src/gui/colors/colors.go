package colors

import "image/color"

var INPUT_COLOR = color.NRGBA{R: 0x88, G: 0x88, B: 0xff, A: 0xff} // blue for circle/square

var BLACK_TEXT = color.NRGBA{R: 0x00, G: 0x00, B: 0x00, A: 0xff} // black text in light mode
var WHITE_TEXT = color.NRGBA{R: 0xff, G: 0xff, B: 0xff, A: 0xff} // white text in dark mode

var UN_SELECTED_RADIO_CHEK = INPUT_COLOR                                  // blue circle/square and Input Text box
var SELECTED_RADIO_CHEK = color.NRGBA{R: 0xb2, G: 0x00, B: 0xff, A: 0xff} // dark purple circle/square

var CHECKBOX_PERIMETER = color.NRGBA{R: 0xd6, G: 0x7f, B: 0xff, A: 0xff} // light purple

var CURRENT_FOCUS_CIRCLE = color.NRGBA{R: 0x88, G: 0xff, B: 0x88, A: 0x44} // light green focus circle

var HOVERING_CIRCLE = color.NRGBA{R: 0x88, G: 0xff, B: 0x88, A: 0x88} // green hover circle

var BUTTON_LIGHT_COLOR = INPUT_COLOR                                    // lighblue circle/square
var BUTTON_DARK_COLOR = color.NRGBA{R: 0x44, G: 0x44, B: 0xff, A: 0xff} // dark-blue buttons for darkmode

var DISABLED_LIGHT_CHECKBOX_TEXT = color.NRGBA{R: 0x55, G: 0x55, B: 0x55, A: 0xFF} // grey disabled checkbox text

var DISABLED_DARK_CHECKBOX_TEXT = color.NRGBA{R: 0xaa, G: 0xaa, B: 0xaa, A: 0xFF} // grey disabled checkbox text

var TRANSLUCENT_BACKGROUND = color.NRGBA{R: 0, G: 255, B: 0, A: 127} // see-through green for background of download button
var BLACK_DOWNLOAD = color.NRGBA{R: 0, G: 0, B: 0, A: 255}           // black text for download button
