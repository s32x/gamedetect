# Gamedetect

[![CircleCI](https://circleci.com/gh/s32x/gamedetect.svg?style=svg)](https://circleci.com/gh/s32x/gamedetect)

Gamedetect is a simple API that uses a trained neural network to identify games that are within the top 100 currently on Twitch (as of March 2019). The full list of supported games can be seen here.

## Try this...

Send a POST request with a (relatively clear) game screenshot in the "image" field of a form to gamedetect.io.

For example:

```bash
curl -X POST https://gamedetect.io -F image=MY_GAME_SCREENSHOT.png
```

Excited yet? I sure am! Gamedetect is a fun project I've been playing with in my free time to learn about Computer Vision, Neural Networks, and Tensorflow. It's sort of my own hello world app that also could potentially serve a real use-case on Twitch or any other streaming platform that requires broadcasters to categorize their stream in the correct directory. I'm still very much a beginner to all of this and I'm sure I'm doing something wrong so feel free to let me know in the issues.