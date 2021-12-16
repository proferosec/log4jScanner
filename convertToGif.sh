ffmpeg -i movie.mp4 -vf "fps=10,scale=1600:-1:flags=lanczos,split[s0][s1];[s0]palettegen[p];[s1][p]paletteuse" -loop 0 movie.gif
