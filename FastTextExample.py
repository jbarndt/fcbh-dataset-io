import fasttext
import tempfile
from DBAdapter import *


def encodeWords(db, modelName):
    resultSet = db.selectScript()
    #filename = "scripture.text" ## This should
    filename = os.path.join(tempfile.mkdtemp(), "scripture.text")
    print("Words to encode in", filename)
    file = open(filename, "w")
    first = 0
    for (id, book_id, chapter_num, script_num, usfm_style, person, word_seq, 
            verse_num, word, punct) in resultSet:
        file.write(word)
        if punct != None:
            file.write(punct)      
        file.write(" ")
    file.close()
    print("start model")
    model = fasttext.train_unsupervised(filename, "cbow")
    print("model finished")
    model.save_model(modelName)
    print("model saved")
    for (id, book_id, chapter_num, script_num, word_seq, 
            verse_num, usfm_style, person, word, punct) in resultSet:
        word_enc = model.get_word_vector(word)
        print(word, type(word_enc.dtype), word_enc.shape)
        db.addWordEncoding(id, word_enc)
    db.updateWordEncoding()




db = DBAdapter("ENG", 2, "WEB")
encodeWords(db, "ENG_2_WEB")




