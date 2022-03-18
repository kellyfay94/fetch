rm -rf *.html *.json

set -x

./fetch https://www.google.com

ls *.html

./fetch -metadata https://www.google.com

rm -rf www.google.com

./fetch https://www.google.com https://www.autify.com 

ls *.html

./fetch -metadata https://www.google.com https://www.autify.com 