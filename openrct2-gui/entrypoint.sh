#!/bin/bash
x11vnc -forever -create &

awk '/]/{tag=$1} (tag=="[general]") && ($1=="game_path"){$3="\"/mnt/rct2\""} 1' ~/.config/OpenRCT2/config.ini > ~/.config/OpenRCT2/config.ini.tmp && mv ~/.config/OpenRCT2/config.ini.tmp ~/.config/OpenRCT2/config.ini
awk '/]/{tag=$1} (tag=="[sound]") && ($1=="master_sound"){$3="false"} 1' ~/.config/OpenRCT2/config.ini > ~/.config/OpenRCT2/config.ini.tmp && mv ~/.config/OpenRCT2/config.ini.tmp ~/.config/OpenRCT2/config.ini

openrct2 "$@"