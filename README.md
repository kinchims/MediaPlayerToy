Simple media player toy for children. Configured for use an MF-RC522 RFID chip to read, and write chips and play the media associated with them. 
Media should be located in /config/data with simple integer valued names (1.mp4, 2.mp4, etc.) (extension should not matter, all types of media should be playable).
To use the RFID trainer, start the app with the trainer commandline args. (./media-player trainer). This will prompt you to scan RFID tags so that they can be written a simple id (1-10).

to start the actual application, use ./media-player or start the systemd service provided. Once the RFID reader detects a chip; it will play the media associated with it through DRM.
Provided dockerfile builds the application for ARM64.

Segments.txt shows a list of commands of how to split a larger mp4 file into smaller chunks.

3D prints for the housing can be found in the prints directory. Housing was made for a 5 inch screen, 1.5 inch speaker, and a MF-RC522 reader which should be located near the slot on the top. (cards are inserted and read similar to a game boy advance catridge).
bottom part of housing not provided due to potential different use of buck converters, batteries, SBC, and screen from which I used.


What it does.

1. Plays media that you provided in the files folder based upon the id that the RFID reader scans. It decodes the streams and sinks it to audio output and the DRM so that it can be displayed to the child via the screen.
2. Device will auto turn off if inactive for 5 minutes (configurable via main.go).


Disabling Console terminal for black screen
1. Running the app alone will not remove the terminal between media plays, you must disable the GettyTTY1 service.
