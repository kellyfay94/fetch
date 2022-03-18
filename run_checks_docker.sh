docker build -t fay_chris_fetch:latest .

docker run --rm -it fay_chris_fetch > output.txt

docker rmi fay_chris_fetch:latest

clear

cat output.txt