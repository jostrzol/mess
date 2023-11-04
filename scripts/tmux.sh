#!/bin/sh

session="mess"

tmux kill-session -t $session 2>/dev/null || true

window=1
tmux new-session  -ds $session -n 'fe'         -c './fe' 'nvim'
tmux split-window -dt $session:$window -l 8    -c './fe' 'yarn dev; zsh -i'
tmux split-window -dt $session:$window.2 -h    -c './fe'

window=2
tmux new-window   -dt $session:$window -n 'be' -c './be' 'nvim'
tmux split-window -dt $session:$window -l 8    -c './be' 'air; zsh -i'
tmux split-window -dt $session:$window.2 -h    -c './be'

tmux switch-client -t $session

