import cv2
import streamlink
import requests
import os
import re
import random
import re

headers = {
    "Accept": "application/vnd.twitchtv.v5+json", 
    "Client-ID": "d4uvtfdr04uq6raoenvj7m86gdk16v",
}


def get_top_games(limit):
    url = "https://api.twitch.tv/kraken/games/top?limit=" + str(limit)
    r = requests.get(url, headers=headers)
    return r.json()["top"]


def get_streams(game, page):
    url = "https://api.twitch.tv/kraken/streams?limit=100&game=" + game
    if page > 1:
        url = url + "&offset=" + str((page-1) * 100)
    r = requests.get(url, headers=headers)
    if "streams" in r.json():
        return r.json()["streams"]
    else:
        return []


def get_stream_url(channel):
    stream_urls = streamlink.streams("https://www.twitch.tv/" + channel)
    if len(stream_urls) == 0:
        print("No streams found")
        return
    return stream_urls["best"].url


def get_frame(stream_url):
    # Create a VideoCapture that will read the first frame from a stream URL
    cap = cv2.VideoCapture(stream_url)
    success, frame = cap.read()
    if not success:
        raise "Failed to read frame from passed stream URL"
    cap.release()
    return frame


def save_frame(out_dir, stream_url):
    try:
        frame = get_frame(stream_url)
    except:
        raise "Failed to retrieve frame"
    out_path = os.path.join(out_dir, str(random.randint(1, 9999999999)) + ".jpg")
    cv2.imwrite(out_path, frame)


pages = range(1, 25)

def download_game_images(game_name, max_images):
    print("Download game images for " + str(game_name))
    game_dir = os.path.join("Screenshots", game_name.replace(":", ""))

    # Create the game directory if it doesn"t already exist
    try:
        os.stat(game_dir)
    except:
        os.mkdir(game_dir)

    num_images = len(os.listdir(game_dir))
    if num_images > max_images:
        print("Already have enough images for " + game_name)
        return
    print("Found " + str(num_images) + " for game " + game_name)

    for page in pages:
        print("Getting streams for " + game_name + " from page " + str(page))
        streams = get_streams(game_name, page)

        for stream in streams:
            channel_name = stream["channel"]["name"]
            game_path = os.path.join("Screenshots", game_name.replace(":", ""))

            print("Saving " + game_name + " frame from " + channel_name + " (" + str(num_images) + " of " + str(max_images) + ")")
            try:
                save_frame(game_path, get_stream_url(channel_name))
            except:
                print("Failed to save frame")

            num_images = len(os.listdir(game_dir))
            if num_images >= max_images:
                return


def download_frames():
    max_images = 1050

    print("Getting top games")
    games = get_top_games(100)

    for game in games:
        game_name = game["game"]["name"]

        if game_name == "Music & Performing Arts":
            continue
        if game_name == "Just Chatting":
            continue
        if game_name == "Art":
            continue
        if game_name == "Poker":
            continue
        if game_name == "Retro":
            continue
        if game_name == "ASMR":
            continue
        if game_name == "Always On":
            continue
        if game_name == "Food & Drink":
            continue


        print("Downloading game screenshots")
        download_game_images(game_name, max_images)

                        
def main():
    download_frames()


main()
