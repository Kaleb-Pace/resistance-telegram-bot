from textblob import TextBlob
import csv

with open('sen-result.txt', 'a') as outfile:
    outfile.write('polarity')
    with open('result.csv', 'r', encoding='utf8') as csvfile:
        reader = csv.reader(csvfile, delimiter=',')
        for row in reader:
            if row[11] != '\\N':
                blob = TextBlob(row[11])
                
                avg_polartiry = 0
                for sentence in blob.sentences:
                    avg_polartiry += sentence.sentiment.polarity

                avg_polartiry /= len(blob.sentences)
                outfile.write(str(avg_polartiry))
            else:
                outfile.write('0')
            outfile.write('\n')

