A Chip-8 emulator written in go.
Thanks to Tobias Langhoff for the [awesome guide](https://tobiasvl.github.io/blog/write-a-chip-8-emulator/)

TODO Architecture of the emulator

TODO instructions breakdown

TODO add some sample programs

TODO add execution instructions
`go run main.go ./roms/Pong2.ch8`

## OpenGL Display setup on WSL2
1. Install mesa display drivers
`apt install mesa-utils libglu1-mesa-dev freeglut3-dev mesa-common-dev`
2. Resolve libxxf86vm dependency
`apt install libxxf86vm-dev`
