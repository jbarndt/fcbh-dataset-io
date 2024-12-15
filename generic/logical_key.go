package generic

import (
	"strconv"
	"strings"
)

// The LogicalKey interface and Reference type are an experimental idea
// for how to abstract out the details for bible references.
// So, that any module that was not specifically about the Bible could use
// The logical key.  This approach would allow modules that are not specifically
// for the Bible to use a Bible reference and be certain it has a consistent meaning.

type LogicalKey interface {
	LogicalKey() string
	IsLogicalKey()
}

type Reference struct {
	PrimaryKey int64
	logicalKey string
	bookId     string
	chapterNum int
	verseStr   string
	verseEnd   string
	chapterEnd int
}

func (r Reference) LogicalKey() string { return r.logicalKey }
func (r Reference) IsLogicalKey()      {}
func (r Reference) BookId() string     { return r.bookId }
func (r Reference) ChapterNum() int    { return r.chapterNum }
func (r Reference) VerseStr() string   { return r.verseStr }
func (r Reference) VerseEnd() string   { return r.verseEnd }
func (r Reference) ChapterEnd() int    { return r.chapterEnd }

func RefByParts(bookId string, chapterNum int, verseStr string, verseEnd string, chapterEnd int) Reference {
	var r Reference
	r.bookId = bookId
	r.chapterNum = chapterNum
	r.verseStr = verseStr
	r.verseEnd = verseEnd
	r.chapterEnd = chapterEnd
	if chapterNum == 0 {
		r.logicalKey = bookId
	} else if verseStr == `` {
		r.logicalKey = bookId + ` ` + strconv.Itoa(chapterNum)
	} else if chapterEnd == 0 && verseEnd == `` {
		r.logicalKey = bookId + ` ` + strconv.Itoa(chapterNum) + `:` + verseStr
	} else if chapterEnd == 0 {
		r.logicalKey = bookId + ` ` + strconv.Itoa(chapterNum) + `:` + verseStr + `-` + verseEnd
	} else {
		r.logicalKey = bookId + ` ` + strconv.Itoa(chapterNum) + `:` + verseStr + `-` + strconv.Itoa(chapterEnd) + `:` + verseEnd
	}
	return r
}

func RefByKey(key string) Reference {
	var r Reference
	r.logicalKey = key
	// This is only partially written
	parts := strings.Split(key, ` `)
	r.bookId = parts[0]
	if len(parts) > 1 {
		parts = strings.Split(parts[1], `:`)
		r.chapterNum, _ = strconv.Atoi(parts[0])
		if len(parts) > 1 {
			r.verseStr = parts[1]
		}
	}
	return r
}

////////////////// Example of other Use of Logical Key ///////////////

type Publish struct {
	key        string
	BookId     string
	ChapterNum int
	Page       int
	Paragraph  int
	Sentence   int
	sequence   int
}

func (p Publish) Key() string   { return p.key }
func (p Publish) IsLogicalKey() {}
func (p Publish) Sequence() int { return p.sequence }

// How will this be used in a generic application?
// The person identifies data by its parts
// Or, the person enters the generic parts
// A Key could first be set to be a book, then we process that
// key for a

// We could permit any number of constructors that is not part of the interface
// When a key is built from parts, we need to update the composite key
// Each time something is set.
// For example, if we permit fields to be set, then we can't have a toKey()
// function until needed, because values can change

// Ergo, I CANNOT have both a full key and the parts in the same type,
// Unless I made it immultable.

// My mms code is always going to use the composite code
// The database will use the parts in the near future, but might later use the
// composite code

// How many constructors would I need for an immultable class
// MRK
// MRK 1
// MRK 1:2
// MRK 1:2-3
// MRK 1:2-2:1
