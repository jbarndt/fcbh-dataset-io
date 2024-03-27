#!/bin/zsh 

if [ -z "$1" ]
then
  echo "Usage: RunScriptUSXDiff.sh  bibleId"
  exit 1 
fi

bibleId=$1
echo $bibleId
database1=$bibleId"_SCRIPT.db"
python3 ExportUtility.py $database1

database2=$bibleId"_USXEDIT.db"
python3 ExportUtility.py $database2

file1=${FCBH_DATASET_DB}/${database1/%.db/.txt}
file2=${FCBH_DATASET_DB}/${database2/%.db/.txt}
diff $file1 $file2 | wc -l
diff -y --suppress-common-lines $file1 $file2 | more
echo "Done!"
echo "diff -y --suppress-common-lines $file1 $file2"
