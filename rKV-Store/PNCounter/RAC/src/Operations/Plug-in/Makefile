build:
	g++ -shared -o statistics.so statistics.cpp -fPIC -ldl
	g++ -shared -o sha256.so -fPIC sha256.cpp -lcrypto -ldl
clean:
	rm *.so