mkdir ../../$1
cp -R Template/ ../../$1
LANG=C
find ../../$1 -type f -exec sed -i .tmp "s/\[PROJECT_NAME\]/$1/g" {} +
find ../../$1/ -iname "*.tmp" -exec rm '{}' ';'
for i in ../../$1/*xcodeproj; do mv "$i" "`echo $i | sed \"s/Template/$1/g\"`"; done
for i in ../../$1/*vcxproj*; do mv "$i" "`echo $i | sed \"s/Template/$1/g\"`"; done
for i in ../../$1/*/*; do if [ -f "$i" ]; then mv "$i" "`echo $i | sed \"s/Template/$1/g\"`"; fi  done
for i in ../../$1/*/*/*; do mv "$i" "`echo $i | sed \"s/Template/$1/g\"`"; done
