set -x

./fetch https://www.google.com

ls *.html

./fetch -metadata https://www.google.com

./fetch https://www.google.com https://autify.com 

ls *.html

./fetch -metadata https://www.google.com https://autify.com 