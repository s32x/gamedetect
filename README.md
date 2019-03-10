# gamedetect

[![CircleCI](https://circleci.com/gh/s32x/gamedetect.svg?style=svg)](https://circleci.com/gh/s32x/gamedetect)

<p align="center">
    <img src="service/static/assets/images/repo.jpg" width="800px" border="0" alt="demo">
</p>

gamedetect is a simple API that uses a trained neural network to identify games that are within the top 100 currently on Twitch (as of March 2019). The full list of supported games can be seen [here](graph/output_labels.txt). The network is trained using [retrain.py](https://github.com/tensorflow/hub/blob/master/examples/image_retraining/retrain.py) which uses InceptionV3 as a pre-trained network. Honestly I'm still at the point where I have no idea what the hell I'm talking about so please bear with me.

## Try this...

Send a POST request with a (relatively clear) game screenshot (one in the supported list) in the "image" field of a form to gamedetect.io.

For example:

```bash
curl -X POST https://gamedetect.io -F image=MY_GAME_SCREENSHOT.png
```

Excited yet? I sure am! gamedetect is a fun project I've been playing with in my free time to learn about Computer Vision, Neural Networks, and Tensorflow. It's sort of my own hello world app that also could potentially serve a real use-case on Twitch or any other streaming platform that requires broadcasters to categorize their stream. That being said, I'm still very much a beginner to all of this and I'm sure I'm doing a number of things wrong - feel free to let me know in the issues if you'd like.