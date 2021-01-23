#!/bin/bash


sudo apt install libsdl2-dev libsdl2-ttf-dev

go get -v github.com/veandco/go-sdl2/{sdl,img,mix,ttf}
