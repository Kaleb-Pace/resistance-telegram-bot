"""Demonstrates how to make a simple call to the Natural Language API."""

import argparse
import time
import sys

from google.cloud import language
from google.cloud.language import enums
from google.cloud.language import types
import csv


def print_result(annotations):
    score = annotations.document_sentiment.score
    magnitude = annotations.document_sentiment.magnitude
    return '{}, {}'.format(score, magnitude)


def analyze(text):
    """Run a sentiment analysis request on text."""
    client = language.LanguageServiceClient()

    document = types.Document(
        content=text,
        type=enums.Document.Type.PLAIN_TEXT)
    annotations = client.analyze_sentiment(document=document)
    return print_result(annotations)


def runfile(index):
    with open('gsen-result' + str(index) + '.txt', 'a') as outfile:
        outfile.write('messageId, polarity, magnitude\n')
        with open('result.csv', 'r', encoding='utf8') as csvfile:
            reader = csv.reader(csvfile, delimiter=',')
            i = 0
            for row in reader:
                if i >= index and i < index + 1000:
                    outfile.write(str(i) + ", ")
                    if row[11] != '\\N':
                        try:
                            outfile.write(analyze(row[11]))
                        except Exception as e:
                            print(e)
                            outfile.write('0, 0')
                        time.sleep(0.15)
                    else:
                        outfile.write('0, 0')
                    outfile.write('\n')
                    print(i)
                i += 1


if __name__ == '__main__':
    print("running " + str(sys.argv[1]))
    runfile(int(sys.argv[1]))

