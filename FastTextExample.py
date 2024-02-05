import fasttext
import tempfile
from DBAdapter import *


def encodeWords(db, modelName):
    resultSet = db.selectWords()
    #filename = "scripture.text" ## This should
    filename = os.path.join(tempfile.mkdtemp(), "scripture.text")
    print("Words to encode in", filename)
    file = open(filename, "w")
    first = 0
    for (word_id, word, punct, src_word) in resultSet:
        file.write(word)
        if punct != None:
            file.write(punct)      
        file.write(" ")
    file.close()
    model = fasttext.train_unsupervised(filename, "cbow")
    model.save_model(modelName)
    for (word_id, word, punct, src_word) in resultSet:
        word_enc = model.get_word_vector(word)
        #print(word, type(word_enc.dtype), word_enc.shape)
        db.addWordEncoding(word_id, word_enc)
    db.updateWordEncoding()




db = DBAdapter("ENG", 3, "Excel")
encodeWords(db, "ENG_3_Excel.model")




