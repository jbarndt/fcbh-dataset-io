#!/bin/zsh 

if [ -z "$1" ]
then
  echo "Usage: RunScriptUSXCompare.sh  bibleId"
  exit 1 
fi

bibleId=$1
echo $bibleId
python3 ScriptAdapter.py $bibleId
database1=$bibleId"_SCRIPT.db"
python3 ExportUtility.py $database1

$HOME/Documents/go2/bin/usx $bibleId
database2=$bibleId"_USXEDIT.db"
python3 WordParser.py $database2
python3 ExportUtility.py $database2

file1=${FCBH_DATASET_DB}/${database1/%.db/.txt}
file2=${FCBH_DATASET_DB}/${database2/%.db/.txt}
diff $file1 $file2 | wc -l
diff -y --suppress-common-lines $file1 $file2 | more
echo "Done!"
echo "diff -y --suppress-common-lines $file1 $file2"
