## scrape.py

Used to scrape game frames from Twitch streams using streamlink, and opencv

```
python scrape.py
```

Running the above command in the scripts directory dumps a collection of the top 100 game screenshots into a folder called Screenshots. Up to 2000 screenshots will be downloaded. You will need to weed out any mis-categorized screenshots.

## retrain.py

Used for training the neural network using the scraped game frame dataset

```
wget https://raw.githubusercontent.com/tensorflow/hub/master/examples/image_retraining/retrain.py

python retrain.py --image_dir=Screenshots --output_graph=output_graph.pb --output_labels=output_labels.txt --how_many_training_steps=100000 --print_misclassified_test_images
```

Running the above command in the scripts directory to train the neural network and output the graph and label files to `output_graph.pb` and `output_graph.txt`. These can be used in the root graph directory 