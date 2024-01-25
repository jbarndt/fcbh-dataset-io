import fasttext
from DBAdapter import *


def encodeWords(db):
    resultSet = db.selectScript()
    filename = "scripture.text"
    file = open(filename, "w")
    first = 0
    for (id, book_id, chapter_num, script_num, word_seq, 
            verse_num, usfm_style, person, word, punct) in resultSet:
        file.write(word)
        if punct != None:
            file.write(punct)      
        file.write(" ")
    file.close()
    sys.exit(1)
    print("start model")
    model = fasttext.train_unsupervised(filename, "cbow")
    print("model finished")
    model.save("sonnet.model")
    print("model saved")
    for (id, book_id, chapter_num, script_num, word_seq, 
            verse_num, usfm_style, person, word) in resultSet:
        word_enc = model.get_word_vector(word)
        print(word, type(word.dtype), word.shape)
        db.updateEncoding(id, word_enc)




db = DBAdapter("ENG", 1, "Sonnet")
encodeWords(db)




