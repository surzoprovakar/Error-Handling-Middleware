FLAGS = -std=c++11

# location of swi-prolog/include in the device
INCLUDES = -I  /usr/lib/swi-prolog/include/ 
# 1st portion is the path of libtinfo.so, 2nd portion is the path of libswipl.so in the device
LI = -L=/usr/lib/x86_64-linux-gnu/ -L/usr/lib/swi-prolog/lib/x86_64-linux/
# the path here is the path where libtinfo.so is
DL = -Wl,-rpath-link,/usr/lib/x86_64-linux-gnu
# find the path of libswipl.so and libtinfo.so in the device
PRELOAD=/usr/lib/swi-prolog/lib/x86_64-linux/libswipl.so



build:
	g++ -c -fPIC Persist-Manager/Persist.cpp Undo-Manager/Single_undo.cpp Undo-Manager/Multiple_undo.cpp -ldl -Wall
	g++ -shared $(LI) $(DL) -o libPersist.so Persist.o Single_undo.o Multiple_undo.o -lrt -ldl
clean:
	@rm -rf *.o libPersist.so DBs/