import spacy

#
# This is not yet working.  I think it is a problem with an ssl module used.
#

# Load the English tokenizer, tagger, parser, NER, and word vectors
nlp = spacy.load("en_core_web_sm")

def preprocess(text):
    # Process the text using Spacy to tokenize and lemmatize the text
    doc = nlp(text)
    # Remove punctuation, spaces, and stop words, and then lower-case all words
    return [token.lemma_.lower() for token in doc if not token.is_punct and not token.is_space and not token.is_stop]

def compare_texts(text1, text2):
    # Preprocess both texts
    tokens1 = set(preprocess(text1))
    tokens2 = set(preprocess(text2))
    
    # Find differences
    diff1 = tokens1 - tokens2
    diff2 = tokens2 - tokens1
    
    return diff1, diff2

# Read your files (replace 'file1.txt' and 'file2.txt' with your actual file names)
with open('file1.txt', 'r') as file:
    text1 = file.read()

with open('file2.txt', 'r') as file:
    text2 = file.read()

# Compare the files
differences1, differences2 = compare_texts(text1, text2)

print("Words in the first text not in the second:", differences1)
print("Words in the second text not in the first:", differences2)